package model

// Department 部门表
type Department struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"name"`
}
