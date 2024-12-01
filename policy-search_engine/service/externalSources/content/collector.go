package content

import (
	"PolicySearchEngine/service"
	"PolicySearchEngine/utils"
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strings"
)

func (s *ExternalSourcesContentColly) getRules() []*service.Rule {
	return []*service.Rule{
		s.zhengceCollector(),
		s.niaCollector(),
		s.zfxxgkmlCollector(),
		s.chinataxCollector(),
		s.ccgpCollector(),
		s.cdeCollector(),
		s.nmpaCollector(),
		s.statsCollector(),
		s.mvaCollector(),
		s.nhcCollector(),
		s.natcmCollector(),
		s.mofCollector(),
		s.motCollector(),
		s.spbCollector(),
		s.samrCollector(),
		s.cncaCollector(),
		s.saacCollector(),
		s.ncacCollector(),
		s.nhsaCollector(),
		s.caeCollector(),
		s.customsCollector(),
		s.forestryCollector(),
		s.lswzCollector(),
		s.meeCollector(),
		s.moagkCollector(),
		s.moagovCollector(),
		s.cnipaCollector(),
		s.miitCollector(),
		s.mcaCollector(),
		s.gjxfjCollector(),
		s.zfxxgkCollector(),
		s.wwwneaCollector(),
		s.mwrzwgkCollector(),
		s.mwrzwCollector(),
		s.mohrssCollector(),
		s.mofcomCollector(),
		s.ndrcCollector(),
		s.satcmCollector(),
		s.sipoCollector(),
	}
}

func (s *ExternalSourcesContentColly) updateTitle(e *colly.HTMLElement) {
	title := utils.TidyString(e.Text)
	s.metaDal.UpdateMetaTitle(title, e.Request.URL.String())
}

func (s *ExternalSourcesContentColly) updateContent(e *colly.HTMLElement) {
	var text []byte
	e.ForEach("*", func(_ int, child *colly.HTMLElement) {
		label := strings.ToLower(child.Name)
		if label == "style" || label == "table" || label == "script" {
			return
		}
		text = append(text, []byte(child.Text)...)
	})
	s.contentDal.InsertContent(e.Request.URL.String(), string(text))

	meta := s.metaDal.GetMetaByUrl(e.Request.URL.String())
	if meta == nil {
		meta = s.metaDal.GetMetaByUrl(e.Request.Headers.Get("Referer"))
	}
	if meta != nil {
		sdIDs := s.dMapDal.GetDepartmentIDsByMetaID(meta.ID)
		fmt.Println(sdIDs)
		//es.IndexDoc(meta.Date, meta.DepartmentID, meta.ProvinceID, meta.Title, meta.Url, string(text), sdIDs)
	} else {
		fmt.Printf("URL: %v, ID: %v, meta:%+v", meta.Url, meta.ID, meta)
		fmt.Println("meta未查询到！！")
	}

}

func (s *ExternalSourcesContentColly) zhengceCollector() *service.Rule {

	rule1 := regexp.MustCompile("https?://www\\.gov\\.cn/zhengce/.*\\.html?")
	rule2 := regexp.MustCompile("https?://www\\.gov\\.cn/xinwen/.*\\.html?")

	combinedRule := regexp.MustCompile(fmt.Sprintf(
		"(%s|%s)",
		rule1.String(),
		rule2.String(),
	))

	hfContent := &service.HtmlFunc{
		QuerySelect: "#UCAP-CONTENT",
		F:           s.updateContent,
	}

	return service.NormalRule(combinedRule, hfContent)
}

