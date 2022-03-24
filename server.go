package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type User struct{
	Id int `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

var(
    users []User
)

// get all users
func GetUsersController(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages" : "success get all users",
		"users" : users,
	})
}

// get user by id
func GetUserController(c echo.Context) error {
	id , _:= strconv.Atoi(c.Param("id"))
	for _, user := range users {
		if user.Id == id {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"messages" : "success get user",
				"user" : user,
			})
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages" : "user not found",
	})
}
// delete user by id
func DeleteUserController(c echo.Context) error{
	id , _:= strconv.Atoi(c.Param("id"))
	for index, user := range users {
		if user.Id == id {
			users = append(users[:index], users[index+1:]...)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"messages" : "success delete user",
				"users" : users,
			})
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages" : "user not found",
	})
}
// update user by id
func UpdateUserController(c echo.Context) error {
	id , _:= strconv.Atoi(c.Param("id"))
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user User
	user.Name = name
	user.Email = email
	user.Password = password
	for index, user := range users {
		if user.Id == id {
			users[index].Name = name
			users[index].Email = email
			users[index].Password = password
			return c.JSON(http.StatusOK, map[string]interface{}{
				"messages" : "success update user",
				"users" : users,
			})
		}
	}
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"messages" : "user not found",
	})
}
// create new user
func CreateUserController(c echo.Context) error {
	// binding data
	user := User{}
	c.Bind(&user)

	if len(users) == 0 {
		user.Id = 1
	} else {
		newId := users[len(users)-1].Id + 1
		user.Id = newId
	}
	users = append(users, user)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages" : "success create user",
		"user" : user,
	})
}

func main(){
	e := echo.New()
	
	// routing with query parameter
	e.GET("/users",GetUsersController)
	e.GET("/users/:id",GetUserController)
	e.POST("/users", CreateUserController)
	e.PUT("/users/:id", UpdateUserController)
	e.DELETE("users/:id", DeleteUserController)

	// START the server, and log if it fails
	e.Logger.Fatal(e.Start(":8000"))
}