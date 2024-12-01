package stateCouncil_center

import (
	"policy-search_engine/service"
	"policy-search_engine/service/stateCouncil-center/content"
	"policy-search_engine/service/stateCouncil-center/meta"
)

const name = "stateCouncil-center" // 国务院

type StateCouncilColly struct {
	content *content.StateCouncilContentColly
	meta    *meta.StateCouncilMetaColly
}

func (s *StateCouncilColly) Meta() service.MetaCrawler {
	return (service.MetaCrawler)(s.meta)
}

func (s *StateCouncilColly) Content() service.ContentCrawler {
	return (service.ContentCrawler)(s.content)
}

func (s *StateCouncilColly) Register(crawlers *service.Crawlers) {
	s.content = new(content.StateCouncilContentColly)
	s.meta = new(meta.StateCouncilMetaColly)

	crawlers.Register(name, s)
}

var _ service.Crawler = (*StateCouncilColly)(nil)
