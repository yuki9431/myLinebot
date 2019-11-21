package main

import (
	"errors"

	"github.com/globalsign/mgo/bson"
	"github.com/yuki9431/mongohelper"
)

// CityInfo 都市情報
type CityInfo struct {
	CityName string `json:"cityname"`
	CityID   string `json:"cityid"`
}

// ShowCityList 都市一覧を取得しメッセージを返す
func ShowCityList() (replyMessage string, err error) {
	cityList := new([]string)
	err = GetAllCityList(cityList)

	replyMessage = "都市一覧\n"
	for _, city := range *cityList {
		replyMessage = replyMessage + city + "\n"
	}

	return
}

// ChangeCity ユーザの所在地を変更する
func ChangeCity(userInfo, cityName string) (replyMessage string, err error) {

	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	if err != nil {
		return
	}

	// 都市IDを取得する
	cityID, err := GetCityID(cityName)
	if err != nil {
		err = errors.New("error: failed get cityID")
	}

	// 都市IDをDBに登録する
	if cityID != "" && cityName != "" {

		selector := bson.M{"userid": userInfo}
		update := bson.M{"$set": bson.M{"cityid": cityID}}

		if err := mongo.UpdateDb(selector, update, "userInfos"); err == nil {
			replyMessage = "選択された都市に変更しました！"
		} else {
			replyMessage = "都市の変更に失敗しました..."
		}

	} else {
		replyMessage = "該当都市が見つかりません💦\n" +
			"\"都市一覧\"と送り頂ければ設定可能な都市が表示されますよ"
	}

	return
}

// GetAllCityList 都市一覧を返す
func GetAllCityList(cityList *[]string) (err error) {
	cityInfos := new([]CityInfo)

	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
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

	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
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

	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	if err == nil {
		// DBから都市一覧を取得 1つだけ取得できる想定
		selector := bson.M{"cityname": bson.M{"$regex": "^" + cityName + ".*"}}
		err = mongo.SearchDb(cityInfos, selector, "cityList")
	}

	if *cityInfos != nil {
		cityID = (*cityInfos)[0].CityID
	} else {
		cityID = ""
	}

	return
}

// GetCityName 都市IDから都市名を取得
func GetCityName(cityID string) (cityName string, err error) {
	cityInfos := new([]CityInfo)

	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	if err == nil {
		// DBから都市一覧を取得 1つだけ取得できる想定
		selector := bson.M{"cityid": cityID}
		err = mongo.SearchDb(cityInfos, selector, "cityList")
	}

	if *cityInfos != nil {
		cityName = (*cityInfos)[0].CityName
	} else {
		cityName = ""
	}

	return
}
