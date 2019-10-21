package main

import (
	"github.com/yuki9431/mongohelper"
)

// HinaResponce セリフをランダムに返す
func HinaResponce() (replyMessage string, err error) {
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	// DBからセリフをランダムで取得する 1つだけ取得する想定
	var quotos []struct{ Quoto string }
	err = mongo.RandomSearchDb(&quotos, "quotos")

	return quotos[0].Quoto, err
}
