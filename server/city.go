package main

import (
	"errors"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/yuki9431/mongoHelper"
)

// 都市情報
type CityInfo struct {
	CityID  string
	Name    string
	Country string
}

// 都市一覧を返す *LineAPIの文字数制限に引っかかるため未使用
func GetAllCityList(cityList *[]string) (err error) {
	cityInfos := new([]CityInfo)

	mongo, err := mongoHelper.NewMongo(mongoDial, mongoName)
	if err == nil {

		// DBから都市一覧を取得
		if err = mongo.SearchDb(cityInfos, nil, "cityList"); err != nil {
			cityList = nil
		}

		for _, cityInfo := range *cityInfos {
			*cityList = append(*cityList, cityInfo.Name)
		}

	}
	mongo.DisconnectDb()

	return
}

// 都市情報を取得する
func GetCityInfo(cityInfo *CityInfo, cityId string) {
	cityInfos := new([]CityInfo)

	mongo, err := mongoHelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	if err == nil {

		// DBから都市一覧を取得
		selector := bson.M{"cityid": cityId}
		if err = mongo.SearchDb(cityInfos, selector, "cityInfo"); err != nil {

		}
		// 取得した情報をcityInfoに渡す
		cityInfo.CityID = (*cityInfos)[0].CityID
		cityInfo.Country = (*cityInfos)[0].Country
		cityInfo.Name = (*cityInfos)[0].Name

	}
}

// 都市名から都市IDを抽出する
func GetCityId(cityName string) (cityId string, err error) {
	cityInfos := new([]CityInfo)

	mongo, err := mongoHelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	if err == nil {
		// DBから都市一覧を取得 1つだけ取得できる想定
		selector := bson.M{"cityname": cityName}
		err = mongo.SearchDb(cityInfos, selector, "cityList")
	}

	return (*cityInfos)[0].CityID, nil
}

// 都市IDから都市名を取得
func GetCityName(cityId string) (cityName string, err error) {
	cityInfos := new([]CityInfo)

	mongo, err := mongoHelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	if err == nil {
		// DBから都市一覧を取得 1つだけ取得できる想定
		selector := bson.M{"cityid": cityId}
		err = mongo.SearchDb(cityInfos, selector, "cityList")
	}

	return (*cityInfos)[0].Name, nil
}

// 都道府県をIDに変換
func ConvertCityToId(cityName string) (cityId string, err error) {
	switch {
	case strings.Contains(cityName, "北海道"):
		cityId = "2130037"
	case strings.Contains(cityName, "青森"):
		cityId = "2130656"
	case strings.Contains(cityName, "岩手"):
		cityId = "2112518"
	case strings.Contains(cityName, "宮城"):
		cityId = "2111888"
	case strings.Contains(cityName, "秋田"):
		cityId = "2113124"
	case strings.Contains(cityName, "山形"):
		cityId = "2110554"
	case strings.Contains(cityName, "福島"):
		cityId = "2112922"
	case strings.Contains(cityName, "茨城"):
		cityId = "2112669"
	case strings.Contains(cityName, "栃木"):
		cityId = "1850310"
	case strings.Contains(cityName, "群馬"):
		cityId = "1863501"
	case strings.Contains(cityName, "埼玉"):
		cityId = "1853226"
	case strings.Contains(cityName, "千葉"):
		cityId = "2113014"
	case strings.Contains(cityName, "東京"):
		cityId = "1850147"
	case strings.Contains(cityName, "神奈川"):
		cityId = "1860291"
	case strings.Contains(cityName, "新潟"):
		cityId = "1855429"
	case strings.Contains(cityName, "富山"):
		cityId = "1849872"
	case strings.Contains(cityName, "石川"):
		cityId = "1861387"
	case strings.Contains(cityName, "福井"):
		cityId = "1863983"
	case strings.Contains(cityName, "山梨"):
		cityId = "1848649"
	case strings.Contains(cityName, "長野"):
		cityId = "1856210"
	case strings.Contains(cityName, "岐阜"):
		cityId = "1863640"
	case strings.Contains(cityName, "静岡"):
		cityId = "1851715"
	case strings.Contains(cityName, "愛知"):
		cityId = "1865694"
	case strings.Contains(cityName, "三重"):
		cityId = "1857352"
	case strings.Contains(cityName, "滋賀"):
		cityId = "1852553"
	case strings.Contains(cityName, "長浜"):
		cityId = "1856239"
	case strings.Contains(cityName, "京都"):
		cityId = "1857910"
	case strings.Contains(cityName, "大阪"):
		cityId = "1853908"
	case strings.Contains(cityName, "兵庫"):
		cityId = "1847966"
	case strings.Contains(cityName, "奈良"):
		cityId = "1855608"
	case strings.Contains(cityName, "和歌山"):
		cityId = "1848938"
	case strings.Contains(cityName, "鳥取"):
		cityId = "1849890"
	case strings.Contains(cityName, "島根"):
		cityId = "1859687"
	case strings.Contains(cityName, "岡山"):
		cityId = "1854381"
	case strings.Contains(cityName, "広島"):
		cityId = "1862413"
	case strings.Contains(cityName, "山口"):
		cityId = "1848681"
	case strings.Contains(cityName, "徳島"):
		cityId = "1850157"
	case strings.Contains(cityName, "香川"):
		cityId = "1860834"
	case strings.Contains(cityName, "愛媛"):
		cityId = "1864226"
	case strings.Contains(cityName, "高知"):
		cityId = "1859146"
	case strings.Contains(cityName, "福岡"):
		cityId = "1863958"
	case strings.Contains(cityName, "佐賀"):
		cityId = "1853299"
	case strings.Contains(cityName, "長崎"):
		cityId = "1856156"
	case strings.Contains(cityName, "熊本"):
		cityId = "1858419"
	case strings.Contains(cityName, "大分"):
		cityId = "1854487"
	case strings.Contains(cityName, "宮崎"):
		cityId = "1856710"
	case strings.Contains(cityName, "鹿児島"):
		cityId = "1860825"
	case strings.Contains(cityName, "沖縄"):
		cityId = "1854345"
	case strings.Contains(cityName, "Brasil") || strings.Contains(cityName, "brasil"):
		cityId = "3448439"
	default:
		cityId = "1850147"
		err = errors.New("err: 該当する都市が見つかりません")
	}

	return
}
