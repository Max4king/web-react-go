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

func csvToJson(data []byte) ([]map[string]string, error) {
	r := csv.NewReader(strings.NewReader(string(data)))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	headers := records[0]
	var jsonArray []map[string]string
	for _, record := range records[1:] {
		dataMap := make(map[string]string)
		for index, value := range record {
			dataMap[headers[index]] = value
		}
		jsonArray = append(jsonArray, dataMap)
	}
	return jsonArray, nil
}

func getData(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok || user == nil {
		return c.JSON(http.StatusUnauthorized, "User not authenticated")
	}
	claims := user.Claims.(jwt.MapClaims)
	role, ok := claims["roles"].(string)
	if !ok {
		return c.JSON(http.StatusForbidden, "Invalid user role")
	}

	var files []string
	switch role {
	case "admin":
		files = append(files, "admin.csv", "client.csv")
	case "client":
		files = append(files, "client.csv")
	default:
		return c.JSON(http.StatusForbidden, "Access denied")
	}

	var combinedData []map[string]string
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error reading data file: "+file)
		}
		jsonData, err := csvToJson(data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error converting CSV to JSON: "+err.Error())
		}
		combinedData = append(combinedData, jsonData...)
	}
	fmt.Println(combinedData)

	return c.JSON(http.StatusOK, combinedData)
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
	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // For frontend's host and port
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Authorization"},
	}))
	e.POST("/login", login)
	e.POST("/register", register)

	g := e.Group("/api")
	g.Use(JWTMiddleware)
	g.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // For frontend's host and port
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Authorization"},
	}))
	g.GET("/data", getData)
	e.Logger.Fatal(e.Start(":1323"))
}
