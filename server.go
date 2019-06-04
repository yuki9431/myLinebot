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
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

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

	// 天気
	_, err = weather.PushMessage(userId, linebot.NewTextMessage("サーバ起動成功...")).Do()
	if err != nil {
		log.Print(err)
	}

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

		// 返信
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Hello")).Do()
				if err != nil {
					log.Print(err)
				}
			} else if event.Type == linebot.EventTypeFollow {
				// get userInfo
				profile, err := bot.GetProfile(event.Source.UserID).Do()

				// jsonファイルに書き込み
				jsonBytes, err := json.Marshal(profile)
				if err != nil {
					log.Print(err)
				}

				jsonFile, err := os.OpenFile("./userInfos.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
				jsonFile.Write(jsonBytes)
				defer jsonFile.Close()

				// フレンド登録時の挨拶
				message := profile.DisplayName + "さん\nはじめまして、毎朝6時に天気情報を教えてあげるね"

				_, err = weather.PushMessage(profile.UserID, linebot.NewTextMessage(message)).Do()
				if err != nil {
					log.Print(err)
				}

				log.Print("UserID: ", profile.UserID)
				log.Print("DisplayName: ", profile.DisplayName)

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
