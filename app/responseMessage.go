package main

import "github.com/yuki9431/mongohelper"

// HinaResponce セリフをランダムに返す
func HinaResponce() (replyMessage string, err error) {
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	// DBからセリフをランダムで取得する 1つだけ取得する想定
	var quotos []struct{ Quoto string }
	err = mongo.RandomSearchDb(&quotos, "quotos")

	return quotos[0].Quoto, err
}

// MorningGreeting 朝の挨拶を返事をDBから取得する
func MorningGreeting() (replyMessage string, err error) {
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	// DBからセリフをランダムで取得する 1つだけ取得する想定
	var quotos []struct{ Quoto string }
	err = mongo.RandomSearchDb(&quotos, "morningGreeting")

	return quotos[0].Quoto, err
}

// NoonGreeting 昼の挨拶を返事をDBから取得する
func NoonGreeting() (replyMessage string, err error) {
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	// DBからセリフをランダムで取得する 1つだけ取得する想定
	var quotos []struct{ Quoto string }
	err = mongo.RandomSearchDb(&quotos, "noonGreeting")

	return quotos[0].Quoto, err
}

// NightGreeting 夜の挨拶を返事をDBから取得する
func NightGreeting() (replyMessage string, err error) {
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	// DBからセリフをランダムで取得する 1つだけ取得する想定
	var quotos []struct{ Quoto string }
	err = mongo.RandomSearchDb(&quotos, "nightGreeting")

	return quotos[0].Quoto, err
}
