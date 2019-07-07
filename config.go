package main

import (
	"encoding/json"
	"io/ioutil"
)

// configファイルを読み込み構造体へ割当
func Read(obj interface{}, filename string) error {

	// 設定ファイルを読み込む
	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// 設定
	err = json.Unmarshal(jsonString, obj)
	if err != nil {
		return err
	}

	return nil
}
