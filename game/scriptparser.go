package game

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"log"
	"path/filepath"
)

type ScriptParser struct {
}

func ScriptParse(scr *SceneManager, i interface{}) ScriptActionInterface {
	m := i.(map[interface{}]interface{})

	for k2, v2 := range m {
		switch k2.(string) {
		case "text":
			return &ScriptActionText{text: v2.(string)}
		case "setStatus":
			return &ScriptActionSetGameStatus{status: v2.(int)}
		case "person":
			scr.person = v2.(bool)
		case "switchon":
			return &ScriptActionSetSwitch{keyword: v2.(string), flag: true}
		case "switchoff":
			return &ScriptActionSetSwitch{keyword: v2.(string), flag: false}
		case "addkeyword":
			return &ScriptActionKeyword{keyword: v2.(string)}
		case "playmusic":
			return &ScriptActionPlayMusic{filename: v2.(string)}
		default:
			log.Fatal("invalid command ", k2, v2)
		}
	}

	//return m.(*ScriptActionInterface)
	return nil
}

func ScriptLoad(filename string) {
	m := make(map[interface{}]interface{})

	b, err := ReadFile(filename) // articles.yaml:파일의 내용을  json 파일의 내용을 읽어서 바이트 슬라이스에 저장
	if err != nil {
		log.Fatalln(err)
		return
	}

	var name = filepath.Ext(filename)

	if name == ".json" {
		err = json.Unmarshal(b, &m)
	} else {
		err = yaml.Unmarshal([]byte(b), &m)
	}

	if err != nil {
		log.Fatalln(err)
		return
	}

	for k, v := range m {
		scr := GameInstance.scriptManager.GetSceneManager(k.(string))

		for _, v2 := range v.([]interface{}) {
			scene := ScriptParse(scr, v2)
			if scene != nil {
				scr.scene = append(scr.scene, scene)
			}
		}
	}

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err != nil {
		log.Fatalln(err)
		return
	}

}
