package main

import (
	"crowdfounding/auth"
	"crowdfounding/handler"
	"crowdfounding/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/crowd_funding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")
	// pendaftaran user
	api.POST("/users", userHandler.RegisterUser)
	// login user
	api.POST("/sessions", userHandler.Login)
	// pengecekan ketersedian email
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	// upload avatars
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run()

	///////////////////////////
	// fmt.Println("Connection database success")

	// var users []user.User
	// length := len(users)
	// fmt.Println(length)

	// db.Find(&users)
	// length = len(users)
	// fmt.Println(length)

	// router := gin.Default()
	// router.GET("/users", handler)
	// router.Run()
}

// func handler(c *gin.Context) {
// 	dsn := "root:@tcp(127.0.0.1:3306)/crowd_founding?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	var users []user.User
// 	db.Find(&users)

// 	c.JSON(http.StatusOK, users)

// 	// input
// 	// handler mapping input ke struct
// 	// service mapping ke struct user
// 	// repostory save struct User ke db
// 	// db
// }
