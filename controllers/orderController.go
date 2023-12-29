package controller

import (
	"context"
	"errors"
	"fmt"
	config "github.com/KhetwalDevesh/restaurant-management/database"
	"github.com/KhetwalDevesh/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItems []*models.Order
		err := config.GetDB().Model(models.Order{}).Find(&orderItems)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "GetOrders got error"})
		}
		c.JSON(http.StatusOK, orderItems)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		orderId := c.Param("order_id")
		var order models.Order
		err := config.GetDB().Model(models.Order{}).Where("id = ?", orderId).First(&order).Error
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the order item"})
		}
		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		var order models.Order
		// orderDate is not send, current time and date is assigned
		order.OrderDate = time.Now()
		// assign the userId to the order to be created
		userId, userIdExists := c.Get("uid")
		if userIdExists {
			order.UserId = userId.(uint32)
		}

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// check if the table exists from where the order is made
		if order.TableId != 0 {
			err := config.GetDB().Model(models.Table{}).Where("id = ?", order.TableId).First(&table).Error
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("table was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
		}

		order.CreatedAt = time.Now()
		order.UpdatedAt = time.Now()

		// validate order data before storing it in db
		if err := validate.Struct(order); err != nil {
			msg := fmt.Sprintf("Order data invalidated : %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		err := config.GetDB().Model(&models.Order{}).Create(&order).Error
		if err != nil {
			msg := fmt.Sprintf("Error creating the order")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Order created successfully!")
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var order *models.Order

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		orderID := c.Param("order_id")

		// Check if the menu exists
		var existingOrder *models.Order
		if err := config.GetDB().Where("id = ?", orderID).First(&existingOrder).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		var tableForOrderToBeUpdated *models.Table
		if order.TableId != 0 {
			if err := config.GetDB().Model(&models.Table{}).Where("table_id = ?", order.TableId).First(&tableForOrderToBeUpdated).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": "GetTable got error"})
				return
			}
			existingOrder.TableId = order.TableId
		}

		existingOrder.UpdatedAt = time.Now()

		// validate order data before storing it in db
		if err := validate.Struct(existingOrder); err != nil {
			msg := fmt.Sprintf("Order data invalidated : %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Save changes to the database
		if err := config.GetDB().Save(&existingOrder).Error; err != nil {
			msg := "Order update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
	}
}

func OrderItemOrderCreator(order models.Order) uint32 {
	// Use GORM to create the order record
	db := config.GetDB()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	db.Create(&order)
	// Convert the order ID to hex format
	return order.ID
}
