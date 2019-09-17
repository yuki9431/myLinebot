package main

import (
	"errors"
	"linebot/config"
	"strconv"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/yuki9431/mongoHelper"
	"github.com/yuki9431/weather"
)

// 天気情報を日本語に変換
func convertWeatherToJp(description string) (jpDescription string) {
	if description != "" {
		switch {
		case description == "clear sky":
			jpDescription = "快晴"

		case description == "few clouds":
			jpDescription = "晴れ"

		case description == "rain":
			jpDescription = "雨"

		case description == "light rain":
			jpDescription = "小雨"

		case strings.Contains(description, "rain"):
			jpDescription = "雨"

		case strings.Contains(description, "cloud"):
			jpDescription = "曇り"

		case strings.Contains(description, "snow"):
			jpDescription = "雪"

		case strings.Contains(description, "thunderstorm"):
			jpDescription = "雷雨"

		default:
			jpDescription = description

		}
	}
	return
}

// 天気情報作成
func createWeatherMessage(apiIds *ApiIds, userInfo UserInfo, cityInfo CityInfo) (message string, err error) {
	// 設定ファイル読み込み
	config := config.NewConfig(configFile)
	if err = config.Read(apiIds); err != nil {
		err = errors.New("err :faild read config")
		return
	}

	cityId := userInfo.CityId
	appId := apiIds.AppId

	// 今日の天気情報を取得　今日の天気情報がない場合は、翌日の天気を取得(0時に近い時を想定)
	w, err := weather.New(cityId, appId)
	w.SetTimezone(time.FixedZone("Asia/Tokyo", 9*60*60))
	if todayInfo := w.GetInfoFromDate(time.Now()); todayInfo.List != nil {
		w.Infos = *todayInfo

	} else {
		todayInfo := w.GetInfoFromDate(time.Now().Add(24 * time.Hour))
		w.Infos = *todayInfo
	}
	dates := w.GetDates()
	icons := w.GetIcons()
	temps := w.GetTemps()
	cityName := cityInfo.Name // APIだと英語表記になるのでDBから持ってくる
	descriptions := w.GetDescriptions()

	// 天気情報メッセージ作成
	message = cityName + "\n" +
		func() string {
			var tempIcon string

			for i, time := range dates {
				tempIcon += time.Format("15:04") + " " +
					convertWeatherToJp(descriptions[i]) +
					w.ConvertIconToWord(icons[i]) + "  " +
					strconv.Itoa(temps[i]) + "℃" + "\n"
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
	userinfos := new([]UserInfo)
	cityInfos := new([]CityInfo)
	mongo, err := mongoHelper.NewMongo(mongoDial, mongoName)

	for {
		t := time.Now()
		if t.Format(layout) == "06:00:00" {
			// DBからユーザ情報を取得
			if err = mongo.SearchDb(userinfos, nil, "userInfos"); err != nil {
				return
			}

			// DBから都市情報を取得
			if err = GetCityInfo(cityInfo, )

			// 抽出した全ユーザ情報に天気情報を配信
			for _, userinfo := range *userinfos {
				if userinfo.UserID != "" {
					var bot *linebot.Client
					if bot, err = linebot.New(apiIds.ChannelSecret, apiIds.ChannelToken); err != nil {
						// エラー時はその場で終了
						return

					} else {
						// 天気情報メッセージ送信
						var message string
						message, err = createWeatherMessage(apiIds, userinfo)
						_, err = bot.PushMessage(userinfo.UserID, linebot.NewTextMessage(message)).Do()
					}
				} else {
					err = errors.New("Error: userId取得失敗")
					return
				}
			}

			mongo.DisconnectDb()

			// 連続送信を防止する
			time.Sleep(1 * time.Second) // sleep 1 second
		}
	}
}
