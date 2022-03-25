package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	users []User
)

func init(){
	InitDB()
	InitialMigration()
}

type Config struct{
	DB_Username string
	DB_Password string
	DB_Port string
	DB_Host string
	DB_Name string
}

type User struct{
	gorm.Model
	ID int `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func InitDB(){
	config := Config{
		DB_Username: "root",
		DB_Password: "Sautmanurung234",
		DB_Port: "3306",
		DB_Host: "localhost",
		DB_Name: "users",
	}
	connectionString := 
	fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DB_Username,
		config.DB_Password,
		config.DB_Host,
		config.DB_Port,
		config.DB_Name,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(connectionString),&gorm.Config{})
	if err != nil {
		panic(err.Error())
	}	
}

func InitialMigration(){
	DB.AutoMigrate(&User{})
}

// get all users
func GetUsersController(c echo.Context) error {
	if err := DB.Find(&users).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages" : "success get all users",
		"users" : users,
	})
}

// get user by id
func GetUserController(c echo.Context) error{
	id , _ := strconv.Atoi(c.Param("id"))
	for _, user := range users{
		if user.ID == id{
			return c.JSON(http.StatusOK, map[string]interface{}{
				"messages" : "success get user",
				"user" : user,
			})
		} else if err := DB.Save(&user).Error; err != nil{
			return echo.NewHTTPError(http.StatusBadRequest,err.Error())
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"messagess" : "User Not Found In Databases",
	})
}

func CreateUserController(c echo.Context) error{
	user := User{}
	c.Bind(&user)

	if err := DB.Save(&user).Error; err != nil{
		return echo.NewHTTPError(http.StatusBadRequest,err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message" : "Success Create New User",
		"User" : user,
	})
}

// delete user by id 
func DeleteUserController(c echo.Context) error{
	id, _ := strconv.Atoi(c.Param("id"))
	for index, user := range users {
		if user.ID == id {
			users = append(users[:index], users[index+1:]...)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"messages" : "success delete user",
				"users" : users,
			})
		} else if err := DB.Save(&user).Error; err != nil{
			return echo.NewHTTPError(http.StatusBadRequest,err.Error())
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"messagess" : "User Not Found In Databases",
	})
}

// Update user by id
func UpdateUserController(c echo.Context) error{
	id , _:= strconv.Atoi(c.Param("id"))
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user User
	user.Name = name
	user.Email = email
	user.Password = password
	for index, user := range users {
		if user.ID == id {
			users[index].Name = name
			users[index].Email = email
			users[index].Password = password
			return c.JSON(http.StatusOK, map[string]interface{}{
				"messages" : "success update user",
				"users" : users,
			})
		} else if err := DB.Save(&user).Error; err != nil{
			return echo.NewHTTPError(http.StatusBadRequest,err.Error())
		}
	}
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"messages" : "Nothing To Delete",
	})
}

func main(){
	// Create a new echo instance
	e := echo.New()
	// Route / to handler function
	e.GET("/users",GetUsersController)
	e.GET("/users/:id",GetUserController)
	e.POST("/users",CreateUserController)
	e.DELETE("/users/:id",DeleteUserController)
	e.PUT("/users/:id",DeleteUserController)

	// start the server, and log if it fails
	e.Logger.Fatal(e.Start(":1234"))
}