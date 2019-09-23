package main

import (
	"math/rand"
	"time"
)

// セリフをランダムに返す
func HinaResponce() (replyMessage string) {

	// セリフ一覧
	replyMessages := []string{
		"ねえ、今から晴れるよ",
		"ねえ、今からキレるよ？",
		"信じられない。気持ち悪い。最悪",
	}

	// 乱数seedを設定
	rand.Seed(time.Now().UnixNano())

	// メッセージ数を取得
	len := len(replyMessages)

	// メッセージをランダムで返す
	return replyMessages[rand.Intn(len)]
}
