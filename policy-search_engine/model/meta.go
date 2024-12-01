package model

import "time"

type Meta struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time `gorm:"index"`
	Date         time.Time `gorm:"type:date"`
	Title        string    `gorm:"title"`
	Url          string    `gorm:"url;unique_index"`
	DepartmentID uint      `gorm:"department_id;index"` // 部门id
	ProvinceID   uint      `gorm:"province_id;index"`   // 省份id
}
