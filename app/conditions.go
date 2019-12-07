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

// Contains 文字列リストと引数の文字列が一致するか判別する
func contains(s string, strs []string) (isContainable bool) {
	isContainable = false

	for _, str := range strs {
		if strings.Contains(s, str) {
			isContainable = true
		}
	}
	return isContainable
}
