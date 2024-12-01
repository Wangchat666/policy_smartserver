package model

// Province 省份表
type Province struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"name"`
}
