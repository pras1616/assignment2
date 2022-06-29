package controllers

import (
	"assignment2/models"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Payload struct {
	OrderedAt    string `form:"orderedAt" json:"orderedAt"`
	CustomerName string `form:"customerName" json:"customerName"`
	Items        []struct {
		ItemCode    string `form:"itemCode" json:"itemCode"`
		Description string `form:"description" json:"description"`
		Quantity    int64  `form:"quantity" json:"quantity"`
	} `form:"items" json:"items"`
}

type Response struct {
	Order_id     uint
	OrderedAt    time.Time
	CustomerName string
	Items        []models.Items
}

type PayloadUpdate struct {
	OrderedAt    string `form:"orderedAt" json:"orderedAt"`
	CustomerName string `form:"customerName" json:"customerName"`
	Items        []struct {
		LineItemId  uint   `form:"lineItemId" json:"lineItemId"`
		ItemCode    string `form:"itemCode" json:"itemCode"`
		Description string `form:"description" json:"description"`
		Quantity    int64  `form:"quantity" json:"quantity"`
	} `form:"items" json:"items"`
}

type Controllers interface {
	CreateOrders(c *gin.Context)
	GetAllOrders(c *gin.Context)
	// GetCarsByID(c *gin.Context)
	UpdateOrders(c *gin.Context)
	DeleteOrders(c *gin.Context)
}

type ControllersStruct struct {
	DB_Order *gorm.DB
	DB_Item  *gorm.DB
}

func NewCarsController(db1 *gorm.DB, db2 *gorm.DB) Controllers {
	return &ControllersStruct{
		DB_Order: db1,
		DB_Item:  db2,
	}
}

func (g *ControllersStruct) CreateOrders(c *gin.Context) {
	var payload Payload

	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println("error found: ", err)
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	t, _ := time.Parse(time.RFC3339, payload.OrderedAt)
	reqOrder := models.Orders{
		CustomerName: payload.CustomerName,
		Ordered_at:   t,
		Created_at:   time.Now(),
		Updated_at:   time.Now(),
	}

	err := g.DB_Order.Create(&reqOrder).Error
	if err != nil {
		fmt.Println("error found: ", err)
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
		return
	}

	for _, item := range payload.Items {
		// user.ID // 1,2,3
		reqItem := models.Items{
			Item_code:   item.ItemCode,
			Description: item.Description,
			Quantity:    item.Quantity,
			Order_id:    reqOrder.Order_id,
			Created_at:  time.Now(),
			Updated_at:  time.Now(),
		}
		err := g.DB_Item.Create(&reqItem).Error
		if err != nil {
			fmt.Println("error found: ", err)
			c.JSON(500, gin.H{
				"message": "internal server error",
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"message": "succesful to create new data",
	})
}

func (g *ControllersStruct) GetAllOrders(c *gin.Context) {
	resultOrder := []models.Orders{}
	resultItem := []models.Items{}

	err := g.DB_Order.Find(&resultOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("error found: ", err)
		c.JSON(404, gin.H{
			"message": "data empty",
		})
		return
	}
	if err != nil {
		fmt.Println("error found: ", err)
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
		return
	}

	var sendResult []Response
	for _, item := range resultOrder {
		params := `order_id = ` + strconv.FormatUint(uint64(item.Order_id), 10)
		errorItem := g.DB_Item.Where(params).Find(&resultItem).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("error found: ", errorItem)
			c.JSON(404, gin.H{
				"message": "data empty",
			})
			return
		}
		if errorItem != nil {
			fmt.Println("error found: ", errorItem)
			c.JSON(500, gin.H{
				"message": "internal server error",
			})
			return
		}

		sendResult = append(sendResult, Response{
			Order_id:     item.Order_id,
			OrderedAt:    item.Ordered_at,
			CustomerName: item.CustomerName,
			Items:        resultItem,
		},
		)

	}

	c.JSON(200, gin.H{
		"message": sendResult,
	})
}

func (g *ControllersStruct) UpdateOrders(c *gin.Context) {
	var payloadUpdate PayloadUpdate
	idStr, _ := c.Params.Get("orderId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Println("error found: ", err)
		c.JSON(404, gin.H{
			"message": "Not Found",
		})
		return
	}
	req := payloadUpdate

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("error found: ", err)
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}

	t, _ := time.Parse(time.RFC3339, req.OrderedAt)
	reqOrder := models.Orders{
		CustomerName: req.CustomerName,
		Ordered_at:   t,
		Updated_at:   time.Now(),
	}

	err = g.DB_Order.Model(&reqOrder).Where("order_id = ?", id).Updates(models.Orders{
		CustomerName: reqOrder.CustomerName,
		Ordered_at:   reqOrder.Ordered_at,
		Updated_at:   reqOrder.Updated_at,
	}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("error found: ", err)
		c.JSON(404, gin.H{
			"message": "data not found",
		})
		return
	}
	if err != nil {
		fmt.Println("error found: ", err)
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
		return
	}

	fmt.Println(req.Items, "<<<<<<<<<<<<<<<<<<<")
	for _, item := range req.Items {
		reqItem := models.Items{
			Item_id:     item.LineItemId,
			Item_code:   item.ItemCode,
			Description: item.Description,
			Quantity:    item.Quantity,
			Updated_at:  time.Now(),
		}
		err = g.DB_Item.Model(&reqItem).Where("item_id = ?", reqItem.Item_id).Updates(models.Items{
			Item_code:   reqItem.Item_code,
			Description: item.Description,
			Quantity:    item.Quantity,
			Updated_at:  reqItem.Updated_at,
		}).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("error found: ", err)
			c.JSON(404, gin.H{
				"message": "data not found",
			})
			return
		}
		if err != nil {
			fmt.Println("error found: ", err)
			c.JSON(500, gin.H{
				"message": "internal server error",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "succesful to update data",
	})
}

func (g *ControllersStruct) DeleteOrders(c *gin.Context) {
	idStr, _ := c.Params.Get("orderId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Println("error found: ", err)
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	reqOrder := models.Orders{}
	reqItem := models.Items{}

	err = g.DB_Order.Delete(&reqOrder, "order_id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("error found: ", err)
		c.JSON(404, gin.H{
			"message": "data not found",
		})
		return
	}
	if err != nil {
		fmt.Println("error found: ", err)
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
		return
	}

	err = g.DB_Item.Delete(&reqItem, "order_id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("error found: ", err)
		c.JSON(404, gin.H{
			"message": "data not found",
		})
		return
	}
	if err != nil {
		fmt.Println("error found: ", err)
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "succesful to delete data",
	})
}
