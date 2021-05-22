package game

import (
	"encoding/json"
	"log"
)

type (
	EntryConfig struct {
		WindowWidth            int        `json:"window_width"`
		WindowHeight           int        `json:"window_height"`
		Title                  string     `json:"title"`
		DoubleClickPixelMargin int        `json:"double_click_pixel_margin"`
		DoubleClickTsMargin    int        `json:"double_click_ts_margin"`
		FontPath               string     `json:"font_path"`
		FontSize               int        `json:"font_size"`
		Scenario               []Scenario `json:"scenario"`
	}

	Scenario struct {
		Name   string `json:"name"`
		Config string `json:"config"`
	}
)

var EntryConfigInstance = EntryConfig{}

func (c *EntryConfig) Load(filename string) {
	b, err := ReadFile(filename) // articles.json 파일의 내용을 읽어서 바이트 슬라이스에 저장
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = json.Unmarshal(b, c)

	if err != nil {
		log.Fatalln(err)
		return
	}
}