func (s *ExternalSourcesContentColly) niaCollector() *service.Rule {
	//https://www.nia.gov.cn/n741440/n741542/c1403794/content.html
	rule := regexp.MustCompile("https?://www\\.nia\\.gov\\.cn/.*content\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: "#content",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) zfxxgkmlCollector() *service.Rule {
	//http://zwgk.mct.gov.cn/zfxxgkml/kjjy/202105/t20210528_924818.html
	//http://zwgk.mct.gov.cn/zfxxgkml/fwzwhyc/202106/t20210610_925136.html
	rule := regexp.MustCompile("https?://zwgk\\.mct\\.gov\\.cn/zfxxgkml.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: ".gsj_htmlcon_bot",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) chinataxCollector() *service.Rule {

	rule := regexp.MustCompile("https?://www\\.chinatax\\.gov\\.cn/.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: "#fontzoom",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) ccgpCollector() *service.Rule {

	//http://www.ccgp.gov.cn/cggg/zygg/gkzb/202211/t20221123_19072939.htm
	rule := regexp.MustCompile("https?://www\\.ccgp\\.gov\\.cn/.*\\.html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".vF_detail_content",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) cdeCollector() *service.Rule {

	//https://www.cde.org.cn/main/news/viewInfoCommon/14aac16a4fc5b5841bc2529988a611cc
	rule := regexp.MustCompile("https?://www\\.cde\\.gov\\.cn/main/news/.*")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".new_detail_content",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) nmpaCollector() *service.Rule {

	//https://www.nmpa.gov.cn/xxgk/ggtg/hzhpggtg/hzhpchjgg/hzhpcjgjj/20230927092112111.html
	rule := regexp.MustCompile("https?://www\\.nmpa\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".text",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) statsCollector() *service.Rule {

	//https://www.stats.gov.cn/xw/tjxw/tzgg/202404/t20240409_1954351.html
	rule := regexp.MustCompile("https?://www\\.stats\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".TRS_PreAppend",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) mvaCollector() *service.Rule {

	//http://www.mva.gov.cn/sy/xx/tzgg/202012/t20201230_44016.html
	rule := regexp.MustCompile("https?://www\\.mva\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: "#div_zhengwen",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) nhcCollector() *service.Rule {

	//http://www.nhc.gov.cn/lljks/tggg/202005/3726889dbf7f4cd0abd9a92105ae53ff.shtml
	rule := regexp.MustCompile("https?://www\\.nhc\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: "#xw_box",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) natcmCollector() *service.Rule {

	//http://www.natcm.gov.cn/bangongshi/gongzuodongtai/2022-10-13/27881.html
	rule := regexp.MustCompile("https?://www\\.natcm\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".zw",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) mofCollector() *service.Rule {

	//http://gks.mof.gov.cn/gongzuodongtai/202308/t20230802_3899924.htm
	rule := regexp.MustCompile("https?://.*\\.mof\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".TRS_Editor",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) motCollector() *service.Rule {

	//https://xxgk.mot.gov.cn/2020/jigou/aqyzljlglj/202107/t20210702_3611040.html
	rule := regexp.MustCompile("https?://xxgk\\.mot\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: "#Zoom",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) spbCollector() *service.Rule {

	//https://www.spb.gov.cn/gjyzj/c100015/c100016/202307/a95b0ebbb5a643e396bdadf262955930.shtml
	rule := regexp.MustCompile("https?://www\\.spb\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".detail-news",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) samrCollector() *service.Rule {

	//https://www.samr.gov.cn/zw/zfxxgk/fdzdgknr/zljds/art/2024/art_ecc33eb3ba634e3aa24930878508f12c.html
	rule := regexp.MustCompile("https?://www\\.samr\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".Three_xilan_07",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) cncaCollector() *service.Rule {

	//https://www.cnca.gov.cn/zwxx/gg/2023/art/2023/art_ef3bb2fec17b44929aae94f807bbf2cf.html
	rule := regexp.MustCompile("https?://www\\.cnca\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".detail_messge",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) saacCollector() *service.Rule {

	//https://www.saac.gov.cn/daj/tzgg/202204/d24b144d1436487a8d9a163322b1e2a5.shtml
	rule := regexp.MustCompile("https?://www\\.saac\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".pages_content",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) ncacCollector() *service.Rule {

	//https://www.ncac.gov.cn/chinacopyright/contents/12547/358243.shtml
	rule := regexp.MustCompile("https?://www\\.ncac\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".m3nEditor",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) nhsaCollector() *service.Rule {

	//http://www.nhsa.gov.cn/art/2024/4/7/art_109_12313.html
	rule := regexp.MustCompile("https?://www\\.nhsa\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: "#zoom",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) caeCollector() *service.Rule {

	//https://www.cae.cn/cae/html/main/col4/2020-11/02/20201102184321555482855_1.html
	rule := regexp.MustCompile("https?://www\\.cae\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: "#zoom",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) customsCollector() *service.Rule {

	//http://www.customs.gov.cn/customs/302249/2480148/5711258/index.html
	rule := regexp.MustCompile("https?://www\\.customs\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: "#easysiteText",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) forestryCollector() *service.Rule {

	//https://www.forestry.gov.cn/c/www/gsgg/511871.jhtml
	rule := regexp.MustCompile("https?://www\\.forestry\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: "#zoomit",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) lswzCollector() *service.Rule {

	//http://www.lswz.gov.cn/html/ywpd/rsrc/2022-09/01/content_272034.shtml
	rule := regexp.MustCompile("https?://www\\.lswz\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".pub-det-content",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) meeCollector() *service.Rule {

	//https://www.mee.gov.cn/ywgz/zcghtjdd/ghxx/202206/t20220628_987021.shtml
	rule := regexp.MustCompile("https?://www\\.mee\\.gov\\.cn/.*html?")
	hfContent := &service.HtmlFunc{
		QuerySelect: ".TRS_Editor",
		F:           s.updateContent,
	}
	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) moagkCollector() *service.Rule {
	//http://www.moa.gov.cn/gk/tzgg_1/tz/202306/t20230623_6430795.htm
	rule := regexp.MustCompile("https?://www\\.moa\\.gov\\.cn/gk/.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: ".TRS_Editor",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) moagovCollector() *service.Rule {
	//http://www.moa.gov.cn/govpublic/CJB/202106/t20210625_6370305.htm
	rule := regexp.MustCompile("https?://www\\.moa\\.gov\\.cn/govpublic/.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: ".gsj_htmlcon_bot",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) cnipaCollector() *service.Rule {
	//https://www.cnipa.gov.cn/art/2020/11/6/art_75_154667.html
	rule := regexp.MustCompile("https?://www\\.cnipa\\.gov\\.cn/.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: ".article",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) miitCollector() *service.Rule {
	//https://www.miit.gov.cn/zwgk/zcwj/wjfb/tz/art/2023/art_a6abdf55cabe446ea447d935e7622366.html?xxgkhide=1
	rule := regexp.MustCompile("https?://www\\.miit\\.gov\\.cn/.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: "#con_con",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) mcaCollector() *service.Rule {
	//https://www.mca.gov.cn/n152/n165/c1662004999979993926/content.html
	rule := regexp.MustCompile("https?://www\\.mca\\.gov\\.cn/.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: "#zoom",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) gjxfjCollector() *service.Rule {
	//https://www.gjxfj.gov.cn/gjxfj/news/ggb/webinfo/2020/06/1590610835303025.htm
	rule := regexp.MustCompile("https?://www\\.gjxfj\\.gov\\.cn/.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: ".ejxxgk_xq_con",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) zfxxgkCollector() *service.Rule {
	//http://zfxxgk.nea.gov.cn/2024-03/09/c_1310768943.htm
	rule := regexp.MustCompile("https?://zfxxgk\\.nea\\.gov\\.cn/.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: ".article-content",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) wwwneaCollector() *service.Rule {
	//http://www.nea.gov.cn/2024-03/18/c_1310768057.htm
	rule := regexp.MustCompile("https?://www\\.nea\\.gov\\.cn/.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: ".article-content",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) mwrzwgkCollector() *service.Rule {
	//http://www.mwr.gov.cn/zwgk/gknr/202302/t20230220_1646227.html
	rule := regexp.MustCompile("https?://www\\.mwr\\.gov\\.cn/zwgk.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: ".gknb_content",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) mwrzwCollector() *service.Rule {
	//http://www.mwr.gov.cn/zw/tzgg/tzgs/202402/t20240226_1704337.html
	rule := regexp.MustCompile("https?://www\\.mwr\\.gov\\.cn/zw.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: ".TRS_Editor",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) mohrssCollector() *service.Rule {
	//http://www.mohrss.gov.cn/SYrlzyhshbzb/zwgk/gggs/tg/202404/t20240415_516839.html
	rule := regexp.MustCompile("https?://www\\.mohrss\\.gov\\.cn/.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: ".TRS_Editor",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) mofcomCollector() *service.Rule {
	//http://www.mofcom.gov.cn/article/h/redht/202110/20211003205975.shtml
	rule := regexp.MustCompile("https?://www\\.mofcom\\.gov\\.cn/.*html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: ".f-cb",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) ndrcCollector() *service.Rule {
	//https://www.ndrc.gov.cn/xxgk/zcfb/tz/202404/t20240416_1365679.html
	rule := regexp.MustCompile("https?://www\\.ndrc\\.gov\\.cn/.*\\.html?")

	hfContent := &service.HtmlFunc{
		QuerySelect: ".article",
		F:           s.updateContent,
	}

	return service.NormalRule(rule, hfContent)
}

func (s *ExternalSourcesContentColly) satcmCollector() *service.Rule {
	//http://www.satcm.gov.cn/renjiaosi/gongzuodongtai/2022-12-14/28482.html
	rule1 := regexp.MustCompile("https?://www\\.satcm\\.gov\\.cn/.*\\.html?")
	rule2 := regexp.MustCompile("https?://yzs\\.satcm\\.gov\\.cn/.*\\.html?")

	combinedRule := regexp.MustCompile(fmt.Sprintf(
		"(%s|%s)",
		rule1.String(),
		rule2.String(),
	))

	return service.NormalRule(combinedRule)
}

func (s *ExternalSourcesContentColly) sipoCollector() *service.Rule {
	//http://www.sipo.gov.cn/gztz/1151437.htm
	rule := regexp.MustCompile("https?://www\\.sipo\\.gov\\.cn/.*")

	return service.NormalRule(rule)
}
