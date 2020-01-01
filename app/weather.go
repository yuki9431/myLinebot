package main

import (
	"errors"
	"myLinebot/config"
	"strconv"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/yuki9431/mongohelper"
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
			jpDescription = "雨　"

		case description == "light rain":
			jpDescription = "小雨"

		case strings.Contains(description, "rain"):
			jpDescription = "雨　"

		case strings.Contains(description, "cloud"):
			jpDescription = "曇り"

		case strings.Contains(description, "snow"):
			jpDescription = "雪　"

		case strings.Contains(description, "thunderstorm"):
			jpDescription = "雷雨"

		default:
			jpDescription = description

		}
	}
	return
}

// CreateWeatherMessage 天気情報作成
func CreateWeatherMessage(userID string, apiIDs *APIIDs, target time.Time) (message string, err error) {

	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	if err != nil {
		return
	}

	// 都市IDを取得するため、DBからユーザ情報を獲得
	userInfos := new([]UserInfo)

	if err = mongo.SearchDb(userInfos, bson.M{"userid": userID}, "userInfos"); err != nil {
		err = errors.New("err search userInfo" + err.Error())
	}

	// 設定ファイル読み込み
	config := config.NewConfig(configFile)
	if err = config.Read(apiIDs); err != nil {
		err = errors.New("err :faild read config")
		return
	}

	cityID := (*userInfos)[0].CityID
	appID := apiIDs.AppID

	// APIだと英語表記になるのでDBから都市名を取得
	cityName, err := GetCityName(cityID)
	if err != nil {
		err = errors.New("err : faild get cityName")
		return
	}

	// 今日の天気情報を取得　今日の天気情報がない場合は、翌日の天気を取得(0時に近い時を想定)
	w, err := weather.New(cityID, appID)
	w.SetTimezone(time.FixedZone("Asia/Tokyo", 9*60*60))
	if todayInfo := w.GetInfoFromDate(target); todayInfo.List != nil {
		w.Infos = *todayInfo

	} else {
		todayInfo := w.GetInfoFromDate(target.Add(24 * time.Hour))
		w.Infos = *todayInfo
	}
	dates := w.GetDates()
	icons := w.GetIcons()
	temps := w.GetTemps()
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

// CreateWeekWeatherMessage 1週間分の天気情報作成
func CreateWeekWeatherMessage(userID string, apiIDs *APIIDs) (message string, err error) {

	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	if err != nil {
		return
	}

	// 都市IDを取得するため、DBからユーザ情報を獲得
	userInfos := new([]UserInfo)

	if err = mongo.SearchDb(userInfos, bson.M{"userid": userID}, "userInfos"); err != nil {
		err = errors.New("err search userInfo" + err.Error())
	}

	// 設定ファイル読み込み
	config := config.NewConfig(configFile)
	if err = config.Read(apiIDs); err != nil {
		err = errors.New("err :faild read config")
		return
	}

	cityID := (*userInfos)[0].CityID
	appID := apiIDs.AppID

	// APIだと英語表記になるのでDBから都市名を取得
	cityName, err := GetCityName(cityID)
	if err != nil {
		err = errors.New("err : faild get cityName")
		return
	}

	// 天気情報を取得
	w, err := weather.New(cityID, appID)
	w.SetTimezone(time.FixedZone("Asia/Tokyo", 9*60*60))

	// 曜日表示の設定
	wdays := [...]string{"日", "月", "火", "水", "木", "金", "土"}

	message = cityName + "\n" + "1週間の天気情報だよ\n"

	var yesterday time.Time
	for i, date := range w.GetDates() {

		//日付を跨いだタイミングで改行するし、日付を挿入する
		if yesterday.Format("2006-01-02") != date.Format("2006-01-02") {
			message += "\n" + date.Format("01/02 (") + wdays[date.Weekday()] + ")\n"
		}

		message += date.Format("15:04 ") +
			convertWeatherToJp(w.GetDescriptions()[i]) +
			w.ConvertIconToWord(w.GetIcons()[i]) + "  " +
			strconv.Itoa(w.GetTemps()[i]) + "℃" + "\n"

		// 前日情報と比較するため、前日情報を保持する
		yesterday = date
	}

	return
}

// SendWeatherInfo 朝の天気配信ジョブ
func SendWeatherInfo(apiIDs *APIIDs) (err error) {
	const layout = "15:04:05" // => hh:mm:ss
	userinfos := new([]UserInfo)
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	for {
		t := time.Now()
		if t.Format(layout) == "06:00:00" {
			// DBからユーザ情報を取得
			if err = mongo.SearchDb(userinfos, nil, "userInfos"); err != nil {
				return
			}

			// 抽出した全ユーザ情報に天気情報を配信
			for _, userinfo := range *userinfos {
				if userinfo.CityID != "" {
					var bot *linebot.Client
					if bot, err = linebot.New(apiIDs.ChannelSecret, apiIDs.ChannelToken); err == nil {
						// 天気情報メッセージ送信
						var message string
						message, err = CreateWeatherMessage(userinfo.UserID, apiIDs, time.Now())
						_, err = bot.PushMessage(userinfo.UserID, linebot.NewTextMessage(message)).Do()
					} else {
						// error
						return
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
