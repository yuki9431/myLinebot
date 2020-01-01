package main

import (
	"strings"
)

// IsOjichan おじいさんの呼びかけに反応する
func IsOjichan(m string) (isOjichan bool) {
	messages := []string{
		"おじいさん",
		"オジイサン",
		"おじいちゃん",
		"オジイチャン",
		"おじさん",
		"オジサン",
		"おじちゃん",
		"オジチャン",
		"じじい",
		"ジジイ",
	}

	return contains(m, messages)
}

// IsAskWeather 天気を尋ねられたとき
func IsAskWeather(m string) (isAskWeather bool) {
	messages := []string{
		"天気",
		"てんき",
		"テンキ",
		"暑い",
		"寒い",
		"あつい",
		"さむい",
	}

	return contains(m, messages)
}

// IsAskTomorrowWeather 明日の天気を教えるとき
func IsAskTomorrowWeather(m string) (isAskTomorrowWeather bool) {
	messages := []string{
		"あした",
		"アシタ",
		"明日",
		"翌日",
		"tomorrow",
		"Tomorrow",
	}

	return contains(m, messages)
}

// IsAskWeekWeather 週の天気を教えるとき
func IsAskWeekWeather(m string) (isAskTomorrowWeather bool) {
	messages := []string{
		"明後日",
		"あさって",
		"アサッテ",
		"未来",
		"週",
	}

	return contains(m, messages)
}

// IsMorningGreeting 朝の挨拶を返事するとき
func IsMorningGreeting(m string) (isMorningGreeting bool) {
	messages := []string{
		"おはよ",
		"はろ",
		"ハロ",
		"hi",
		"Hi",
		"hello",
		"Hello",
		"morning",
		"Morning",
		"oi",
		"Oi",
		"dia",
		"Dia",
		"朝",
	}

	return contains(m, messages)
}

// IsNoonGreeting 朝の挨拶を返事するとき
func IsNoonGreeting(m string) (isNoonGreeting bool) {
	messages := []string{
		"こんにち",
		"noon",
		"Noon",
		"tarde",
		"Tarde",
		"昼",
	}

	return contains(m, messages)
}

// IsNightGreeting 朝の挨拶を返事するとき
func IsNightGreeting(m string) (isNightGreeting bool) {
	messages := []string{
		"こんばん",
		"night",
		"Night",
		"noite",
		"Noite",
		"おやすみ",
		"夜",
	}

	return contains(m, messages)
}

// IsChangeCity 所在地を変更するとき
func IsChangeCity(m string) (isChangeCity bool) {
	messages := []string{
		"都市変更:",
	}

	return contains(m, messages)
}

// IsShowCityList 都市一覧を表示するとき
func IsShowCityList(m string) (isShowCityList bool) {
	messages := []string{
		"都市一覧",
	}

	return contains(m, messages)
}

// IsShowHelp 都市一覧を表示するとき
func IsShowHelp(m string) (isShowCityList bool) {
	messages := []string{
		"へるぷ",
		"ヘルプ",
		"助けて",
		"たすけて",
		"Help",
		"help",
	}

	return contains(m, messages)
}

// 文字列リストと引数の文字列が一致するか判別する
func contains(s string, strs []string) (isContainable bool) {
	isContainable = false

	for _, str := range strs {
		if strings.Contains(s, str) {
			isContainable = true
		}
	}
	return isContainable
}
