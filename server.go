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
	"os"
	"time"
	"weather"

	"github.com/docopt/docopt-go"
	"github.com/greymd/ojichat/generator"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"gopkg.in/mgo.v2"
)

const (
	// TODO プロパティファイルに外だし
	channelSecret = "f7c8a8c3df6f23c2549f3f1ed484dc47"
	channelToken  = "fmgf96KJrTdN7B/T2aS39L9XDycHqS86H0F09ekR/mtUadt+R3sY1eYba8R6h0ifJ3yqmATJq9117er8GtipA2LgN81xluam/udbmUoluWJeS2GQQyFSKsl9djd/yytyEh9Q/8un3gFIZJ/op1Dz+wdB04t89/1O/w1cDnyilFU="
	appId         = "63ef79e871474934c1bd707239475660"
	cityId        = "1850147" // Tokyo
	certFile      = "/etc/letsencrypt/live/blacksnowpi.f5.si-0001/fullchain.pem"
	keyFile       = "/etc/letsencrypt/live/blacksnowpi.f5.si-0001/privkey.pem"
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
	w := weather.New(cityId, appId).GetInfoFromDate(time.Now())
	dates := w.GetDates()
	icons := w.GetIcons()
	cityName := w.GetCityName()

	// 天気情報メッセージ作成
	message := "TODO:\n降水確率追加\n天気アイコン追加\nフォーマットを見やすく\n\n" +
		cityName + "の天気情報です♪\n\n"
	for i, date := range dates {
		message = message +
			date.Format("01月02日 15時04分") + "時点の天気は" +
			w.ConvertIconToWord(icons[i]) + "でしょう。\n\n"
	}

	return message
}

// 天気配信ジョブ
func sendWeatherInfo(c *linebot.Client) {
	const layout = "15:04:05" // => hh:mm:ss
	for {
		t := time.Now()
		if t.Format(layout) == "06:00:00" {
			// DBからユーザ情報を取得
			var userinfos []UserInfos
			searchDb(&userinfos, "userInfos")

			// 抽出した全ユーザ情報に天気情報を配信
			for _, userinfo := range userinfos {

				// 天気情報メッセージ送信
				message := createWeatherMessage()

				_, err := c.PushMessage(userinfo.UserID, linebot.NewTextMessage(message)).Do()
				if err != nil {
					log.Print(err)
				}
			}

			// 連続送信を防止する
			time.Sleep(1 * time.Second) // sleep 1 second
		}
	}
}

// ojichat実装
func ojichat(name string) string {
	parser := &docopt.Parser{
		OptionsFirst: true,
	}
	args, _ := parser.ParseArgs("", nil, "")
	config := generator.Config{}
	config.TargetName = name
	err := args.Bind(&config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	result, err := generator.Start(config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return result
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

// mondoDB抽出
func searchDb(obj interface{}, colectionName string) {
	col := connectDb().C(colectionName)
	if err := col.Find(nil).All(obj); err != nil {
		log.Fatal(err)
	}
}

func main() {
	handler, err := httphandler.New(channelSecret, channelToken)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := handler.NewClient()
	if err != nil {
		log.Print(err)
		return
	}

	go sendWeatherInfo(bot)

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {

		// イベント処理
		for _, event := range events {
			// get userInfo
			profile, err := bot.GetProfile(event.Source.UserID).Do()
			userInfos := new(UserInfos)
			userInfos.UserID = profile.UserID
			userInfos.DisplayName = profile.DisplayName
			userInfos.PictureURL = profile.PictureURL
			userInfos.StatusMessage = profile.StatusMessage

			if event.Type == linebot.EventTypeMessage {
				// 天気情報をあげる
				//message := event.Message

				//if message == "天気"

				_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ojichat(userInfos.DisplayName))).Do()
				if err != nil {
					log.Print(err)
				}
			} else if event.Type == linebot.EventTypeFollow {

				// ユーザ情報をDBに登録
				db := connectDb()
				defer disconnectDb(db)
				insertDb(userInfos, "userInfos")

				// フレンド登録時の挨拶
				message := profile.DisplayName + "さん\nはじめまして、毎朝6時に天気情報を教えてあげるね"

				_, err = bot.PushMessage(profile.UserID, linebot.NewTextMessage(message)).Do()
				if err != nil {
					log.Print(err)
				}
			}
		}
	})
	http.Handle("/callback", handler)
	if err := http.ListenAndServeTLS(":443", certFile, keyFile, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
