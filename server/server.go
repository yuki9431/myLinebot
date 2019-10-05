package main

import (
	"linebot/config"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/greymd/ojichat/generator"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"github.com/yuki9431/logger"
	"github.com/yuki9431/mongoHelper"
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
	apiIds := new(ApiIds)
	config := config.NewConfig(configFile)
	if err := config.Read(apiIds); err != nil {
		logger.Fatal(err)
	}

	// 開始メッセージ
	logger.Write("start server linebot")

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
			// DB設定
			mongo, err := mongoHelper.NewMongo(mongoDial, mongoName)
			if err != nil {
				logger.Fatal(err)
			}

			logger.Write("start event : " + event.Type)

			// ユーザのIDを取得
			userId := event.Source.UserID
			logger.Write("userid :" + userId)

			// 都市IDを取得するため、DBからユーザ情報を獲得
			userInfos := new([]UserInfo)
			if err := mongo.SearchDb(userInfos, bson.M{"userid": userId}, "userInfos"); err != nil {
				logger.Write("err search userInfo" + err.Error())
				return
			}

			// APIからユーザのプロフィールを取得後、レスポンスする
			if profile, err := bot.GetProfile(userId).Do(); err == nil {
				if event.Type == linebot.EventTypeMessage {
					// 返信メッセージ
					var replyMessage string

					switch message := event.Message.(type) {
					case *linebot.TextMessage:
						if strings.Contains(message.Text, "天気") {
							if replyMessage, err = createWeatherMessage(apiIds, (*userInfos)[0]); err != nil { // (*userInfos)[0]は一意の値しか取れない想定
								logger.Write(err)
							}

						} else if strings.Contains(message.Text, "おじさん") || strings.Contains(message.Text, "オジサン") {
							if replyMessage, err = ojichat(profile.DisplayName); err != nil {
								logger.Write(err)
							}

						} else if strings.Contains(message.Text, "都市変更") {
							cityName := strings.Replace(message.Text, " ", "", -1) // 全ての半角スペースを消す
							cityName = strings.Replace(cityName, "都市変更:", "", 1)   // 頭の都市変更:を消す

							// 都市IDを取得する
							cityId, err := GetCityId(cityName)
							if err != nil {
								logger.Write(err)
							}

							// 都市IDをDBに登録する
							if cityId != "" {
								selector := bson.M{"userid": profile.UserID}
								update := bson.M{"$set": bson.M{"cityid": cityId}}
								if err := mongo.UpdateDb(selector, update, "userInfos"); err != nil {
									replyMessage = "都市の変更に失敗しました..."
									logger.Write("failed update ciyId")

								} else {
									replyMessage = "選択された都市に変更しました！"
									logger.Write("success update ciyId")
								}
							} else {
								replyMessage = "該当都市がな見つかりません💦\n" +
									"都市一覧と送り頂ければ設定可能な都市が表示されますよ"
							}

						} else if strings.Contains(message.Text, "都市一覧") {
							cityList := new([]string)
							GetAllCityList(cityList)

							replyMessage = "都市一覧\n"
							for _, city := range *cityList {
								replyMessage = replyMessage + city + "\n"
							}

						} else if strings.Contains(message.Text, "ヘルプ") || strings.Contains(message.Text, "help") {
							// botの機能を返信する
							replyMessage = usage
						} else {
							// 100%の晴れ女
							replyMessage = HinaResponce()
						}

						// 返信処理
						if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
							logger.Write(err)
						}
						logger.Write("message.Text: " + message.Text)
					}
				} else if event.Type == linebot.EventTypeFollow {
					userInfo := new(UserInfo)
					userInfo.UserID = profile.UserID
					userInfo.DisplayName = profile.DisplayName
					userInfo.CityId, _ = GetCityId("東京") //初回登録時には問答無用で東京民や
					userInfo.PictureURL = profile.PictureURL
					userInfo.StatusMessage = profile.StatusMessage

					// ユーザ情報をDBに登録
					if err := mongo.InsertDb(userInfo, "userInfos"); err != nil {
						logger.Write(err)
					}

					// フレンド登録時の挨拶
					var replyMessages [5]string
					replyMessages[0] = profile.DisplayName + "さん\nはじめまして、毎朝6時に天気情報を教えてあげるね"
					replyMessages[1] = usage
					replyMessages[2] = "お住まいの都市を変更するには、下記の通りメッセージをお送りください"
					replyMessages[3] = "都市変更:東京"
					replyMessages[4] = "都市変更:Brasil"

					for _, replyMessage := range replyMessages {
						if _, err = bot.PushMessage(userId, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
							logger.Write(err)
						}
					}

				}
			}

			// ブロック処理時はプロフィールを取得できないので、if文の外に記載
			if event.Type == linebot.EventTypeUnfollow {

				// ユーザ情報をDBから削除
				selector := bson.M{"userid": userId}
				if err := mongo.RemoveDb(selector, "userInfos"); err != nil {
					logger.Write(err)
				} else {
					logger.Write("success delete:" + userId)
				}
			}

			mongo.DisconnectDb()
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
