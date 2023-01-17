package main

import (
	"crowdfounding/auth"
	"crowdfounding/campaign"
	"crowdfounding/handler"
	"crowdfounding/helper"
	"crowdfounding/payment"
	"crowdfounding/transaction"
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
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := gin.Default()

	router.Static("/images", "./images")

	api := router.Group("/api/v1")
	// pendaftaran user
	api.POST("/users", userHandler.RegisterUser)
	// login user
	api.POST("/sessions", userHandler.Login)
	// pengecekan ketersedian email
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	// upload avatars
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	// campaign routes
	// show all campaigns
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	// detail of campaign
	api.GET("/campaign/:id", campaignHandler.GetCampaign)
	// create a campaign
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	// update a campaign
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	// upload image campaign
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	// transactions routes
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
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

		userData, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", userData)
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
