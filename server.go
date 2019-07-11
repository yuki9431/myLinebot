package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"weather"

	"github.com/docopt/docopt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/greymd/ojichat/generator"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"gopkg.in/mgo.v2"
)

const configFile = "config.json"

const usage = `機能説明
天気　　 : 本日の天気情報を取得
おじさん : オジさん？に呼びかける

comming soon...
都市変更 : 天気情報取得の所在地を変更する
時間変更 : 毎日の天気配信時刻を変更する

https://github.com/yuki9431/myLinebot`

// ユーザプロフィール情報
type UserInfos struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

// API等の設定
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
	descriptions := w.GetDescriptions()

	// 天気情報メッセージ作成
	message := cityName + "\n" +
		func() string {
			var tempIcon string
			for i, time := range dates {
				tempIcon += time.Format("15:04") + " " +
					w.ConvertIconToWord(icons[i]) + " " +
					descriptions[i] + "\n"
			}

			wdays := [...]string{"日", "月", "火", "水", "木", "金", "土"}

			return dates[0].Format("01/02 (") + wdays[dates[0].Weekday()] + ")" +
				"の天気情報だよ" + "\n" + tempIcon
		}()

	return message
}

// 天気配信ジョブ
func sendWeatherInfo() {
	const layout = "15:04:05" // => hh:mm:ss
	userinfos := new([]UserInfos)
	var c linebot.Client
	for {
		t := time.Now()
		if t.Format(layout) == "06:00:00" {
			// DBからユーザ情報を取得
			searchDb(userinfos, "userInfos")

			// 抽出した全ユーザ情報に天気情報を配信
			for _, userinfo := range *userinfos {

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
	db := connectDb()
	defer disconnectDb(db)
	col := db.C(colectionName)
	if err := col.Insert(obj); err != nil {
		log.Fatal(err)
	}
}

// mongoDB削除
func removeDb(obj interface{}, colectionName string) {
	db := connectDb()
	defer disconnectDb(db)
	col := db.C(colectionName)
	if _, err := col.RemoveAll(obj); err != nil {
		log.Println(err)
	}
}

// mondoDB抽出
func searchDb(obj interface{}, colectionName string) {
	db := connectDb()
	defer disconnectDb(db)
	col := db.C(colectionName)
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
			log.Print("start event : " + event.Type + "\n")

			// ユーザのIDを取得
			userId := event.Source.UserID

			// ユーザのプロフィールを取得後、レスポンスする
			if profile, err := bot.GetProfile(userId).Do(); err == nil {
				if event.Type == linebot.EventTypeMessage {
					// 返信メッセージ
					var replyMessage string

					switch message := event.Message.(type) {
					case *linebot.TextMessage:
						if strings.Contains(message.Text, "天気") {
							replyMessage = createWeatherMessage()

						} else if strings.Contains(message.Text, "おじさん") {
							replyMessage = ojichat(profile.DisplayName)

						} else {
							replyMessage = usage
						}

						// 返信処理
						_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
						if err != nil {
							log.Print(err)
						}
					}
				} else if event.Type == linebot.EventTypeFollow {
					// TODO insert前に存在の確認

					// ユーザ情報をDBに登録
					insertDb(profile, "userInfos")

					// フレンド登録時の挨拶
					replyMessage := profile.DisplayName + "さん\nはじめまして、毎朝6時に天気情報を教えてあげるね"

					_, err = bot.PushMessage(userId, linebot.NewTextMessage(replyMessage)).Do()
					if err != nil {
						log.Print(err)
					}
				}
			}

			// ブロック時の処理、ユーザ情報をDBから削除する
			if event.Type == linebot.EventTypeUnfollow {
				query := bson.M{"userid": userId}
				removeDb(query, "userInfos")
			}

			log.Println("end event")
		}
	})
	http.Handle("/callback", handler)
	if err := http.ListenAndServeTLS(":443", config.CertFile, config.KeyFile, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatal(err)
	// }
}
