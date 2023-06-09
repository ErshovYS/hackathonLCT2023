package string_generator

import (
	"strconv"
)

func GenerateReport(strVal map[string]string, val map[string]float64) ([]byte, error) {
	r := NewRequestPdf("")

	//html template path
	templatePath := "templates/sample.html"

	pageWidth := 932
	data := make(map[string]string)

	data["PageHeight"] = "1313"
	data["PageHeight7"] = "1312"
	data["PageHeight8"] = "1309"
	data["PageWidth"] = strconv.FormatInt(int64(pageWidth), 10)
	FontSizeText1 := 14
	ImagePosX := 120

	data["FontSizeText1"] = getString(FontSizeText1)
	data["ImageLineWidth"] = getString(pageWidth - ImagePosX*2)

	//страница 1

	textMarkup1 := initTextMarkup(220)

	data["PageLeftPos1"] = "125"
	data["PageTitleLeftPos1"] = "231"
	data["PagePos10"] = getNextRow(textMarkup1)
	data["PagePos11"] = getNextRow(textMarkup1)
	data["PagePos12"] = getNextParagraph(textMarkup1)

	for i := 3; i < 14; i++ {
		data["PagePos1"+getString(i)] = getNextRow(textMarkup1)
	}

	//страница 2
	numLinePage2 := 0
	data["LeftColumnPosPage20"] = "140"
	data["LeftColumnPosPage21"] = "105"

	data["RightColumnPosPage2"] = "480"

	Organization := parseString(strVal["Industry"])

	pageMarkUp20 := initPageMarkup(200, 5, 580)

	data["IndustryColumn1"] = getBodyNextPositionRow(pageMarkUp20)
	data["IndustryText1"] = Organization[0]

	for i := 1; i < len(Organization); i++ {
		if len(Organization[i]) > 0 {
			data["IndustryColumn"+getString(i+1)] = getBodyNextPosition(pageMarkUp20)
			data["IndustryText"+getString(i+1)] = Organization[i]
		}
	}

	data["IndustryPos1"] = getBodyNextLeftPositionRow(pageMarkUp20)
	data["Industry"] = "ОТРАСЛЬ:"

	addImageLine(&data, pageMarkUp20, "2", &numLinePage2, FontSizeText1, ImagePosX, true)

	data["TypeOfOrganizationColumn1"] = getBodyNextPositionRow(pageMarkUp20)
	data["TypeOfOrganizationPos1"] = getBodyNextLeftPositionRow(pageMarkUp20)
	data["TypeOfOrganization"] = "ТИП ОРГАНИЗАЦИИ:"
	data["TypeOfOrganizationText1"] = strVal["Organization"]

	addImageLine(&data, pageMarkUp20, "2", &numLinePage2, FontSizeText1, ImagePosX, true)

	data["NumberOfStaffPos1"] = getBodyNextPositionRow(pageMarkUp20)
	data["NumberOfStaffPos2"] = getBodyNextPosition(pageMarkUp20)
	data["NumberOfStaffColumn1"] = getBodyNextLeftPositionRow(pageMarkUp20)
	data["NumberOfStaff1"] = "КОЛИЧЕСТВО"
	data["NumberOfStaff2"] = "СОТРУДНИКОВ:"
	data["NumberOfStaffText1"] = strVal["WorkersCount"]

	addImageLine(&data, pageMarkUp20, "2", &numLinePage2, FontSizeText1, ImagePosX, true)

	data["LocationPos1"] = getBodyNextPositionRow(pageMarkUp20)
	data["LocationPos2"] = getBodyNextPosition(pageMarkUp20)
	data["LocationPos3"] = getBodyNextPosition(pageMarkUp20)
	data["LocationColumn1"] = getBodyNextLeftPositionRow(pageMarkUp20)
	data["Location1"] = "РАЙОН"
	data["Location2"] = "РАСПОЛОЖЕНИЯ"
	data["Location3"] = "ПРОИЗВОДСТВА:"
	data["LocationText1"] = strVal["District"]

	addImageLine(&data, pageMarkUp20, "2", &numLinePage2, FontSizeText1, ImagePosX, true)

	pageMarkUp21 := initPageMarkup(pageMarkUp20.BodyCurrent, 2, pageMarkUp20.BodyCurrent+150)

	data["ResultPos1"] = getBodyNextPositionRow(pageMarkUp21)
	getBodyNextLeftPositionRow(pageMarkUp21)
	data["ResultText1"] = "ИТОГОВЫЕ ЗНАЧЕНИЯ ВОЗМОЖНЫХ ЗАТРАТ"

	data["TotalPos1"] = getBodyNextPositionRow(pageMarkUp21)
	data["TotalPos2"] = getBodyNextPosition(pageMarkUp21)
	data["TotalColumn1"] = getBodyNextLeftPositionRow(pageMarkUp21)
	data["TotalText1"] = "ИТОГО ВОЗМОЖНЫХ"
	data["TotalText2"] = "РАСХОДОВ"
	data["TotalFrom"] = floatToStringFrom(val["TotalFrom"], val["TotalTo"])
	data["TotalTo"] = floatToString(val["TotalTo"])
	data["TotalUnits"] = getUnits(val["TotalTo"])

	addImageLine(&data, pageMarkUp21, "2", &numLinePage2, FontSizeText1, ImagePosX, true)

	pageMarkUp22 := initPageMarkup(pageMarkUp21.BodyCurrent, 5, 1200)

	addMapDataText1(&data, pageMarkUp22, "Detail", "ДЕТАЛИ РАСХОДОВ")
	getBodyNextLeftPositionRow(pageMarkUp22)

	addMapData1FromTo(&data, &val, pageMarkUp22, "Personal", "ПЕРСОНАЛ")
	addImageLine(&data, pageMarkUp22, "2", &numLinePage2, FontSizeText1, ImagePosX, true)

	addMapData1FromTo(&data, &val, pageMarkUp22, "Estate", "ИМУЩЕСТВО ОРГАНИЗАЦИИ")
	addImageLine(&data, pageMarkUp22, "2", &numLinePage2, FontSizeText1, ImagePosX, true)

	addMapData1FromTo(&data, &val, pageMarkUp22, "Tax", "НАЛОГИ")
	addImageLine(&data, pageMarkUp22, "2", &numLinePage2, FontSizeText1, ImagePosX, true)

	addMapData1FromTo(&data, &val, pageMarkUp22, "Service", "УСЛУГИ")

	//страница 3
	numLinePage3 := 0
	textMarkup3 := initTextMarkup(220)

	data["PagePos30"] = getNextRow(textMarkup3)
	data["PagePos31"] = getNextRow(textMarkup3)
	data["PagePos32"] = getNextRow(textMarkup3)
	data["PagePos33"] = getNextRow(textMarkup3)
	data["PagePos34"] = getNextRow(textMarkup3)
	data["PagePos35"] = getNextRow(textMarkup3)
	data["PagePos36"] = getNextParagraph(textMarkup3)
	data["PagePos37"] = getNextRow(textMarkup3)
	data["PagePos38"] = getNextRow(textMarkup3)

	pageMarkUp3 := initPageMarkup(textMarkup3.CurPos, 4, 1100)

	addMapData2Val(&data, &val, pageMarkUp3, "PersonalCount",
		[]string{"ПЛАНИРУЕМАЯ ЧИСЛЕННОСТЬ", "ПЕРСОНАЛА"})
	data["PersonalCount"] = strconv.FormatInt(int64(val["PersonalCount"]), 10)

	addImageLine(&data, pageMarkUp3, "3", &numLinePage3, FontSizeText1, ImagePosX, true)

	addMapData2FromTo(&data, &val, pageMarkUp3, "PersonalSalary",
		[]string{"ВЫПЛАТА ЗАРПЛАТЫ", "ПЕРСОНАЛА"})

	addImageLine(&data, pageMarkUp3, "3", &numLinePage3, FontSizeText1, ImagePosX, true)

	addMapData2FromTo(&data, &val, pageMarkUp3, "PersonalSocial",
		[]string{"СТРАХОВЫЕ ВЗНОСЫ", "(МЕДИЦИНСКОЕ СТРАХОВАНИЕ)"})

	addImageLine(&data, pageMarkUp3, "3", &numLinePage3, FontSizeText1, ImagePosX, true)

	addMapData2FromTo(&data, &val, pageMarkUp3, "PersonalPension",
		[]string{"СТРАХОВЫЕ ВЗНОСЫ", "(ПЕНСИОННОЕ СТРАХОВАНИЕ)"})

	addImageLine(&data, pageMarkUp3, "3", &numLinePage3, FontSizeText1, ImagePosX, true)

	addMapData2FromTo(&data, &val, pageMarkUp3, "PersonalNDFL", []string{"НАЛОГ НА ДОХОДЫ", "ФИЗИЧЕСКИХ ЛИЦ", "(НДФЛ)"})

	//страница 4
	numLinePage4 := 0

	textMarkup4 := initTextMarkup(220)

	data["PagePos40"] = getNextRow(textMarkup4)
	data["PagePos41"] = getNextRow(textMarkup4)
	data["PagePos42"] = getNextRow(textMarkup4)
	data["PagePos43"] = getNextRow(textMarkup4)
	data["PagePos44"] = getNextParagraph(textMarkup4)
	data["PagePos45"] = getNextRow(textMarkup4)
	data["PagePos46"] = getNextParagraph(textMarkup4)
	data["PagePos47"] = getNextRow(textMarkup4)

	pageMarkUp4 := initPageMarkup(textMarkup4.CurPos, 3, 1100)

	res := addMapData1FromTo(&data, &val, pageMarkUp4, "EstatePrice", "СТОИМОСТЬ ЗЕМЛИ")

	addImageLine(&data, pageMarkUp4, "4", &numLinePage4, FontSizeText1, ImagePosX, res)

	addMapData1FromTo(&data, &val, pageMarkUp4, "EstateTax", "НАЛОГ НА ЗЕМЛЮ")

	addImageLine(&data, pageMarkUp4, "4", &numLinePage4, FontSizeText1, ImagePosX, res)

	addMapData2Val(&data, &val, pageMarkUp4, "EquipmentPrice", []string{"ОБОРУДОВАНИЕ И", "ИНЫЕ РАСХОДЫ"})
	data["EquipmentPriceUnits"] = getUnits(val["EquipmentPrice"])

	//страница 5
	numLinePage5 := 0
	textMarkup5 := initTextMarkup(220)

	data["PagePos50"] = getNextRow(textMarkup5)
	data["PagePos51"] = getNextRow(textMarkup5)
	data["PagePos52"] = getNextRow(textMarkup5)
	data["PagePos53"] = getNextParagraph(textMarkup5)
	data["PagePos54"] = getNextRow(textMarkup5)
	data["PagePos55"] = getNextRow(textMarkup5)

	pageMarkUp5 := initPageMarkup(textMarkup5.CurPos, 7, 1100)

	res = addMapData2FromTo(&data, &val, pageMarkUp5, "MoscowTax",
		[]string{"НАЛОГ В БЮДЖЕТ", "МОСКВЫ"})

	addImageLine(&data, pageMarkUp5, "5", &numLinePage5, FontSizeText1, ImagePosX, res)

	addMapData1FromTo(&data, &val, pageMarkUp5, "PropertyTax", "НАЛОГ НА ИМУЩЕСТВО")
	addImageLine(&data, pageMarkUp5, "5", &numLinePage5, FontSizeText1, ImagePosX, res)

	addMapData1FromTo(&data, &val, pageMarkUp5, "ProfitTax", "НАЛОГ НА ПРИБЫЛЬ")
	addImageLine(&data, pageMarkUp5, "5", &numLinePage5, FontSizeText1, ImagePosX, res)

	addMapData1FromTo(&data, &val, pageMarkUp5, "TransportTax", "ТРАНСПОРТНЫЙ НАЛОГ")
	addImageLine(&data, pageMarkUp5, "5", &numLinePage5, FontSizeText1, ImagePosX, res)

	addMapData1FromTo(&data, &val, pageMarkUp5, "OtherTax", "ПРОЧИЕ НАЛОГИ")
	addImageLine(&data, pageMarkUp5, "5", &numLinePage5, FontSizeText1, ImagePosX, res)

	addMapData1Val(&data, &val, pageMarkUp5, "PatentPrice", "СТОИМОСТЬ ПАТЕНТА")

	//страница 6
	numLinePage6 := 0

	textMarkup6 := initTextMarkup(220)

	data["PagePos60"] = getNextRow(textMarkup6)
	data["PagePos61"] = getNextRow(textMarkup6)
	data["PagePos62"] = getNextRow(textMarkup6)
	data["PagePos63"] = getNextRow(textMarkup6)
	data["PagePos64"] = getNextParagraph(textMarkup6)
	data["PagePos65"] = getNextRow(textMarkup6)

	pageMarkUp6 := initPageMarkup(textMarkup6.CurPos, 3, 1100)

	addMapData1FromTo(&data, &val, pageMarkUp6, "CapBuild", "КАПИТАЛЬНОЕ СТРОИТЕЛЬСТВО")
	addImageLine(&data, pageMarkUp6, "6", &numLinePage6, FontSizeText1, ImagePosX, res)

	addMapData1FromTo(&data, &val, pageMarkUp6, "CapRebuild", "КАПИТАЛЬНЫЙ РЕМОНТ")
	addImageLine(&data, pageMarkUp6, "6", &numLinePage6, FontSizeText1, ImagePosX, res)

	addMapData2FromTo(&data, &val, pageMarkUp6, "Financial",
		[]string{"СТОИМОСТЬ ОКАЗАНИЯ", "УСЛУГ ПО БУХГАЛТЕРСКОМУ", "УЧЕТУ"})

	//страница 7
	textMarkup7 := initTextMarkup(220)
	data["PageTitleTopPos7"] = "190"

	data["PagePos70"] = getNextRow(textMarkup7)
	data["PagePos71"] = getNextRow(textMarkup7)
	data["PagePos72"] = getNextRow(textMarkup7)
	data["PagePos73"] = getNextParagraph(textMarkup7)
	data["PagePos74"] = getNextRow(textMarkup7)
	data["PagePos75"] = getNextRow(textMarkup7)
	data["PagePos76"] = getNextRow(textMarkup7)
	data["PagePos77"] = getNextRow(textMarkup7)
	data["PagePos78"] = getNextParagraph(textMarkup7)
	data["PagePos79"] = getNextRow(textMarkup7)
	data["PagePos710"] = getNextRow(textMarkup7)
	data["PagePos711"] = getNextRow(textMarkup7)
	data["PagePos712"] = getNextRow(textMarkup7)
	data["PagePos713"] = getNextRow(textMarkup7)
	data["PagePos714"] = getNextParagraph(textMarkup7)
	data["PagePos715"] = getNextRow(textMarkup7)

	//страница 8
	data["PageTitleTopPos8"] = "190"
	textMarkup8 := initTextMarkup(220)
	data["PagePos80"] = getNextRow(textMarkup8)
	data["PagePos81"] = getNextRow(textMarkup8)
	data["PagePos82"] = getNextRow(textMarkup8)
	data["PagePos83"] = getNextRow(textMarkup8)
	data["PagePos84"] = getNextParagraph(textMarkup8)
	data["PagePos85"] = getNextRow(textMarkup8)
	data["PagePos86"] = getNextRow(textMarkup8)
	data["PagePos87"] = getNextParagraph(textMarkup8)
	data["PagePos88"] = getNextRow(textMarkup8)
	data["PagePos89"] = getNextRow(textMarkup8)
	data["PagePos810"] = getNextRow(textMarkup8)
	data["PagePos811"] = getNextParagraph(textMarkup8)
	data["PagePos812"] = getNextRow(textMarkup8)
	data["PagePos813"] = getNextRow(textMarkup8)
	data["PagePos814"] = getNextRow(textMarkup8)
	data["PagePos815"] = getNextParagraph(textMarkup8)
	data["PagePos816"] = getNextRow(textMarkup8)
	data["PagePos817"] = getNextRow(textMarkup8)

	//разметка

	data["LeftColumnPosPage3"] = "130"
	data["LeftColumnPosPage4"] = "130"
	data["LeftColumnPosPage5"] = "130"
	data["LeftColumnPosPage6"] = "130"

	data["LeftCalloutPos1"] = data["PageLeftPos1"]
	textCalloutPos1 := initTextMarkup2(1185, 18)
	data["CloudPos1"] = getNextRow(textCalloutPos1)
	data["CloudPos2"] = getNextRow(textCalloutPos1)
	data["CloudPos3"] = getNextRow(textCalloutPos1)
	data["CloudPos4"] = getNextRow(textCalloutPos1)

	data["RightColumnPosPage3"] = "490"
	data["RightColumnPosPage4"] = "490"
	data["RightColumnPosPage5"] = "490"
	data["RightColumnPosPage6"] = "490"

	data["PageTitleTopPos1"] = "115"
	data["PageTitleTopPos2"] = "190"

	data["Color0"] = "000000"
	data["Color1"] = "274183"
	data["Color2"] = "274183"
	data["Color3"] = "274183"
	data["Color4"] = "274183"
	data["Color5"] = "14153d"

	data["StyleTitle0"] = "sans-serif"

	data["StylePages0"] = "sans-serif"
	data["StylePages1"] = "sans-serif"
	data["StylePages2"] = "sans-serif"
	data["StylePages3"] = "sans-serif"
	data["StylePages4"] = "sans-serif"
	data["StylePages5"] = "sans-serif"

	if err := r.ParseTemplate(templatePath, data); err != nil {
		return nil, err
	}

	return r.GeneratePDF()
}

