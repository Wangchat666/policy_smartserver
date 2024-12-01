package science_tianjin1

import (
	"PolicySearchEngine/service"
	"PolicySearchEngine/service/tianjin/science-tianjin1/content"
	"PolicySearchEngine/service/tianjin/science-tianjin1/meta"
)

const name = "science-tianjin" // 辽宁科技部

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
