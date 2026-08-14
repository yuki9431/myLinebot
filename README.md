weather
====
## Overview

- [OpenWeather](https://openweathermap.org)のAPIをGo言語で使用するためのパッケージ

## Description
5日分の天気情報を3時間ごとに取得できます。

### 天気情報
- 都市名	string型
- 日付	[]string型
- 天気	[]string型
- アイコン	[]string型
- 気温	[]int型

## Requirement
- Go 1.10 or later


## Install
```bash:#
go get github.com/yuki9431/weather
```

## Configuration
```go:main.go
import (
	"fmt"

	"github.com/yuki9431/weather"
)

func main() {
	w, err := weather.New("<cityid>", "<appid>")
	if err != nil {
		fmt.Println(err)
	}

	for i, date := range w.GetDates() {
		fmt.Println(date.String() + " : " + w.GetDescriptions()[i])
	}
}
```

## How to start

はじめに[OpenWeather](https://openweathermap.org)のアカウントを作成し  
API key(APPID)を取得する必要があります。

1. New(cityId, appid)でweatherInfos型を取得する
2. 下記メソッドで必要な情報を取得


### New(cityId string, appid string)
OpenWeatherから5日間の天気情報を取得します。  
天気情報はweatherInfos型で取得できます。

### GetCityName()
New()で取得した天気情報に含まれる都市名を取得します。  
(TODO 現在は東京で固定しているが、可変にしたい)

### GetIcons()
New()で取得した天気情報に含まれる天気アイコンを表す文字列を取得します。  
5日分*3時間毎の天気アイコンを[]string型で返します。

### GetDates()
New()で取得した天気情報に含まれる日付を取得します。  
5日分*3時間毎の日付を[]time.Time型で返します。

### GetDescriptions()
New()で取得した天気情報に含まれる天気(sun, rain等)を取得します。  
5日分*3時間毎の天気を[]string型で返します。

### GetTemps()
New()で取得した天気情報に含まれる気温を取得します。  
小数点以下はの値は四捨五入します。  
5日分*3時間毎の気温を[]int型で返します。

### ConvertIconToWord(icon string)
天気アイコンを絵文字変換します。  
(01d ⇒ ☀️)

### GetInfoFromDate(target time.Time)
日付を指定し、1日分の天気情報を取得できます。  



## Contribution
1. Fork ([https://github.com/yuki9431/weather](https://github.com/yuki9431/weather))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Create new Pull Request


## Author
[Dillen H. Tomida](https://github.com/yuki9431)
