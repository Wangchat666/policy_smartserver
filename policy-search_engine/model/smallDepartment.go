package model

// SmallDepartment 小部门表
type SmallDepartment struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"name"`
}
