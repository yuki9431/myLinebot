package main

import (
	"math/rand"
	"time"

	"github.com/yuki9431/mongohelper"
)

// HinaResponce セリフをランダムに返す
func HinaResponce() (replyMessage string, err error) {

	// TODO DBからランダムでセリフを取得
	mongo, err := mongohelper.NewMongo(mongoDial, mongoName)
	defer mongo.DisconnectDb()

	// ドキュメント数を取得し、ランダムでセリフをDB取得する
	// words, err := mongo.Count("quotes")
	// if err != nil {
	// 	return
	// }

	//mongo.SearchDb

	// セリフ一覧
	replyMessages := []string{
		"ねえ、今から晴れるよ",
		"ねえ、今からキレるよ？",
		"信じられない。気持ち悪い。最悪",
		"あの日私達は世界の形を決定的に変えてしまったんだ",
		"私、好きだな、この仕事。晴れ女の仕事",
	}

	// 乱数seedを設定
	rand.Seed(time.Now().UnixNano())

	// メッセージ数を取得
	len := len(replyMessages)

	// メッセージをランダムで返す
	return replyMessages[rand.Intn(len)], nil
}
