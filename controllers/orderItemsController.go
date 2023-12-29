package controller

import (
	"context"
	"fmt"
	config "github.com/KhetwalDevesh/restaurant-management/database"
	"github.com/KhetwalDevesh/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type OrderItemPack struct {
	TableId    uint32
	OrderItems []models.OrderItem
}

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var allOrderItems []*models.OrderItem
		err := config.GetDB().Find(&allOrderItems).Error
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing ordered items"})
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderIdAsString := c.Param("order_id")
		// Convert orderID from string to uint32
		orderId, err := strconv.ParseUint(orderIdAsString, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
			return
		}
		allOrderItems, err := ItemsByOrder(uint32(orderId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing order items by order"})
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func ItemsByOrder(id uint32) ([]*models.OrderItem, error) {
	var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var orderItems []*models.OrderItem
	err := config.GetDB().Preload("Order.Table").Preload("Food.Menu").Preload("Order.User").
		Joins("JOIN orders o ON o.id = order_items.order_id").
		Where("o.id = ?", id).
		Order("order_items.id asc").
		Find(&orderItems).Error
	defer cancel()
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderItemId := c.Param("order_item_id")
		var orderItem *models.OrderItem
		err := config.GetDB().Model(models.OrderItem{}).Where("id = ?", orderItemId).First(&orderItem)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing order item"})
			return
		}
		c.JSON(http.StatusOK, orderItem)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderItemId := c.Param("order_item_id")
		var updatedOrderItem models.OrderItem
		if err := c.BindJSON(&updatedOrderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Fetch the existing orderItem from the database
		var existingOrderItem models.OrderItem
		err := config.GetDB().Where("id = ?", orderItemId).First(&existingOrderItem).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the existing order item"})
			return
		}

		// Update fields if they are provided in the request
		if updatedOrderItem.Quantity != 0 {
			existingOrderItem.Quantity = updatedOrderItem.Quantity
		}

		if updatedOrderItem.UnitPrice != 0 {
			existingOrderItem.UnitPrice = updatedOrderItem.UnitPrice
		}

		if updatedOrderItem.FoodId != 0 {
			existingOrderItem.FoodId = updatedOrderItem.FoodId
		}
		existingOrderItem.UpdatedAt = time.Now()

		// Save the updated order item back to the database
		err = config.GetDB().Save(&existingOrderItem).Error
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the order item"})
			return
		}
		c.JSON(http.StatusOK, existingOrderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItemPack *OrderItemPack
		var order models.Order
		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		order.OrderDate = time.Now()
		order.TableId = orderItemPack.TableId
		// assign the userId to the order to be created
		userId, userIdExists := c.Get("uid")
		if userIdExists {
			order.UserId = userId.(uint32)
		}
		// if Order exists already we just assign the orderItems to it, otherwise we create a new Order
		orderId := orderItemPack.OrderItems[0].OrderId
		if orderId == 0 {
			orderId = OrderItemOrderCreator(order)
		}
		for i := range orderItemPack.OrderItems {
			orderItemPack.OrderItems[i].OrderId = orderId
			orderItemPack.OrderItems[i].CreatedAt = time.Now()
			orderItemPack.OrderItems[i].UpdatedAt = time.Now()
			var num = toFixed(orderItemPack.OrderItems[i].UnitPrice, 2)
			orderItemPack.OrderItems[i].UnitPrice = num
		}
		config.GetDB().Create(&orderItemPack.OrderItems)
		defer cancel()
		c.JSON(http.StatusOK, fmt.Sprintf("orderItems Created successfully for orderId : %v", orderId))
	}
}
