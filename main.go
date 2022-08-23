package main

import (
	// "fmt"

	"fmt"
	"log"
	"net/http"
	"strings"
	"tesa/auth"
	"tesa/campaign"
	"tesa/handler"
	"tesa/helper"
	"tesa/user"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/cobago?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	campaignsRepository := campaign.NewRepository(db)
	// campaigns, err := campaignsRepository.FindByUserID(21)

	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println(len(campaigns))
	// for _, campaign := range campaigns {
	// 	fmt.Println(campaign.Name)
	// 	if len(campaign.CampaignImages) > 0 {
	// 		fmt.Println(campaign.CampaignImages[0].FileName)
	// 	}
	// }

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignsRepository)

	campaigns, _ := campaignService.FindCampaigns(23)
	fmt.Println(len(campaigns))

	authService := auth.NewService()

	// fmt.Println(authService.GenerateToken(1001))
	// userService.SaveAvatar(1, "images/avatar.png")

	// token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyMH0.XyGJF39ZDuW-s7-TutG-YD2rhjvjFVuViYiyABZIXeUs")
	// if err != nil {
	// 	fmt.Println("Error")
	// 	fmt.Println("Error")
	// 	fmt.Println("Error")
	// }

	// if token.Valid {
	// 	fmt.Println("valid")
	// 	fmt.Println("valid")
	// 	fmt.Println("valid")
	// } else {
	// 	fmt.Println("invalid")
	// 	fmt.Println("invalid")
	// 	fmt.Println("invalid")
	// }

	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/checkEmail", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	router.Run()

	//input dari user
	//handler, mapping input dari user ke struct user
	//service, melakukan mapping dari struct input ke struct user
	//repository, melakukan penyimpanan data ke database
	//db

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

// middleware
// ambil nilai header Authorization: bearer tokentoken
// dari header authorization, kita ambil nilai tokennya saja
// kita validasi token
// kita ambil user_id
// ambil user dari db berdasarkan user_id lewat service
// kita set contect isinya user
