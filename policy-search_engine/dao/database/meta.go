package database

import (
	"PolicySearchEngine/model"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type MetaDal struct{ Db *gorm.DB }

// InsertMeta 添加meta
func (m *MetaDal) InsertMeta(date time.Time, title string, url string, departmentID uint, provinceID uint) uint {
	meta := model.Meta{
		Date:         date,
		Title:        title,
		Url:          url,
		DepartmentID: departmentID,
		ProvinceID:   provinceID,
	}
	// 存在则忽略，不存在则插入(根据url判断）
	result := m.Db.Where(model.Meta{Url: meta.Url}).FirstOrCreate(&meta) //确保了在数据库中存在一个 Url 字段等于 meta.Url 的记录，无论是通过查找现有记录还是创建新记录来实现。
	if result.Error != nil {
		fmt.Printf("InsertMeta... %s, %v", date.String(), meta)
		log.Fatal(result.Error) //记录错误信息并终止程序。log.Fatal 会输出错误信息并调用 os.Exit(1)，从而终止程序的执行。
	}
	return meta.ID
}

func (m *MetaDal) UpdateMetaTitle(title string, url string) {
	meta := model.Meta{
		Title: title,
		Url:   url,
	}
	result := m.Db.Where(model.Meta{Url: meta.Url}).Updates(&meta) //Updates(&meta) 则将 meta 对象中的非零字段（Title字段）更新到数据库中
	if result.Error != nil {
		fmt.Printf("UpdateMetaTitle... %v", meta)
		log.Fatal(result.Error)
	}
}

func (m *MetaDal) GetAllMeta(departmentID, provinceID uint) *[]model.Meta {
	var metaList []model.Meta
	result := m.Db.Where(model.Meta{
		DepartmentID: departmentID,
		ProvinceID:   provinceID,
	}).Find(&metaList) //m.Db 是数据库连接实例，Where 方法用于指定查询条件是 DepartmentID 和 ProvinceID 等于传入的参数值。Find 方法将查询结果存储到 metaList 中。
	if result.Error != nil {
		fmt.Printf("读取数据失败: %v\n", result.Error)
		return nil
	}
	return &metaList
} //etAllMeta方法的作用是根据部门ID和省份ID从数据库中检索元数据，并返回这些数据的列表。如果检索过程中出现错误，它会返回nil。

func (m *MetaDal) GetAllMetaByIDs(provinceID uint, id uint) *[]model.Meta {
	var metaList []model.Meta
	result := m.Db.Where("province_id = ? AND id > ?", provinceID, id).Find(&metaList) //即查找 province_id 等于 provinceID 且 id 大于 id 的记录。Find(&metaList) 将查询结果存储到 metaList 切片中。
	if result.Error != nil {
		fmt.Printf("读取数据失败: %v\n", result.Error)
		return nil
	}
	return &metaList
}

func (m *MetaDal) GetMetaByUrl(url string) *model.Meta {
	var meta model.Meta
	result := m.Db.Where(model.Meta{
		Url: url,
	}).First(&meta) //First(&meta) 将查询结果存储到 meta 对象中，并且只返回第一个匹配的记录。
	if result.Error != nil {
		fmt.Printf("读取数据失败: %v\n", result.Error)
		return nil
	}
	return &meta
}
