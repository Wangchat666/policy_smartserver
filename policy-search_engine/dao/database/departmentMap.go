package database

import (
	"PolicySearchEngine/model"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SmallDepartmentMapDal struct{ Db *gorm.DB }

func (d *SmallDepartmentMapDal) InsertDID(metaID uint, sdID uint) {

	// 查询现有记录:
	var dMap model.SmallDepartmentMap
	err := d.Db.Where("meta_id = ? and small_department_id = ?", metaID, sdID).First(&dMap).Error //它查询数据库中是否存在一个meta_id和small_department_id与给定的metaID和sdID匹配的SmallDepartmentMap记录。
	//检查查询错误:
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("查找dMap记录失败 metaID:%d, dID:%d, err:%+v\n", metaID, sdID, err)
		return
	}
	//检查记录是否已存在:
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	// 创建新记录:
	dMap = model.SmallDepartmentMap{
		MetaID:            metaID,
		SmallDepartmentID: sdID,
	}

	//插入新记录
	result := d.Db.Create(&dMap)
	//检查插入错误:
	if result.Error != nil {
		fmt.Printf("插入dMap记录失败 err:%+v", result.Error)
	}
} //这个方法主要用于向数据库中插入一个新的SmallDepartmentMap记录，但如果该记录已经存在（基于meta_id和small_department_id的组合），则不会进行任何操作。如果在查询或插入过程中发生错误（除了记录已存在的情况），则会打印错误。但请注意，此方法没有返回任何错误值，因此调用者无法知道是否成功插入了记录。

func (d *SmallDepartmentMapDal) GetDepartmentIDsByMetaID(id uint) (sdIDs []uint) {
	//这个方法接受一个uint类型的参数id，代表MetaID，并返回一个uint类型的切片sdIDs，该切片包含所有关联的SmallDepartmentID。
	var DepartmentMaps []model.SmallDepartmentMap                                    //用于存储从数据库中查询到的SmallDepartmentMap记录。
	result := d.Db.Where(model.SmallDepartmentMap{MetaID: id}).Find(&DepartmentMaps) //MetaID字段等于传入的id参数。然后调用Find方法来查询数据库，并将结果存储在DepartmentMaps切片中。
	if result.Error != nil {
		fmt.Printf("读取DepartmentMap失败: %v\n", result.Error)
		return sdIDs //出现错误，方法会返回一个空的sdIDs切片。
	}
	for _, departmentMap := range DepartmentMaps {
		sdIDs = append(sdIDs, departmentMap.SmallDepartmentID) //将每个SmallDepartmentMap记录的SmallDepartmentID字段添加到sdIDs切片中。
	}
	return sdIDs
}
