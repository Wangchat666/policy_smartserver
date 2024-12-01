package externalSources

import (
	"PolicySearchEngine/service"
	"PolicySearchEngine/service/externalSources/content"
	"PolicySearchEngine/service/externalSources/meta"
)

const name = "external-sources" // 国务院

type ExternalSourcesColly struct {
	content *content.ExternalSourcesContentColly
	meta    *meta.ExternalSourcesMetaColly
}

func (s *ExternalSourcesColly) Meta() service.MetaCrawler {
	return (service.MetaCrawler)(s.meta)
}

func (s *ExternalSourcesColly) Content() service.ContentCrawler {
	return (service.ContentCrawler)(s.content)
}

func (s *ExternalSourcesColly) Register(crawlers *service.Crawlers) {
	s.content = new(content.ExternalSourcesContentColly)
	s.meta = new(meta.ExternalSourcesMetaColly)

	crawlers.Register(name, s)
}

var _ service.Crawler = (*ExternalSourcesColly)(nil)
