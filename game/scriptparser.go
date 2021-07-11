package game

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"log"
	"path/filepath"
	"strings"
)

type ScriptParser struct {
}

func ScriptParse(scr *SceneManager, i interface{}) ScriptActionInterface {
	m := i.(map[interface{}]interface{})

	for k2, v2 := range m {
		switch strings.ToLower(k2.(string)) {
		case "text":
			return &ScriptActionText{text: v2.(string)}
		case "setstatus":
			return &ScriptActionSetGameStatus{status: v2.(int)}
		case "person":
			scr.person = v2.(bool)
		case "nonexclusive":
			scr.nonexclusive = v2.(bool)
		case "switchon":
			return &ScriptActionSetSwitch{keyword: v2.(string), flag: true}
		case "switchoff":
			return &ScriptActionSetSwitch{keyword: v2.(string), flag: false}
		case "addkeyword":
			return &ScriptActionAddKeyword{keyword: v2.(string)}
		case "addlocation":
			return &ScriptActionAddLocation{keyword: v2.(string)}
		case "addperson":
			return &ScriptActionAddPerson{keyword: v2.(string)}
		case "playmusic":
			return &ScriptActionPlayMusic{filename: v2.(string)}
		case "condition":
			scr.condition = v2.([]interface{})
		case "invalidkeywordresponse":
			scr.invalidKeywordResponse = v2.(string)

		default:
			log.Fatal("invalid command :", k2, ": ", v2)
		}
	}

	//return m.(*ScriptActionInterface)
	return nil
}

func ScriptLoad(filename string) {
	m := []map[interface{}]interface{}{}

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

	for _, v := range m {
		for k2, v2 := range v {
			//scr := GameInstance.scriptManager.GetSceneManager(k2.(string))
			scr := GameInstance.scriptManager.NewSceneManager(k2.(string))

			for _, v2 := range v2.([]interface{}) {
				scene := ScriptParse(scr, v2)
				if scene != nil {
					scr.scene = append(scr.scene, scene)
				}
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
