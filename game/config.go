package game

import (
	"encoding/json"
	"log"
)

type Config struct {
	Window_width           int    `json:window_width`
	Window_height          int    `json:window_height`
	Doubclick_pixel_margin int    `json:doubclick_pixel_margin`
	Doubclick_ts_margin    int64  `json:doubclick_ts_margin`
	Title                  string `json:title`
	Font_path              string `json:font_path`
	Font_size              int    `json:font_size`
	Log_lines              int    `json:log_lines`
}

var ConfigInstance = Config{}

func (c *Config) Load(filename string) {
	b, err := ReadFile(filename) // articles.json 파일의 내용을 읽어서 바이트 슬라이스에 저장
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = json.Unmarshal(b, &ConfigInstance)

	if err != nil {
		log.Fatalln(err)
		return
	}

}
