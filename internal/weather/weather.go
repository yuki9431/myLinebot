package weather

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"time"
)

// 絶対零度 気温変換で使用する
const absoluteTmp = -273.15

type weather struct {
	cityId   string
	appid    string
	timezone *time.Location
	Infos    weatherInfos
}

type weatherInfos struct {
	List []struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		Name string `json:"name"`
	} `json:"city"`
}

func New(cityId string, appid string) (w *weather, err error) {
	apiUrl := "http://api.openweathermap.org/data/2.5/forecast?id=" +
		cityId +
		"&" +
		"appid=" +
		appid

	resp, err := http.Get(apiUrl)
	if err != nil {
		err = errors.New("天気情報の取得に失敗しました")
		return
	}
	defer resp.Body.Close()

	// jsonデコード
	var wI weatherInfos
	if err = json.NewDecoder(resp.Body).Decode(&wI); err != nil {
		err = errors.New("jsonデコードに失敗しました")
	}

	timezone := time.FixedZone("JST", 9*60*60)

	w = &weather{
		cityId:   cityId,
		appid:    appid,
		Infos:    wI,
		timezone: timezone,
	}

	return
}

func (w *weather) SetTimezone(t *time.Location) {
	w.timezone = t
}

func (w *weather) GetCityName() string {
	return w.Infos.City.Name
}

func (w *weather) GetIcons() []string {
	var icons []string
	for _, l := range w.Infos.List {
		for _, lw := range l.Weather {
			icons = append(icons, lw.Icon)
		}
	}

	return icons
}

func (w *weather) GetDates() []time.Time {
	var times []time.Time

	for _, l := range w.Infos.List {
		date, _ := time.Parse("2006-01-02 15:04:05", l.DtTxt)
		if w.timezone != nil {
			times = append(times, date.In(w.timezone))
		} else {
			times = append(times, date)
		}
	}

	return times
}

func (w *weather) GetDescriptions() []string {
	var descriptions []string
	for _, l := range w.Infos.List {
		for _, w := range l.Weather {
			descriptions = append(descriptions, w.Description)
		}
	}

	return descriptions
}

func (w *weather) GetTemps() []int {
	var maxTemps []int
	for _, l := range w.Infos.List {
		maxTemps = append(maxTemps, (int)(math.Round(l.Main.Temp+absoluteTmp)))
	}

	return maxTemps
}

func (w *weather) ConvertIconToWord(icon string) string {
	var word string

	switch icon {
	case "01d", "01n":
		word = "☀️"
	case "02d", "02n":
		word = "🌤"
	case "03d", "04d", "03n", "04n":
		word = "☁️"
	case "09d", "09n":
		word = "☂️"
	case "10d", "10n":
		word = "☔️"
	case "11d", "11n":
		word = "⚡️"
	case "13d", "13n":
		word = "☃️"
	case "50d", "50n":
		word = "💨"
	default:
		word = "😇" // 不正な値

	}
	return word
}

func (w *weather) GetInfoFromDate(target time.Time) *weatherInfos {
	const layout = "2006-01-02" // => YYYY-MM-DD

	var weatherInfosToday weatherInfos
	weatherInfosToday.City.Name = w.Infos.City.Name

	for i, date := range w.GetDates() {
		if target.Format(layout) == date.Format(layout) {
			weatherInfosToday.List = append(weatherInfosToday.List, w.Infos.List[i])
		}
	}
	return &weatherInfosToday
}
