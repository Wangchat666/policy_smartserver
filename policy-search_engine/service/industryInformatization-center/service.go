package industryInformatization_center

import (
	"PolicySearchEngine/service"
	"PolicySearchEngine/service/industryInformatization-center/content"
	"PolicySearchEngine/service/industryInformatization-center/meta"
)

const name = "industryInformatization-center" // 中央工信部

type IndustryInformatizationColly struct {
	content *content.IndustryInformatizationContentColly
	meta    *meta.IndustryInformatizationMetaColly
}

func (s *IndustryInformatizationColly) Meta() service.MetaCrawler {
	return (service.MetaCrawler)(s.meta)
}

func (s *IndustryInformatizationColly) Content() service.ContentCrawler {
	return (service.ContentCrawler)(s.content)
}

func (s *IndustryInformatizationColly) Register(crawlers *service.Crawlers) {
	s.content = new(content.IndustryInformatizationContentColly)
	s.meta = new(meta.IndustryInformatizationMetaColly)

	crawlers.Register(name, s)
}

var _ service.Crawler = (*IndustryInformatizationColly)(nil)
