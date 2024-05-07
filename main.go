package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Data struct {
	User string
	Pass string
	Role string
}

func readDatabase() ([]Data, error) {
	content, err := os.ReadFile("./user-login.json")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	var payload []Data
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return nil, fmt.Errorf("error during Unmarshal: %v", err)
	}
	return payload, nil
}

type User struct {
	User string `json:"name"`
	Pass string `json:"password"`
}

func createToken(name Data) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name.User
	claims["roles"] = name.Role
	claims["exp"] = time.Now().Add(time.Hour * 24 * 10).Unix()
	t, err := token.SignedString([]byte("secret"))
	return t, err
}

func login(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid data")
	}
	payload, err := readDatabase()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	for _, userData := range payload {
		if userData.User == user.User {
			if userData.Pass == user.Pass {
				token, err := createToken(userData)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, "Error in creating token")
				}
				return c.JSON(http.StatusOK, token)
			}
			return c.JSON(http.StatusUnauthorized, "Incorrect password")
		}
	}
	fmt.Println(user)
	return c.JSON(http.StatusNotFound, "User not found")
}

type newUser struct {
	User string `json:"name"`
	Pass string `json:"password"`
	Role string `json:"role"`
}

func register(c echo.Context) error {
	var user newUser
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid data")
	}
	payload, err := readDatabase()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	for _, userData := range payload {
		if userData.User == user.User {
			return c.JSON(http.StatusConflict, "User already Exist")
		}
	}
	newUser := Data{user.User, user.Pass, user.Role}
	payload = append(payload, newUser)

	newPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Error in converting user-login with the new user", err)
	}
	err = os.WriteFile("user-login.json", newPayload, 0644)
	if err != nil {
		log.Fatal("Error in Writing to the file", err)
	}
	fmt.Println("User Registered")
	return c.JSON(http.StatusOK, "Registered")
}

func csvToJson(data []byte) (string, error) {
	// Convert byte data to a reader, which can be used by csv.NewReader
	r := csv.NewReader(strings.NewReader(string(data)))
	records, err := r.ReadAll()
	if err != nil {
		return "", err
	}

	// The first record is assumed to be headers
	headers := records[0]
	var jsonArray []map[string]string

	for _, record := range records[1:] {
		dataMap := map[string]string{}
		for index, value := range record {
			dataMap[headers[index]] = value
		}
		jsonArray = append(jsonArray, dataMap)
	}

	jsonData, err := json.Marshal(jsonArray)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func getData(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok || user == nil {
		return c.JSON(http.StatusUnauthorized, "User not authenticated")
	}
	claims := user.Claims.(jwt.MapClaims)
	// name := claims["name"].(string)
	roles := claims["roles"].(string)

	var dataFiles []string

	if roles == "admin" {
		dataFiles = append(dataFiles, "admin.csv")
	} else if roles == "client" {
		dataFiles = append(dataFiles, "client.csv") // Client only gets client.csv.
	} else {
		return c.JSON(http.StatusForbidden, "Access denied") // If role is neither admin nor client.
	}

	var combinedData strings.Builder
	for _, dataFile := range dataFiles {
		data, err := os.ReadFile(dataFile)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error reading data file: "+dataFile)
		}
		combinedData.WriteString(string(data) + "\n")
	}
	// fmt.Println(combinedData.String())
	return c.Blob(http.StatusOK, "text/csv", []byte(combinedData.String()))
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
		}

		// Typically Authorization header is in the format `Bearer token`
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid Authorization header format")
		}

		tokenStr := splitToken[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
		}

		// Set user context from token
		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "error parsing token claims")
		}

		c.Set("user", token) // Now `user` is set to the jwt.Token containing claims
		return next(c)
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(JWTMiddleware)
	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // For frontend's host and port
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.POST("/login", login)
	e.POST("/register", register)
	e.GET("/data", getData)
	e.Logger.Fatal(e.Start(":1323"))
}
