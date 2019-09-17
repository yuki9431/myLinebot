package main

import (
	"errors"
	"flag"
	"linebot/config"
	"log"
	"logger"
	"mongoHelper"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	logfile    = "/var/log/linebot.log"
	configFile = "../config/config.json"
	mongoDial  = "mongodb://localhost/mongodb"
	mongoName  = "mongodb"
)

// ユーザプロフィール情報
type UserInfo struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
	CityId        string `json:"cityId"`
}

// API等の設定
type ApiIds struct {
	ChannelSecret string `json:"channelSecret"`
	ChannelToken  string `json:"channelToken"`
	AppId         string `json:"appId"`
	CityId        string `json:"cityId"`
	CertFile      string `json:"certFile"`
	KeyFile       string `json:"keyFile"`
}

func main() {
	// log出力設定
	file, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logger := logger.New(file)

	// 設定ファイル読み込み
	apiIds := new(ApiIds)
	config := config.NewConfig(configFile)
	if err := config.Read(apiIds); err != nil {
		logger.Fatal(err)
	}

	// DB設定
	mongo, err := mongoHelper.NewMongo(mongoDial, mongoName)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Write("start information")

	// DBからユーザ情報を取得
	userinfos := new([]UserInfo)
	if err = mongo.SearchDb(userinfos, nil, "userInfos"); err != nil {
		return
	}

	// 抽出した全ユーザ情報にメッセージを配信
	for _, userinfo := range *userinfos {
		if userinfo.UserID != "" {
			var bot *linebot.Client
			if bot, err = linebot.New(apiIds.ChannelSecret, apiIds.ChannelToken); err != nil {
				// エラー時はその場で終了
				return

			} else {
				// CLIからメッセージ取得
				flag.Parse()
				messages := flag.Args()

				// メッセージ送信
				for _, message := range messages {
					_, err = bot.PushMessage(userinfo.UserID, linebot.NewTextMessage(message)).Do()
				}
			}
		} else {
			err = errors.New("Error: userId取得失敗")
			return
		}
	}

}
