package main

import (
	"crowdfounding/auth"
	"crowdfounding/handler"
	"crowdfounding/helper"
	"crowdfounding/user"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strings"

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
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	router.Run()

}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}

// Middleware
// ambil nilai header Authorization: Bearer tokentokentoken
// dari header authorization, ambil nilai tokennya saja
// lakukan validasi token
// ambil nilai user_id
// ambil user dari db berdasarkan user_id lewat service
// set context isinya user

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
