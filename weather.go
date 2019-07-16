package main

import (
	"errors"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/yuki9431/weather"
)

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
	w, err := weather.New(cityId, appId)
	w.SetTimezone(*time.FixedZone("Asia/Tokyo", 9*60*60))
	w.Infos = *w.GetInfoFromDate(time.Now())
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
	mongo, err := NewMongo(mongoDial, mongoName)

	for {
		t := time.Now()
		if t.Format(layout) == "06:00:00" {
			// DBからユーザ情報を取得
			if err = mongo.searchDb(userinfos, "userInfos"); err != nil {
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
