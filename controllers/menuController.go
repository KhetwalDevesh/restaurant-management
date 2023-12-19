package controller

import (
	"context"
	"fmt"
	config "github.com/KhetwalDevesh/restaurant-management/database"
	"github.com/KhetwalDevesh/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menus *models.Menu
		err := config.GetDB().Model(models.Menu{}).Find(&menus).Error
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "GetMenus got error"})
		}
		c.JSON(http.StatusOK, menus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		menuId := c.Param("menu_id")
		var menu models.Menu
		err := config.GetDB().Model(models.Menu{}).Where("id = ?", menuId).First(&menu).Error
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the menu item"})
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		menu.CreatedAt = time.Now()
		menu.UpdatedAt = time.Now()
		err := config.GetDB().Model(&models.Food{}).Create(menu).Error
		if err != nil {
			msg := fmt.Sprintf("Error creating the menu")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Menu created successfully!")
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		menuID := c.Param("menu_id")

		// Check if the menu exists
		var existingMenu models.Menu
		if err := config.GetDB().Where("id = ?", menuID).First(&existingMenu).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
			return
		}

		// Validate time span
		if !menu.StartDate.IsZero() && !menu.EndDate.IsZero() {
			if !inTimeSpan(menu.StartDate, menu.EndDate, time.Now()) {
				msg := "Kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
		}

		// Update menu fields
		if !menu.StartDate.IsZero() {
			existingMenu.StartDate = menu.StartDate
		}
		if !menu.EndDate.IsZero() {
			existingMenu.EndDate = menu.EndDate
		}
		if menu.Name != "" {
			existingMenu.Name = menu.Name
		}
		if menu.Category != "" {
			existingMenu.Category = menu.Category
		}

		existingMenu.UpdatedAt = time.Now()

		// Save changes to the database
		if err := config.GetDB().Save(&existingMenu).Error; err != nil {
			msg := "Menu update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Menu updated successfully"})
	}
}

// inTimeSpan checks if a given time is within a time span
func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}
