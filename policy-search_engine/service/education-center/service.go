package education_center

import (
	"PolicySearchEngine/service"
	"PolicySearchEngine/service/education-center/content"
	"PolicySearchEngine/service/education-center/meta"
)

const name = "education-center" // 中央教育部

type EducationColly struct {
	content *content.EducationContentColly
	meta    *meta.EducationMetaColly
}

func (s *EducationColly) Meta() service.MetaCrawler {
	return (service.MetaCrawler)(s.meta)
}

func (s *EducationColly) Content() service.ContentCrawler {
	return (service.ContentCrawler)(s.content)
}

func (s *EducationColly) Register(crawlers *service.Crawlers) {
	s.content = new(content.EducationContentColly)
	s.meta = new(meta.EducationMetaColly)

	crawlers.Register(name, s)
}

var _ service.Crawler = (*EducationColly)(nil)
