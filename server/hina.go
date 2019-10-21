package main

import (
	"github.com/yuki9431/mongohelper"
)

// HinaResponce セリフをランダムに返す
func HinaResponce() (replyMessage string, err error) {
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	// DBからセリフをランダムで取得する 1つだけ取得する想定
	quoto := new([]string)
	err = mongo.RandomSearchDb(quoto, nil, "quotos")

	return (*quoto)[1], err
}
