package main

import (
	"PolicySearchEngine/config"
	"PolicySearchEngine/dao/database"
	"PolicySearchEngine/dao/es"
	"PolicySearchEngine/http"
	"PolicySearchEngine/service"
	science_liaoning "PolicySearchEngine/service/liaoning/science-liaoning"
)

func main() {
	// 配置初始化
	config.Init()
	database.Init()
	database.InitTable()
	es.Init()

	var crawler service.Crawlers

	//中央科技部
	//var scienceColly33 science_center.ScienceColly
	//scienceColly33.Register(&crawler)

	//辽宁科学技术厅
	var scienceColly6 science_liaoning.ScienceColly
	scienceColly6.Register(&crawler)

	//天津市科学技术局
	//var scienceColly2 science_tianjin.ScienceColly
	//scienceColly2.Register(&crawler)
	//var scienceColly2_1 science_tianjin1.ScienceColly
	//scienceColly2_1.Register(&crawler)

	//辽宁省政府
	//var PolicyColly6 government_liaoning.ScienceColly
	//PolicyColly6.Register(&crawler)

	//工业信息部
	//var industryInformatizationColly industryInformatization_center.IndustryInformatizationColly
	//industryInformatizationColly.Register(&crawler)

	//中央教育部
	//var educationColly education_center.EducationColly
	//educationColly.Register(&crawler)

	//各省份
	//var stateCouncilColly stateCouncil_center.StateCouncilColly
	//stateCouncilColly.Register(&crawler)

	//var externalSourcesColly externalSources.ExternalSourcesColly
	//externalSourcesColly.Register(&crawler)

	crawler.Run()
	http.Router()
}
