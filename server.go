// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// +build !appengine

package main

import (
	"log"
	"net/http"
	"time"
	"weather"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"gopkg.in/mgo.v2"
)

type UserInfos struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

// 天気情報作成
func createWeatherMessage() string {
	// 今日の天気情報を取得
	w := weather.New().GetInfoFromDate(time.Now())
	dates := w.GetDates()
	icons := w.GetIcons()
	cityName := w.GetCityName()

	// 天気情報メッセージ作成
	message := cityName + "の天気情報です♪\n\n"
	for i, date := range dates {
		message = message +
			date.Format("01月02日 15時04分") + "時点の天気は" +
			w.ConvertIconToWord(icons[i]) + "でしょう。\n\n"
	}

	return message
}

// 天気配信ジョブ
func sendWeatherInfo(c *linebot.Client, userId string) {
	const layout = "15:04:05" // => hh:mm:ss
	for {
		t := time.Now()
		if t.Format(layout) == "16:40:00" {
			// 天気情報メッセージ送信
			message := createWeatherMessage()
			_, err := c.PushMessage(userId, linebot.NewTextMessage(message)).Do()
			if err != nil {
				log.Print(err)
			}

			// 連続送信を防止する
			time.Sleep(1 * time.Second) // sleep 1 second
		}
	}
}

// mongoDB接続
func connectDb() *mgo.Database {
	session, err := mgo.Dial("mongodb://localhost/mongodb")
	if err != nil {
		log.Fatal(err)
	}

	return session.DB("mongodb")
}

func disconnectDb(db *mgo.Database) {
	db.Session.Close()
}

// mongoDB挿入
func insertDb(obj interface{}, colectionName string) {
	col := connectDb().C(colectionName)
	if err := col.Insert(obj); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	const (
		channelSecret = "f7c8a8c3df6f23c2549f3f1ed484dc47"
		channelToken  = "fmgf96KJrTdN7B/T2aS39L9XDycHqS86H0F09ekR/mtUadt+R3sY1eYba8R6h0ifJ3yqmATJq9117er8GtipA2LgN81xluam/udbmUoluWJeS2GQQyFSKsl9djd/yytyEh9Q/8un3gFIZJ/op1Dz+wdB04t89/1O/w1cDnyilFU="
		userId        = "U0b00920127574259c8ac979e5f59f0ea"
	)

	weather, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		log.Print(err)
	}

	// サーバ起動確認
	_, err = weather.PushMessage(userId, linebot.NewTextMessage("サーバ起動成功...")).Do()
	if err != nil {
		log.Print(err)
	}

	go sendWeatherInfo(weather, userId)

	handler, err := httphandler.New(channelSecret, channelToken)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		bot, err := handler.NewClient()
		if err != nil {
			log.Print(err)
			return
		}

		// イベント処理
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(createWeatherMessage())).Do()
				if err != nil {
					log.Print(err)
				}
			} else if event.Type == linebot.EventTypeFollow {
				// get userInfo
				profile, err := bot.GetProfile(event.Source.UserID).Do()
				userInfos := new(UserInfos)
				userInfos.UserID = profile.UserID
				userInfos.DisplayName = profile.DisplayName
				userInfos.PictureURL = profile.PictureURL
				userInfos.StatusMessage = profile.StatusMessage

				// ユーザ情報をDBに登録
				db := connectDb()
				defer disconnectDb(db)
				insertDb(userInfos, "userInfos")

				// フレンド登録時の挨拶
				message := profile.DisplayName + "さん\nはじめまして、毎朝6時に天気情報を教えてあげるね"

				_, err = weather.PushMessage(profile.UserID, linebot.NewTextMessage(message)).Do()
				if err != nil {
					log.Print(err)
				}
			}
		}
	})
	http.Handle("/callback", handler)
	// This is just a sample code.
	// For actually use, you must support HTTPS by using `ListenAndServeTLS`, reverse proxy or etc.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
