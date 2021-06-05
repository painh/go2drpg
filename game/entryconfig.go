package game

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"log"
	"path/filepath"
)

type (
	EntryConfig struct {
		WindowWidth            int        `yaml:"window_width" json:"window_width"`
		WindowHeight           int        `yaml:"window_height" json:"window_height"`
		Title                  string     `yaml:"title" json:"title"`
		DoubleClickPixelMargin int        `yaml:"double_click_pixel_margin" json:"double_click_pixel_margin"`
		DoubleClickTsMargin    int        `yaml:"double_click_ts_margin" json:"double_click_ts_margin"`
		FontPath               string     `yaml:"font_path" json:"font_path"`
		FontSize               int        `yaml:"font_size" json:"font_size"`
		Scenario               []Scenario `yaml:"scenario" json:"scenario"`
	}

	Scenario struct {
		Name   string `yaml:"name" json:"name"`
		Config string `yaml:"config" json:"config"`
	}
)

var EntryConfigInstance = EntryConfig{}

func (c *EntryConfig) Load(filename string) {
	b, err := ReadFile(filename) // articles.json 파일의 내용을 읽어서 바이트 슬라이스에 저장
	if err != nil {
		log.Fatalln(err)
		return
	}

	var name = filepath.Ext(filename)

	if name == ".json" {
		err = json.Unmarshal(b, c)
	} else {
		err = yaml.Unmarshal([]byte(b), c)
	}

	if err != nil {
		log.Fatalln(err)
		return
	}
}
