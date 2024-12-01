package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"policy-search_engine/model"
	"strings"
)

type ContentDal struct{ Db *gorm.DB }

func (c *ContentDal) InsertContent(url string, article string) {

	var meta model.Meta
	if err := c.Db.Where("url = ?", url).First(&meta).Error; err != nil { //查找 url 对应的元数据记录。如果查询失败，进入错误处理逻辑。
		if errors.Is(err, gorm.ErrRecordNotFound) { //如果错误是记录未找到，尝试将 url 中的 https:// 替换为 http://，再次查询。
			url = strings.Replace(url, "https://", "http://", 1)
			err = c.Db.Where("url = ?", url).First(&meta).Error //再次查询
		}
		if err != nil {
			// 处理查找失败的情况，例如返回错误或者进行其他逻辑
			fmt.Printf("查找Meta记录失败 url:%s err:%+v\n", url, err)
			return
		}
	}

	content := model.Content{
		MetaID:  meta.ID,
		Article: article,
	} //创建一个 model.Content 类型的变量 content，设置其 MetaID 为找到的元数据记录的 ID，Article 为传入的文章内容。

	result := c.Db.Where(model.Content{MetaID: meta.ID}).
		Assign(model.Content{Article: article}).
		FirstOrCreate(&content) //查找或创建内容记录。如果找到匹配的记录，则更新其 Article 字段；如果没有找到，则创建新记录。
	if result.Error != nil {
		fmt.Printf("插入Content记录失败 err:%+v", result.Error)
	}

} //根据给定的 URL 查找对应的元数据记录，然后将文章内容插入到内容表中，并与元数据记录关联。如果 URL 对应的元数据记录不存在，会尝试替换 URL 中的协议再次查询。

func (c *ContentDal) GetContentByMetaID(id uint) *model.Content {
	var content model.Content
	result := c.Db.Where(model.Content{MetaID: id}).First(&content)
	if result.Error != nil {
		fmt.Printf("读取文章内容失败: %v\n", result.Error)
		return nil
	}
	return &content
}