func addMapDataText1(data *map[string]string, pageMarkup *PageMarkup, name string, text string) {
	(*data)[name+"Pos1"] = getBodyNextPositionRow(pageMarkup)
	(*data)[name+"Text1"] = text
}

func addMapData1Val(data *map[string]string, val *map[string]float64, pageMarkup *PageMarkup, name string, text string) bool {

	if (*val)[name] > 0 {
		(*data)[name+"Pos1"] = getBodyNextPositionRow(pageMarkup)
		(*data)[name+"Column1"] = getBodyNextLeftPositionRow(pageMarkup)
		(*data)[name+"Text1"] = text
		(*data)[name] = floatToString((*val)[name])
		(*data)[name+"Units"] = getUnits((*val)[name])
		return true
	}
	return false
}

func addMapData1FromTo(data *map[string]string, val *map[string]float64, pageMarkup *PageMarkup, name string, text string) bool {

	if (*val)[name+"To"] > 0 {
		(*data)[name+"Pos1"] = getBodyNextPositionRow(pageMarkup)
		(*data)[name+"Column1"] = getBodyNextLeftPositionRow(pageMarkup)
		(*data)[name+"Text1"] = text
		(*data)[name+"From"] = floatToStringFrom((*val)[name+"From"], (*val)[name+"To"])
		(*data)[name+"To"] = floatToString((*val)[name+"To"])
		(*data)[name+"Units"] = getUnits((*val)[name+"To"])
		return true
	}
	return false
}

