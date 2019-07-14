package main

import (
	"errors"
	"log"
	"logger"
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

const (
	logfile       = "/var/log/linebot.log"
	configFile    = "config.json"
	followMessage = "さん\nはじめまして、毎朝6時に天気情報を教えてあげるね"
	usage         = `機能説明
天気　　 : 本日の天気情報を取得
おじさん : オジさん？に呼びかける

comming soon...
都市変更 : 天気情報取得の所在地を変更する
時間変更 : 毎日の天気配信時刻を変更する

https://github.com/yuki9431/myLinebot`
)

// ユーザプロフィール情報
type UserInfos struct {
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

// 天気情報作成
func createWeatherMessage(apiIds *ApiIds) (message string, err error) {
	// 設定ファイル読み込み
	config := NewConfig(configFile)
	if err = config.Read(apiIds); err != nil {
		return
	}

	cityId := apiIds.CityId
	appId := apiIds.AppId

	// 今日の天気情報を取得
	w := weather.New(cityId, appId).GetInfoFromDate(time.Now())
	dates := w.GetDates()
	icons := w.GetIcons()
	cityName := w.GetCityName()
	descriptions := w.GetDescriptions()

	// 天気情報メッセージ作成
	message = cityName + "\n" +
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

	return
}

// 天気配信ジョブ
func sendWeatherInfo(apiIds *ApiIds) (err error) {
	const layout = "15:04:05" // => hh:mm:ss
	userinfos := new([]UserInfos)

	for {
		t := time.Now()
		if t.Format(layout) == "06:00:00" {
			// DBからユーザ情報を取得
			if err = searchDb(userinfos, "userInfos"); err != nil {
				return
			}

			// 抽出した全ユーザ情報に天気情報を配信
			for _, userinfo := range *userinfos {
				if userinfo.UserID != "" {
					var bot *linebot.Client
					if bot, err = linebot.New(apiIds.ChannelSecret, apiIds.ChannelToken); err != nil {
						return

					} else {
						// 天気情報メッセージ送信
						var message string
						message, err = createWeatherMessage(apiIds)
						_, err = bot.PushMessage(userinfo.UserID, linebot.NewTextMessage(message)).Do()
					}
				} else {
					err = errors.New("Error: userId取得失敗")
					return
				}
			}

			// 連続送信を防止する
			time.Sleep(1 * time.Second) // sleep 1 second
		}
	}
}

// ojichat実装
func ojichat(name string) (result string, err error) {
	parser := &docopt.Parser{
		OptionsFirst: true,
	}
	args, _ := parser.ParseArgs("", nil, "")
	config := generator.Config{}
	config.TargetName = name
	err = args.Bind(&config)

	result, err = generator.Start(config)

	return result, err
}

// mongoDB接続
func connectDb() (d *mgo.Database, err error) {
	session, err := mgo.Dial("mongodb://localhost/mongodb")
	d = session.DB("mongodb")

	return
}

func disconnectDb(db *mgo.Database) {
	db.Session.Close()
}

// mongoDB挿入
func insertDb(obj interface{}, colectionName string) (err error) {
	db, err := connectDb()
	if err != nil {
		return
	}
	defer disconnectDb(db)

	col := db.C(colectionName)
	return col.Insert(obj)
}

// mongoDB削除
func removeDb(obj interface{}, colectionName string) (err error) {
	db, err := connectDb()
	if err != nil {
		return
	}
	defer disconnectDb(db)

	col := db.C(colectionName)
	_, err = col.RemoveAll(obj)
	return
}

// mondoDB抽出
func searchDb(obj interface{}, colectionName string) (err error) {
	db, err := connectDb()
	if err != nil {
		return
	}
	defer disconnectDb(db)

	col := db.C(colectionName)
	return col.Find(nil).All(obj)
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
	config := NewConfig(configFile)
	if err := config.Read(apiIds); err != nil {
		logger.Fatal(err)
	}

	// 指定時間に天気情報を配信
	go func() {
		if err := sendWeatherInfo(apiIds); err != nil {
			logger.Write(err)
		}
	}()

	handler, err := httphandler.New(apiIds.ChannelSecret, apiIds.ChannelToken)
	if err != nil {
		logger.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		bot, err := handler.NewClient()
		if err != nil {
			logger.Fatal(err)
			return
		}

		// イベント処理
		for _, event := range events {
			logger.Write("start event : " + event.Type)

			// ユーザのIDを取得
			userId := event.Source.UserID
			logger.Write("userid :" + userId)

			// ユーザのプロフィールを取得後、レスポンスする
			if profile, err := bot.GetProfile(userId).Do(); err == nil {
				if event.Type == linebot.EventTypeMessage {
					// 返信メッセージ
					var replyMessage string

					switch message := event.Message.(type) {
					case *linebot.TextMessage:
						if strings.Contains(message.Text, "天気") {
							if replyMessage, err = createWeatherMessage(apiIds); err != nil {
								logger.Write(err)
							}

						} else if strings.Contains(message.Text, "おじさん") {
							if replyMessage, err = ojichat(profile.DisplayName); err != nil {
								logger.Write(err)
							}

						} else {
							replyMessage = usage
						}

						// 返信処理
						if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
							logger.Write(err)
						}
						logger.Write("message.Text: " + message.Text)
					}
				} else if event.Type == linebot.EventTypeFollow {
					// TODO insert前に存在の確認

					// ユーザ情報をDBに登録
					if err := insertDb(profile, "userInfos"); err != nil {
						logger.Write(err)
					}

					// フレンド登録時の挨拶
					replyMessage := profile.DisplayName + followMessage

					if _, err = bot.PushMessage(userId, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						logger.Write(err)
					}
				}
			}

			// ブロック時の処理、ユーザ情報をDBから削除する
			if event.Type == linebot.EventTypeUnfollow {
				query := bson.M{"userid": userId}
				if err := removeDb(query, "userInfos"); err != nil {
					logger.Write(err)
				}
			}

			logger.Write("end event")
		}
	})
	http.Handle("/callback", handler)
	if err := http.ListenAndServeTLS(":443", apiIds.CertFile, apiIds.KeyFile, nil); err != nil {
		logger.Fatal("ListenAndServe: ", err)
	}

	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	logger.Fatal(err)
	// }
}
