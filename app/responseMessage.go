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

// morningGreeting 朝の挨拶を返事をDBから取得する
func morningGreeting() (replyMessage string, err error) {
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	// DBからセリフをランダムで取得する 1つだけ取得する想定
	var quotos []struct{ Quoto string }
	err = mongo.RandomSearchDb(&quotos, "morningGreeting")

	return quotos[0].Quoto, err
}

// noonGreeting 昼の挨拶を返事をDBから取得する
func noonGreeting() (replyMessage string, err error) {
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	// DBからセリフをランダムで取得する 1つだけ取得する想定
	var quotos []struct{ Quoto string }
	err = mongo.RandomSearchDb(&quotos, "noonGreeting")

	return quotos[0].Quoto, err
}

// nightGreeting 夜の挨拶を返事をDBから取得する
func nightGreeting() (replyMessage string, err error) {
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	// DBからセリフをランダムで取得する 1つだけ取得する想定
	var quotos []struct{ Quoto string }
	err = mongo.RandomSearchDb(&quotos, "nightGreeting")

	return quotos[0].Quoto, err
}
