package main

import (
	"github.com/globalsign/mgo/bson"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/yuki9431/mongohelper"
)

// Follow ユーザがフォローした際の処理
func Follow(profile *linebot.UserProfileResponse) (replyMessages []string, err error) {
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
	replyMessages = append(replyMessages, profile.DisplayName+"さん\nはじめまして、毎朝6時に天気情報を教えてあげるね")
	replyMessages = append(replyMessages, usage)
	replyMessages = append(replyMessages, "お住まいの都市を変更するには、下記の通りメッセージをお送りください")
	replyMessages = append(replyMessages, "都市変更:東京")
	replyMessages = append(replyMessages, "都市変更:Brasil")

	return
}

// UnFollow ユーザがブロックした際の処理
func UnFollow(userID string) (err error) {

	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	if err != nil {
		return
	}

	selector := bson.M{"userid": userID}
	err = mongo.RemoveDb(selector, "userInfos")
	return
}
