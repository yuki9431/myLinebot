package main

import (
	"logger"

	"github.com/globalsign/mgo/bson"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/yuki9431/mongohelper"
)

// Follow ユーザがフォローした際の処理
func Follow(profile *linebot.UserProfileResponse) (replyMessage string, err error) {
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	if err != nil {
		return
	}

	userInfo := new(UserInfo)
	userInfo.UserID = profile.UserID
	userInfo.DisplayName = profile.DisplayName
	userInfo.CityID, _ = GetCityID("東京") //初回登録時には問答無用で東京民や
	userInfo.PictureURL = profile.PictureURL
	userInfo.StatusMessage = profile.StatusMessage

	// ユーザ情報をDBに登録
	mongo.InsertDb(userInfo, "userInfos")

	// フレンド登録時の挨拶
	var replyMessages [5]string
	replyMessages[0] = profile.DisplayName + "さん\nはじめまして、毎朝6時に天気情報を教えてあげるね"
	replyMessages[1] = usage
	replyMessages[2] = "お住まいの都市を変更するには、下記の通りメッセージをお送りください"
	replyMessages[3] = "都市変更:東京"
	replyMessages[4] = "都市変更:Brasil"

	for _, replyMessage = range replyMessages {
		replyMessage += replyMessage
	}

	return
}

// UnFollow ユーザがブロックした際の処理
func UnFollow(userID string, logger logger.Logger) {

	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	if err != nil {
		return
	}

	selector := bson.M{"userid": userID}
	if err := mongo.RemoveDb(selector, "userInfos"); err != nil {
		logger.Write(err)
	} else {
		logger.Write("success delete:" + userID)
	}
}
