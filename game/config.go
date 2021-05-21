package game

import (
	"encoding/json"
	"log"
)

type (
	Config struct {
		WindowWidth            int    `json:"window_width"`
		WindowHeight           int    `json:"window_height"`
		DoubleClickPixelMargin int    `json:"double_click_pixel_margin"`
		DoubleClickTsMargin    int64  `json:"double_click_ts_margin"`
		Title                  string `json:"title"`
		FontPath               string `json:"font_path"`
		FontSize               int    `json:"font_size"`
		LogLines               int    `json:"log_lines"`
		MapX                   int    `json:"map_x"`
		MapY                   int    `json:"map_y"`
		MapWidth               int    `json:"map_width"`
		MapHeight              int    `json:"map_height"`
		LogX                   int    `json:"log_x"`
		LogY                   int    `json:"log_y"`
		LogWidth               int    `json:"log_width"`
		LogHeight              int    `json:"log_height"`
		LineSpacing            int    `json:"line_spacing"`
		BtnPersonX             int    `json:"btn_person_x"`
		BtnPersonY             int    `json:"btn_person_y"`
		BtnPersonWidth         int    `json:"btn_person_width"`
		BtnPersonHeight        int    `json:"btn_person_height"`
		BtnLocationX           int    `json:"btn_location_x"`
		BtnLocationY           int    `json:"btn_location_y"`
		BtnLocationWidth       int    `json:"btn_location_width"`
		BtnLocationHeight      int    `json:"btn_location_height"`
		BtnItemX               int    `json:"btn_item_x"`
		BtnItemY               int    `json:"btn_item_y"`
		BtnItemWidth           int    `json:"btn_item_width"`
		BtnItemHeight          int    `json:"btn_item_height"`
		CursorX                int    `json:"cursor_x"`
		CursorY                int    `json:"cursor_y"`
		StartTimeMin           int    `json:"start_time_min"`
	}
)

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
