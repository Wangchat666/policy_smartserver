package science_center

import (
	"policy-search_engine/service"
	"policy-search_engine/service/science-center/content"
	"policy-search_engine/service/science-center/meta"
)

const name = "science-center" // 中央科技部

type ScienceColly struct {
	content *content.ScienceContentColly
	meta    *meta.ScienceMetaColly
}

func (s *ScienceColly) Meta() service.MetaCrawler {
	return (service.MetaCrawler)(s.meta) //将 s.meta 成员变量从其原始类型转换为 service.MetaCrawler 接口类型，并返回这个接口实例。
} //使得 ScienceColly 实例可以被当作 service.MetaCrawler 接口类型来使用。

func (s *ScienceColly) Content() service.ContentCrawler {
	return (service.ContentCrawler)(s.content)
}

func (s *ScienceColly) Register(crawlers *service.Crawlers) {
	s.content = new(content.ScienceContentColly) //初始化 s 的 content 成员变量为一个新的 content.ScienceContentColly 实例。
	s.meta = new(meta.ScienceMetaColly)          //初始化 s 的 meta 成员变量为一个新的 meta.ScienceMetaColly 实例。

	crawlers.Register(name, s)
}

var _ service.Crawler = (*ScienceColly)(nil)
