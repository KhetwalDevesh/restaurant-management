package helpers

import (
	"fmt"
	config "github.com/KhetwalDevesh/restaurant-management/database"
	"github.com/KhetwalDevesh/restaurant-management/models"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"time"
)

type SignedDetails struct {
	Email   string
	Name    string
	UserId  uint32
	IsAdmin bool
	jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("TOKEN_SECRET_KEY")

func GenerateAllTokens(email string, name string, userId uint32, isAdmin bool) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:   email,
		Name:    name,
		UserId:  userId,
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(240)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId uint32) {
	// Assuming you have a User model
	user := models.User{
		ID:           userId,
		Token:        signedToken,
		RefreshToken: signedRefreshToken,
		UpdatedAt:    time.Now(),
	}

	// Use GORM to update the user details
	if err := config.GetDB().Model(&user).Updates(user).Error; err != nil {
		log.Panic(err)
		return
	}
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	//the token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprint("token is expired")
		msg = err.Error()
		return
	}
	return claims, msg
}
