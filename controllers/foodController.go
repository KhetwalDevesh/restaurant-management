package controller

import (
	"context"
	"errors"
	"fmt"
	config "github.com/KhetwalDevesh/restaurant-management/database"
	"github.com/KhetwalDevesh/restaurant-management/helpers"
	"github.com/KhetwalDevesh/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
	"time"
)

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var foodItems []*models.Food
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		offset, limit := helpers.GetLimitOffset(20, page)
		config.GetDB().Model(models.Food{}).Offset(offset).Limit(limit + 1).Find(&foodItems)
		nextPage := false
		if len(foodItems) > limit {
			nextPage = true
			foodItems = foodItems[0 : len(foodItems)-1]
		}
		type GetFoodsResponse struct {
			FoodItems []*models.Food
			NextPage  bool
		}

		getFoodsResponse := GetFoodsResponse{
			FoodItems: foodItems,
			NextPage:  nextPage,
		}
		c.JSON(http.StatusOK, getFoodsResponse)
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")
		var food models.Food
		err := config.GetDB().Model(models.Food{}).Where("food_id = ?", foodId).First(&food).Error
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the food item"})
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		isUserAdmin, _ := c.Get("isAdmin")
		if isUserAdmin == false {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You need to be an admin to create a food"})
			return
		}
		var menu *models.Menu
		var food *models.Food

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := config.GetDB().Where("id = ?", food.MenuId).First(&menu).Error
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		food.CreatedAt = time.Now()
		food.UpdatedAt = time.Now()
		var num = toFixed(food.Price, 2)
		food.Price = num

		// validate food data before storing it in db
		if err := validate.Struct(food); err != nil {
			msg := fmt.Sprintf("Food data invalidated : %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		err = config.GetDB().Model(&models.Food{}).Create(&food).Error
		if err != nil {
			msg := fmt.Sprintf("Error creating the food")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Food created successfully!")
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		foodID := c.Param("food_id")

		// Check if the menu exists
		var existingFood models.Food
		if err := config.GetDB().Where("id = ?", foodID).First(&existingFood).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
			return
		}

		if food.Name != "" {
			existingFood.Name = food.Name
		}

		if food.Price != 0 {
			existingFood.Price = food.Price
		}

		if food.Image != "" {
			existingFood.Image = food.Image
		}
		var menuForFoodToBeUpdated *models.Menu
		if food.MenuId != 0 {
			if err := config.GetDB().Model(&models.Menu{}).Where("menu_id = ?", food.MenuId).First(&menuForFoodToBeUpdated).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": "GetMenu got error"})
				return
			}
			existingFood.MenuId = food.MenuId
		}

		existingFood.UpdatedAt = time.Now()

		// validate food data before saving it
		if err := validate.Struct(existingFood); err != nil {
			msg := fmt.Sprintf("Food data invalidated : %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Save changes to the database
		if err := config.GetDB().Save(&existingFood).Error; err != nil {
			msg := "Food update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Food updated successfully"})
	}
}
