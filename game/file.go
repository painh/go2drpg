package game

import (
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"io/ioutil"
	"log"
)

//func ReadFile(filename string) ([]byte, error) {
//	statikFS, err := fs.New()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Access individual files by their paths.
//	r, err := statikFS.Open(filename)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer r.Close()
//	contents, err := ioutil.ReadAll(r)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return contents, err
//}

func ReadFile(filename string) ([]byte, error) {
	r, err := ebitenutil.OpenFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	return contents, err
}
