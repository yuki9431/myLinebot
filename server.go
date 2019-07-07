package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"weather"

	"github.com/docopt/docopt-go"
	"github.com/greymd/ojichat/generator"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"gopkg.in/mgo.v2"
)

const configFile = "config.json"

// 構造体定義
type UserInfos struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

// 構造体定義
type Config struct {
	ChannelSecret string `json:"channelSecret"`
	ChannelToken  string `json:"channelToken"`
	AppId         string `json:"appId"`
	CityId        string `json:"cityId"`
	CertFile      string `json:"certFile"`
	KeyFile       string `json:"keyFile"`
}

// 天気情報作成
func createWeatherMessage() string {
	// 設定ファイル読み込み
	config := new(Config)
	if err := Read(config, configFile); err != nil {
		log.Fatal(err)
	}

	cityId := config.CityId
	appId := config.AppId

	// 今日の天気情報を取得
	w := weather.New(cityId, appId).GetInfoFromDate(time.Now())
	dates := w.GetDates()
	icons := w.GetIcons()
	cityName := w.GetCityName()

	// 天気情報メッセージ作成
	message := cityName + "\n" +
		func() string {
			var times string
			for _, time := range dates {
				times = times + " " + time.Format("15時")
			}
			return times
		}()

	for _, icon := range icons {
		message = message + w.ConvertIconToWord(icon) + "    "
	}

	return message
}

// 天気配信ジョブ
func sendWeatherInfo() {
	const layout = "15:04:05" // => hh:mm:ss
	var userinfos []UserInfos
	var c linebot.Client
	for {
		t := time.Now()
		if t.Format(layout) == "06:00:00" {
			// DBからユーザ情報を取得
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
		log.Print(err)
		os.Exit(1)
	}

	result, err := generator.Start(config)
	if err != nil {
		log.Print(err)
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
		log.Fatal(err)
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
	// 設定ファイル読み込み
	config := new(Config)
	if err := Read(config, configFile); err != nil {
		log.Fatal(err)
	}

	// 指定時間に天気情報を配信
	go sendWeatherInfo()

	handler, err := httphandler.New(config.ChannelSecret, config.ChannelToken)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		bot, err := handler.NewClient()
		if err != nil {
			log.Fatal(err)
			return
		}

		// イベント処理
		for _, event := range events {
			// get userInfo
			var profile *linebot.UserProfileResponse
			if profile, err = bot.GetProfile(event.Source.UserID).Do(); err != nil {
				log.Print("err:ユーザの情報取得失敗")
				return
			}

			if event.Type == linebot.EventTypeMessage {
				// 返信メッセージ
				var replyMessage string

				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if strings.Contains(message.Text, "天気") {
						replyMessage = createWeatherMessage()
					} else if strings.Contains(message.Text, "ヘルプ") {
						replyMessage = "TODO 機能説明"
					} else {
						replyMessage = ojichat(profile.DisplayName)
					}

					// 返信処理
					_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
					if err != nil {
						log.Print(err)
					}
				}
			} else if event.Type == linebot.EventTypeFollow {
				// ユーザ情報をDBに登録
				db := connectDb()
				defer disconnectDb(db)
				insertDb(profile, "userInfos")

				// フレンド登録時の挨拶
				message := profile.DisplayName + "さん\nはじめまして、毎朝6時に天気情報を教えてあげるね"

				_, err = bot.PushMessage(profile.UserID, linebot.NewTextMessage(message)).Do()
				if err != nil {
					log.Print(err)
				}
			} else if event.Type == linebot.EventTypeUnfollow {
				// DBから削除する処理
			}
		}
	})
	http.Handle("/callback", handler)
	if err := http.ListenAndServeTLS(":443", config.CertFile, config.KeyFile, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
