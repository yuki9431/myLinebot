package main

import (
	"github.com/globalsign/mgo/bson"
	"github.com/yuki9431/mongoHelper"
)

// CityInfo 都市情報
type CityInfo struct {
	CityName string `json:"cityname"`
	CityID   string `json:"cityid"`
}

// GetAllCityList 都市一覧を返す *LineAPIの文字数制限に引っかかるため未使用
func GetAllCityList(cityList *[]string) (err error) {
	cityInfos := new([]CityInfo)

	mongo, err := mongoHelper.NewMongo(mongoDial, mongoName)
	if err == nil {

		// DBから都市一覧を取得
		if err = mongo.SearchDb(cityInfos, nil, "cityList"); err != nil {
			cityList = nil
		}

		for _, cityInfo := range *cityInfos {
			*cityList = append(*cityList, cityInfo.CityName)
		}

	}
	mongo.DisconnectDb()

	return
}

// GetCityInfo 都市情報を取得する
func GetCityInfo(cityInfo *CityInfo, cityID string) {
	cityInfos := new([]CityInfo)

	mongo, err := mongoHelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	if err == nil {

		// DBから都市一覧を取得
		selector := bson.M{"cityid": cityID}
		if err = mongo.SearchDb(cityInfos, selector, "cityList"); err != nil {

		}
		// 取得した情報をcityInfoに渡す
		cityInfo.CityID = (*cityInfos)[0].CityID
		cityInfo.CityName = (*cityInfos)[0].CityName

	}
}

// GetCityID 都市名から都市IDを抽出する
func GetCityID(cityName string) (cityID string, err error) {
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

// GetCityName 都市IDから都市名を取得
func GetCityName(cityID string) (cityName string, err error) {
	cityInfos := new([]CityInfo)

	mongo, err := mongoHelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	if err == nil {
		// DBから都市一覧を取得 1つだけ取得できる想定
		selector := bson.M{"cityid": cityID}
		err = mongo.SearchDb(cityInfos, selector, "cityList")
	}

	return (*cityInfos)[0].CityName, nil
}
