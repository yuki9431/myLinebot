package main

import (
	"flag"
	"log"
	"myLinebot/config"
	"net/http"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/greymd/ojichat/generator"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"github.com/yuki9431/logger"
	"github.com/yuki9431/mongohelper"
)

const (
	logfile    = "/var/log/linebot.log"
	configFile = "config/config.json"
	mongoDial  = "mongodb://localhost/mongodb"
	mongoName  = "mongodb"

	usage = "レスポンス説明\n" +
		"[天気]\n" +
		"  本日の天気情報を取得\n\n" +
		"[おじさん]\n" +
		"  オジさん？に呼びかける\n\n" +
		"[都市変更:都道府県]\n" +
		"  天気情報取得の所在地を変更する"
)

// UserInfo ユーザプロフィール情報
type UserInfo struct {
	UserID        string `json:"userID"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
	CityID        string `json:"cityId"`
}

// APIIDs API等の設定
type APIIDs struct {
	ChannelSecret string `json:"channelSecret"`
	ChannelToken  string `json:"channelToken"`
	AppID         string `json:"appId"`
	CityID        string `json:"cityId"`
	CertFile      string `json:"certFile"`
	KeyFile       string `json:"keyFile"`
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

func main() {
	// log出力設定
	file, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logger := logger.New(file)

	// 設定ファイル読み込み
	apiIDs := new(APIIDs)
	config := config.NewConfig(configFile)
	if err := config.Read(apiIDs); err != nil {
		logger.Fatal(err)
	}

	// 指定時間に天気情報を配信
	go func() {
		if err := sendWeatherInfo(apiIDs); err != nil {
			logger.Write(err)
		}
	}()

	handler, err := httphandler.New(apiIDs.ChannelSecret, apiIDs.ChannelToken)
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
			// DB設定
			mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
			if err != nil {
				logger.Fatal(err)
			}

			logger.Write("start event : " + event.Type)

			// ユーザのIDを取得
			userID := event.Source.UserID
			logger.Write("userid :" + userID)

			// 都市IDを取得するため、DBからユーザ情報を獲得
			userInfos := new([]UserInfo)
			if err := mongo.SearchDb(userInfos, bson.M{"userid": userID}, "userInfos"); err != nil {
				logger.Write("err search userInfo" + err.Error())
				return
			}

			// APIからユーザのプロフィールを取得後、レスポンスする
			if profile, err := bot.GetProfile(userID).Do(); err == nil {
				if event.Type == linebot.EventTypeMessage {
					// 返信メッセージ
					var replyMessage string

					switch message := event.Message.(type) {
					case *linebot.TextMessage:

						if IsAskWeather(message.Text) {
							if replyMessage, err = createWeatherMessage(apiIDs, (*userInfos)[0]); err != nil { // (*userInfos)[0]は一意の値しか取れない想定
								logger.Write(err)
							}

						} else if IsOjichan(message.Text) {
							if replyMessage, err = ojichat(profile.DisplayName); err != nil {
								logger.Write(err)
							}

						} else if IsChangeCity(message.Text) {
							cityName := strings.Replace(message.Text, " ", "", -1) // 全ての半角スペースを消す
							cityName = strings.Replace(cityName, "都市変更:", "", 1)   // 頭の都市変更:を消す

							replyMessage = ChangeCity(profile.UserID, cityName, logger)

						} else if IsShowCityList(message.Text) {
							if replyMessage, err = ShowCityList(); err != nil {
								logger.Write(err)
							}

						} else if IsShowHelp(message.Text) {
							// botの機能を返信する
							replyMessage = usage

						} else {
							// 100%の晴れ女
							if replyMessage, err = HinaResponce(); err != nil {
								logger.Write(err)
							}

						}

						// 返信処理
						if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
							logger.Write(err)
						}
						logger.Write("message.Text: " + message.Text)
					}

				} else if event.Type == linebot.EventTypeFollow {
					replyMessage, err := Follow(profile)
					if err != nil {
						logger.Write(err)
					}

					if _, err = bot.PushMessage(userID, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						logger.Write(err)
					}
				}
			}

			// ブロック処理時はプロフィールを取得できないので、if文の外に記載
			if event.Type == linebot.EventTypeUnfollow {

				// ユーザ情報をDBから削除
				selector := bson.M{"userid": userID}
				if err := mongo.RemoveDb(selector, "userInfos"); err != nil {
					logger.Write(err)
				} else {
					logger.Write("success delete:" + userID)
				}
			}

			mongo.DisconnectDb()
			logger.Write("end event")
		}
	})

	// 使用するポートを取得
	var addr = flag.String("addr", ":443", "アプリケーションのアドレス")
	flag.Parse()

	logger.Write("start server linebot port", *addr)

	http.Handle("/callback", handler)
	if err := http.ListenAndServeTLS(*addr, apiIDs.CertFile, apiIDs.KeyFile, nil); err != nil {
		logger.Fatal("ListenAndServe: ", err)
	}

	// if err := http.ListenAndServe(":80", nil); err != nil {
	// 	logger.Fatal(err)
	// }

}
