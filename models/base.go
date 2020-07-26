package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

func CreateBaseRecord(record interface{}) error {
	return DB().Create(record).Error
}

func DeleteBaseRecord(record interface{}) error {
	return DB().Delete(record).Error
}

func UpdateBaseRecord(record interface{}) error {
	return DB().Model(record).Update(record).Error
}

func QueryBaseRecord(orgModel interface{}, offset, limit uint, filters map[string]string) (uint, []map[string]interface{}, error) {

	filterTmp := make(map[string]interface{})

	for k, v := range filters {
		filterTmp[k] = v
	}

	fmt.Println("filterTmp->", filterTmp)

	// 统计总数
	var total uint = 0
	err := DB().Model(orgModel).Where(filterTmp).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}
	fmt.Println("total->", total)

	// 查询记录
	var ret []map[string]interface{}

	err = DB().Model(orgModel).Where(filterTmp).Offset(offset).Limit(limit).Find(&ret).Error

	if err != nil {
		if gorm.ErrRecordNotFound == err { // 没有记录
			return 0, ret, nil
		}
		return 0, nil, err
	}

	return 0, ret, nil
}
