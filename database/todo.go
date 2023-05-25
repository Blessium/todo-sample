package database

import (
	"time"
)

type Todo struct {
	ID           uint `gorm:"column:id;primaryKey"`
    Title        string `gorm:"column:title"`
    Notes        string `gorm:"column:notes"`
    CreationDate time.Time `gorm:"column:creation_date"`
    DueDate      time.Time `gorm:"column:due_date"`
    Completed    bool `gorm:"column:completed"`
}
