package model

// SmallDepartmentMap metaId与小部门的映射关系表
type SmallDepartmentMap struct {
	ID     uint `gorm:"primarykey"`
	MetaID uint `gorm:"not null"`
	// 对应SmallDepartment表id
	SmallDepartmentID uint `gorm:"not null"`
}
