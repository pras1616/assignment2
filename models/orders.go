package models

import "time"

type Orders struct {
	Order_id     uint      `form:"order_id" json:"order_id" gorm:"primaryKey"`
	CustomerName string    `form:"customerName" json:"customerName" xml:"customerName" binding:"required"`
	Ordered_at   time.Time `form:"ordered_at" json:"ordered_at" xml:"ordered_at" binding:"required"`
	Created_at   time.Time `form:"created_at" json:"created_at" xml:"created_at"`
	Updated_at   time.Time `form:"updated_at" json:"updated_at" xml:"updated_at"`
}
