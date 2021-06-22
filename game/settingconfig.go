package game

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"log"
	"path/filepath"
)

type (
	LocationInfo struct {
		Name        string `yaml:"name" json:"name"`
		DisplayName string `yaml:"display_name" json:"display_name"`
		Filename    string `yaml:"filename" json:"filename"`
	}
	SettingConfig struct {
		DevMode                bool           `yaml:"dev_mode" json:"dev_mode"`
		WindowWidth            int            `yaml:"window_width" json:"window_width"`
		WindowHeight           int            `yaml:"window_height" json:"window_height"`
		DoubleClickPixelMargin int            `yaml:"double_click_pixel_margin" json:"double_click_pixel_margin"`
		DoubleClickTsMargin    int64          `yaml:"double_click_ts_margin" json:"double_click_ts_margin"`
		RenderTileSize         float64        `yaml:"render_tile_size" json:"render_tile_size"`
		RealTileSize           float64        `yaml:"real_tile_size" json:"real_tile_size"`
		Title                  string         `yaml:"title" json:"title"`
		FontPath               string         `yaml:"font_path" json:"font_path"`
		FontSize               int            `yaml:"font_size" json:"font_size"`
		LogLines               int            `yaml:"log_lines" json:"log_lines"`
		MapX                   int            `yaml:"map_x" json:"map_x"`
		MapY                   int            `yaml:"map_y" json:"map_y"`
		MapWidth               int            `yaml:"map_width" json:"map_width"`
		MapHeight              int            `yaml:"map_height" json:"map_height"`
		LogX                   int            `yaml:"log_x" json:"log_x"`
		LogY                   int            `yaml:"log_y" json:"log_y"`
		LogWidth               int            `yaml:"log_width" json:"log_width"`
		LogHeight              int            `yaml:"log_height" json:"log_height"`
		LineSpacing            int            `yaml:"line_spacing" json:"line_spacing"`
		BtnPersonX             int            `yaml:"btn_person_x" json:"btn_person_x"`
		BtnPersonY             int            `yaml:"btn_person_y" json:"btn_person_y"`
		BtnPersonWidth         int            `yaml:"btn_person_width" json:"btn_person_width"`
		BtnPersonHeight        int            `yaml:"btn_person_height" json:"btn_person_height"`
		BtnLocationX           int            `yaml:"btn_location_x" json:"btn_location_x"`
		BtnLocationY           int            `yaml:"btn_location_y" json:"btn_location_y"`
		BtnLocationWidth       int            `yaml:"btn_location_width" json:"btn_location_width"`
		BtnLocationHeight      int            `yaml:"btn_location_height" json:"btn_location_height"`
		BtnKeywordX            int            `yaml:"btn_keyword_x" json:"btn_keyword_x"`
		BtnKeywordY            int            `yaml:"btn_keyword_y" json:"btn_keyword_y"`
		BtnKeywordWidth        int            `yaml:"btn_keyword_width" json:"btn_keyword_width"`
		BtnKeywordHeight       int            `yaml:"btn_keyword_height" json:"btn_keyword_height"`
		BtnTalkEndX            int            `yaml:"btn_talk_end_x" json:"btn_talk_end_x"`
		BtnTalkEndY            int            `yaml:"btn_talk_end_y" json:"btn_talk_end_y"`
		BtnTalkEndWidth        int            `yaml:"btn_talk_end_width" json:"btn_talk_end_width"`
		BtnTalkEndHeight       int            `yaml:"btn_talk_end_height" json:"btn_talk_end_height"`
		CursorX                int            `yaml:"cursor_x" json:"cursor_x"`
		CursorY                int            `yaml:"cursor_y" json:"cursor_y"`
		StartTimeMin           int            `yaml:"start_time_min" json:"start_time_min"`
		BtnZoomoutX            int            `yaml:"btn_zoomout_x" json:"btn_zoomout_x"`
		BtnZoomoutY            int            `yaml:"btn_zoomout_y" json:"btn_zoomout_y"`
		BtnZoomoutWidth        int            `yaml:"btn_zoomout_width" json:"btn_zoomout_width"`
		BtnZoomoutHeight       int            `yaml:"btn_zoomout_height" json:"btn_zoomout_height"`
		BtnZoominX             int            `yaml:"btn_zoomin_x" json:"btn_zoomin_x"`
		BtnZoominY             int            `yaml:"btn_zoomin_y" json:"btn_zoomin_y"`
		BtnZoominWidth         int            `yaml:"btn_zoomin_width" json:"btn_zoomin_width"`
		BtnZoominHeight        int            `yaml:"btn_zoomin_height" json:"btn_zoomin_height"`
		BtnCenterX             int            `yaml:"btn_center_x" json:"btn_center_x"`
		BtnCenterY             int            `yaml:"btn_center_y" json:"btn_center_y"`
		BtnCenterWidth         int            `yaml:"btn_center_width" json:"btn_center_width"`
		BtnCenterHeight        int            `yaml:"btn_center_height" json:"btn_center_height"`
		WorkFolder             string         `yaml:"work_folder" json:"work_folder"`
		LocationList           []LocationInfo `yaml:"location_list" json:"location_list"`
		DefaultMoveMin         int            `yaml:"default_move_min" json:"default_move_min"`
		DefaultTalkMin         int            `yaml:"default_talk_min" json:"default_talk_min"`
		DefaultLocationMin     int            `yaml:"default_location_min" json:"default_location_min"`
		PlayerObjectName       string         `yaml:"player_object_name" json:"player_object_name"`
		ZoomStep               int            `yaml:"zoom_step" json:"zoom_step"`
		Scripts                []string       `yaml:"scripts" json:"scripts"`
		BtnClickSound          string         `yaml:"btn_click_sound" json:"btn_click_sound"`
		DefaultBGMVolume       float64        `yaml:"default_bgm_volume" json:"default_bgm_volume"`
		DefaultSFXVolume       float64        `yaml:"default_sfx_volume" json:"default_sfx_volume"`
	}
)

var SettingConfigInstance = SettingConfig{}

func (c *SettingConfig) Load(filename string) {
	b, err := ReadFile(filename) // articles.yaml:파일의 내용을  json 파일의 내용을 읽어서 바이트 슬라이스에 저장
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
		log.Fatalf("error: %v", err)
	}

	if err != nil {
		log.Fatalln(err)
		return
	}
}
