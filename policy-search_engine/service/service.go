package service

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"policy-search_engine/config"
	"time"
)

type Crawlers struct {
	Name    string
	Crawler []Crawler
}

func (c *Crawlers) Run() {
	cr := cron.New()                    //创建一个新的定时任务调度器实例。
	for _, crawler := range c.Crawler { //遍历 Crawlers 实例中的所有爬虫实例。

		meta := crawler.Meta()       //获取当前爬虫实例的 Meta 对象。
		content := crawler.Content() //获取当前爬虫实例的 Content 对象。

		// todo 先运行一遍，防止本来就有问题，代码稳定后可删除
		meta.ExecuteWorkflow()    //立即执行 Meta 对象的工作流。
		content.ExecuteWorkflow() //立即执行 Content 对象的工作流。

		spec := config.V.GetString("cron." + c.Name) //从配置文件中获取当前爬虫的定时任务规则。
		_, err := cr.AddFunc(spec, func() {
			fmt.Printf("定时任务运行 time:%s name:%s \n", time.Now(), c.Name)
			meta.ExecuteWorkflow()
			content.ExecuteWorkflow()
		})
		if err != nil {
			fmt.Printf("定时任务添加失败 err: %+v \n", err)
			return
		}
	}
	cr.Start() //启动调度器
}

// Register 新部门加入Crawler
func (c *Crawlers) Register(name string, crawler Crawler) { //参数是一个字符串 name 和一个 Crawler 类型的实例 crawler
	c.Name = name                          //将 c 的 Name 成员变量设置为传入的 name 参数。这个名称用于标识注册的爬虫。
	c.Crawler = append(c.Crawler, crawler) //将传入的 crawler 实例添加到 c 的 Crawler 切片中。
}
