package controller

import (
	"errors"
	"fmt"
	config "github.com/KhetwalDevesh/restaurant-management/database"
	helper "github.com/KhetwalDevesh/restaurant-management/helpers"
	"github.com/KhetwalDevesh/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []*models.User
		var totalUsersCount int64

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		// Fetch total count of users
		config.GetDB().Model(&models.User{}).Count(&totalUsersCount)

		// Fetch paginated users
		config.GetDB().Offset(startIndex).Limit(recordPerPage).Find(&users)

		response := gin.H{
			"total_count": totalUsersCount,
			"user_items":  users,
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		var user models.User

		// Assuming that you have a GORM database instance configured using config.GetDB()
		db := config.GetDB()

		// Find the user by ID
		err := db.Where("id = ?", userId).First(&user).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching user"})
			}
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		// Convert the JSON data coming from Postman to something that Golang understands
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the email has already been used by another user
		var existingUser models.User
		if err := config.GetDB().Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "this email already exists"})
			return
		}

		// Hash password
		password := HashPassword(user.Password)
		user.Password = password
		//user.IsAdmin
		// Create some extra details for the user object - created_at, updated_at, ID
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		// If all is ok, then insert this new user into the users table
		if err := config.GetDB().Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while creating user"})
			return
		}

		// Generate token and refresh token
		token, refreshToken, _ := helper.GenerateAllTokens(user.Email, user.Name, user.ID, user.IsAdmin)
		user.Token = token
		user.RefreshToken = refreshToken

		// validate user data before storing it in db
		if err := validate.Struct(user); err != nil {
			msg := fmt.Sprintf("User data invalidated : %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if err := config.GetDB().Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while generating and assigning token to user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": user})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User

		// Convert the login data from Postman which is in JSON to Golang readable format
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find a user with that email and see if that user even exists
		if err := config.GetDB().Where("email = ?", user.Email).First(&foundUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found, login seems to be incorrect"})
			return
		}

		// Then you will verify the password
		passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// If all goes well, then you'll generate tokens
		token, refreshToken, _ := helper.GenerateAllTokens(foundUser.Email, foundUser.Name, foundUser.ID, foundUser.IsAdmin)

		// Update tokens - token and refresh token using GORM
		helper.UpdateAllTokens(token, refreshToken, foundUser.ID)
		// Return status OK
		foundUser.Token = token
		foundUser.RefreshToken = refreshToken
		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or password is incorrect")
		check = false
	}
	return check, msg
}
