package models

import (
	"time"

	"gorm.io/gorm"
)

type Expenses struct {
	gorm.Model
	Date        time.Time
	Sum         int
	Discription *string
	UserID      int
	User        Users
	CategoryID  int
	Category    Categories
}
