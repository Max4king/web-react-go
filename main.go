package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	// For json
)

type Data struct {
	User string
	Pass string
	Role string
}

func readDatabase() []Data {
	content, err := os.ReadFile("./user-login.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload []Data
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}

func getUserRole(c echo.Context) error {
	name := c.Param("name")
	password := c.Param("password")
	payload := readDatabase()
	for _, user := range payload {
		if user.User == name {
			if user.Pass == password {
				return c.String(http.StatusOK, user.Role)
			}
			return c.String(http.StatusUnauthorized, "Incorrect password")
		}
	}

	return c.String(http.StatusNotFound, "User not found")
}

func registerUser(c echo.Context) error {
	name := c.FormValue("name")
	pwd := c.FormValue("password")
	role := c.FormValue("role")
	payload := readDatabase()
	for _, user := range payload {
		if user.User == name {
			return c.String(http.StatusConflict, "User already Exist")
		}
	}
	newUser := Data{name, pwd, role}
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
	e.POST("/login", getUserRole)
	e.POST("/register", registerUser)
	e.Logger.Fatal(e.Start(":1323"))
}
