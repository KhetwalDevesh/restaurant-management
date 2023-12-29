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

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var allTables []*models.Table
		err := config.GetDB().Model(models.Table{}).Find(&allTables)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "GetTables got error"})
		}
		c.JSON(http.StatusOK, allTables)
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		tableId := c.Param("table_id")
		var table models.Table
		err := config.GetDB().Model(models.Table{}).Where("id = ?", tableId).First(&table).Error
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the table item"})
		}
		c.JSON(http.StatusOK, table)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		isUserAdmin, _ := c.Get("isAdmin")
		if isUserAdmin == false {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You need to be an admin to create a table"})
			return
		}
		var table *models.Table
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		table.CreatedAt = time.Now()
		table.UpdatedAt = time.Now()

		// validate table data before storing it in db
		if err := validate.Struct(table); err != nil {
			msg := fmt.Sprintf("Table data invalidated : %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		err := config.GetDB().Model(models.Table{}).Create(&table).Error
		if err != nil {
			msg := fmt.Sprintf("Error creating the table")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Table created successfully!")
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var table models.Table

		tableID := c.Param("table_id")

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj = make(map[string]interface{})

		if table.NumberOfGuests != 0 {
			updateObj["numberOfGuests"] = table.NumberOfGuests
		}

		if table.TableNumber != 0 {
			updateObj["tableNumber"] = table.TableNumber
		}

		table.UpdatedAt = time.Now()

		// validate table data before storing it in db
		if err := validate.Struct(table); err != nil {
			msg := fmt.Sprintf("Table data invalidated : %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Update fields in the database
		if err := config.GetDB().Model(&models.Table{}).Where("id = ?", tableID).Updates(updateObj).Error; err != nil {
			msg := "Table item update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Table item updated successfully"})
	}
}
