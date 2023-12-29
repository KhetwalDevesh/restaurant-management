package controller

import (
	"context"
	"fmt"
	config "github.com/KhetwalDevesh/restaurant-management/database"
	"github.com/KhetwalDevesh/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

var validate = validator.New()

type InvoiceViewFormat struct {
	Id             uint32
	PaymentMethod  models.PaymentMethod
	OrderId        uint32
	PaymentStatus  models.PaymentStatus
	PaymentDue     interface{}
	TableNumber    interface{}
	TotalAmount    float64
	PaymentDueDate time.Time
	OrderDetails   interface{}
}

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var allInvoices []*models.Invoice
		err := config.GetDB().Model(models.Invoice{}).Find(&allInvoices)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "GetInvoices got error"})
		}
		c.JSON(http.StatusOK, allInvoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		invoiceId := c.Param("invoice_id")
		var invoice *models.Invoice
		err := config.GetDB().Model(models.Invoice{}).Where("id = ?", invoiceId).First(&invoice).Error
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing invoice item"})
		}
		var invoiceView InvoiceViewFormat
		allOrderItems, er := ItemsByOrder(invoice.OrderID)
		if er != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing invoice item"})
		}
		_ = config.GetDB().Model(models.OrderItem{}).Select("SUM(quantity * unit_price) as total_amount").Where("order_id = ?", invoice.OrderID).Scan(&invoiceView.TotalAmount).Error
		invoiceView.OrderId = invoice.OrderID
		invoiceView.PaymentDueDate = invoice.PaymentDueDate
		invoiceView.PaymentMethod = invoice.PaymentMethod
		invoiceView.Id = invoice.ID
		invoiceView.PaymentStatus = invoice.PaymentStatus
		invoiceView.TableNumber = allOrderItems[0].Order.Table.TableNumber
		invoiceView.OrderDetails = allOrderItems
		c.JSON(http.StatusOK, invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		isUserAdmin, _ := c.Get("isAdmin")
		if isUserAdmin == false {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You need to be an admin to create an invoice"})
			return
		}
		var invoice *models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var order *models.Order

		if err := config.GetDB().Where("id = ?", invoice.OrderID).First(&order).Error; err != nil {
			msg := fmt.Sprintf("Order was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		status := "pending"
		if invoice.PaymentStatus == "" {
			invoice.PaymentStatus = models.PaymentStatus(status)
		}

		invoice.PaymentDueDate = time.Now().AddDate(0, 0, 1)
		invoice.CreatedAt = time.Now()
		invoice.UpdatedAt = time.Now()

		if err := validate.Struct(invoice); err != nil {
			msg := fmt.Sprintf("Invoice invalidated : %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Save the invoice to the database
		if err := config.GetDB().Create(&invoice).Error; err != nil {
			msg := fmt.Sprintf("Invoice item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, invoice)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var invoice *models.Invoice
		invoiceID := c.Param("invoice_id")

		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find the existing invoice by invoiceID
		var existingInvoice *models.Invoice
		if err := config.GetDB().Where("id = ?", invoiceID).First(&existingInvoice).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
			return
		}
		defer cancel()

		// Update the invoice fields
		if invoice.PaymentMethod != "" {
			existingInvoice.PaymentMethod = invoice.PaymentMethod
		}

		if invoice.PaymentStatus != "" {
			existingInvoice.PaymentStatus = invoice.PaymentStatus
		}

		existingInvoice.UpdatedAt = time.Now()

		if err := validate.Struct(existingInvoice); err != nil {
			msg := fmt.Sprintf("Invoice invalidated : %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Save the updated invoice
		if err := config.GetDB().Save(&existingInvoice).Error; err != nil {
			msg := fmt.Sprintf("Invoice update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, existingInvoice)
	}
}