func addMapData1FromToFull(data *map[string]string, val *map[string]float64, pageMarkup *PageMarkup, name string, text string) bool {

	if (*val)[name+"To"] > 0 {
		(*data)[name+"Pos1"] = getBodyNextPositionRow(pageMarkup)
		(*data)[name+"Column1"] = getBodyNextLeftPositionRow(pageMarkup)
		(*data)[name+"Text1"] = text
		(*data)[name+"From"] = strconv.FormatInt(int64((*val)[name+"From"]), 10)
		(*data)[name+"To"] = strconv.FormatInt(int64((*val)[name+"To"]), 10)
		(*data)[name+"Units"] = getUnits(1.0)
		return true
	}
	return false
}

func addMapData2FromTo(data *map[string]string, val *map[string]float64, pageMarkup *PageMarkup, name string, vec []string) bool {
	if (*val)[name+"To"] > 0 {
		(*data)[name+"Pos1"] = getBodyNextPositionRow(pageMarkup)
		(*data)[name+"Text1"] = vec[0]

		for i := 1; i < len(vec); i++ {
			(*data)[name+"Pos"+getString(i+1)] = getBodyNextPosition(pageMarkup)
			(*data)[name+"Text"+getString(i+1)] = vec[i]
		}
		(*data)[name+"Column1"] = getBodyNextLeftPositionRow(pageMarkup)

		(*data)[name+"From"] = floatToStringFrom((*val)[name+"From"], (*val)[name+"To"])
		(*data)[name+"To"] = floatToString((*val)[name+"To"])
		(*data)[name+"Units"] = getUnits((*val)[name+"To"])
		return true
	}
	return false
}

func addMapData2Val(data *map[string]string, val *map[string]float64, pageMarkup *PageMarkup, name string, vec []string) bool {
	if (*val)[name] > 0 {
		(*data)[name+"Pos1"] = getBodyNextPositionRow(pageMarkup)
		(*data)[name+"Pos2"] = getBodyNextPosition(pageMarkup)
		(*data)[name+"Column1"] = getBodyNextLeftPositionRow(pageMarkup)
		(*data)[name+"Text1"] = vec[0]
		(*data)[name+"Text2"] = vec[1]
		(*data)[name] = floatToString((*val)[name])
		return true
	}
	return false
}

func addImageLine(data *map[string]string, pageMarkup *PageMarkup, numPage string, numLine *int, fontSize int, imagePosX int, res bool) {
	if res {
		(*data)["ImagePosY"+numPage+getString(*numLine)] = getImageLinePosition(pageMarkup, fontSize)
		(*data)["ImagePosX1"] = getString(imagePosX)
		*numLine++
	}
}
