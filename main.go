package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
				return c.JSON(http.StatusOK, userData.Role)
			}
			return c.JSON(http.StatusUnauthorized, "Incorrect password")
		}
	}
	fmt.Println(user)
	return c.JSON(http.StatusNotFound, "User not found")
}

func register(c echo.Context) error {
	var user Data
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "Invalid data")
	}
	payload, err := readDatabase()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	for _, userData := range payload {
		if userData.User == user.User {
			return c.String(http.StatusConflict, "User already Exist")
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
	return c.String(http.StatusOK, "Registered")
}

func main() {
	fmt.Println("vim-go")
	e := echo.New()

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // For frontend's host and port
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.POST("/login", login)
	e.POST("/register", register)
	e.Logger.Fatal(e.Start(":1323"))
}
