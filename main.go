package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello, Users",
			"status":  true,
		})
	}
}

type User struct {
	ID       int
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	listData []User
)

func init() {
	listData = []User{}
}

func regist(newUser User) User {
	if len(listData) == 0 {
		newUser.ID = 1
	} else {
		newUser.ID = listData[len(listData)-1].ID + 1
	}
	listData = append(listData, newUser)

	return newUser
}

func validateData(email string, password string) User {
	for _, val := range listData {
		if val.Email == email && val.Password == password {
			return val
		}
	}
	return User{}
}

func main() {
	var e = echo.New()
	e.Use(middleware.CORS())
	e.GET("/users", Hello())

	e.POST("/users", func(c echo.Context) error {
		var newData User
		err := c.Bind(&newData)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "error from server",
				"status":  false,
			})
		}
		data := regist(newData)

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "success",
			"status":  true,
			"data":    data,
		})
	})

	e.POST("/login", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "congrats!",
			"status":  true,
		})
	}, middleware.BasicAuth(func(email, password string, ctx echo.Context) (bool, error) {
		var u = validateData(email, password)
		if u.ID != 0 {
			return true, nil
		}
		return false, nil
	}))

	e.Start(":8000")
}
