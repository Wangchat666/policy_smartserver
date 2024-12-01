package government_liaoning

import (
	"policy-search_engine/service"
	"policy-search_engine/service/liaoning/government-liaoning/content"
	"policy-search_engine/service/liaoning/government-liaoning/meta"
)

const name = "government-liaoning" // 辽宁省政策

type ScienceColly struct {
	content *content.ScienceContentColly
	meta    *meta.ScienceMetaColly
}

func (s *ScienceColly) Meta() service.MetaCrawler {
	return (service.MetaCrawler)(s.meta)
}

func (s *ScienceColly) Content() service.ContentCrawler {
	return (service.ContentCrawler)(s.content)
}

func (s *ScienceColly) Register(crawlers *service.Crawlers) {
	s.content = new(content.ScienceContentColly)
	s.meta = new(meta.ScienceMetaColly)

	crawlers.Register(name, s)
}

var _ service.Crawler = (*ScienceColly)(nil)
