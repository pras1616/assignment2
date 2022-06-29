package models

import "time"

type Items struct {
	Item_id     uint      `form:"Item_id" json:"Item_id" gorm:"primaryKey"`
	Item_code   string    `form:"item_code" json:"item_code" xml:"item_code" binding:"required"`
	Description string    `form:"description" json:"description" xml:"description" binding:"required"`
	Quantity    int64     `form:"quantity" json:"quantity" xml:"quantity" binding:"required"`
	Order_id    uint      `form:"order_id" json:"order_id" xml:"order_id" binding:"required"`
	Created_at  time.Time `form:"created_at" json:"created_at" xml:"created_at"`
	Updated_at  time.Time `form:"updated_at" json:"updated_at" xml:"updated_at"`
}
