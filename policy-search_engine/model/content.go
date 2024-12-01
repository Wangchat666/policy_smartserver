package model

import "time"

type Content struct {
	ID        uint   `gorm:"primarykey"`
	MetaID    uint   `gorm:"not null;unique_index"`
	Article   string `gorm:"mediumtext"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"index"`
}

//`gorm:""`表示位标记，
//`gorm:"primarykey"`表示主键，`gorm:"not null;unique_index"`表示非空且唯一索引，
