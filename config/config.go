package config

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	file string
}

type Config interface {
	Read(interface{}) error
}

func NewConfig(f string) Config {
	return &config{file: f}
}

// configファイルを読み込み構造体へ割当
func (c *config) Read(obj interface{}) error {

	// 設定ファイルを読み込む
	jsonString, err := ioutil.ReadFile(c.file)
	if err != nil {
		return err
	}

	// 設定
	err = json.Unmarshal(jsonString, &obj)
	if err != nil {
		return err
	}

	return nil
}
