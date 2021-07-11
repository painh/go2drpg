package game

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const (
	sampleRate = 32000
)

type musicType int

const (
	typeOgg musicType = iota
	typeMP3
)

func (t musicType) String() string {
	switch t {
	case typeOgg:
		return "Ogg"
	case typeMP3:
		return "MP3"
	default:
		panic("not reached")
	}
}

type AudioPlayer struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	volume128    int
	musicType    musicType
}

func NewPlayer(audioContext *audio.Context, filename string) (*AudioPlayer, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	const bytesPerSample = 4 // TODO: This should be defined in audio package

	var s audioStream

	var err error

	b, err := ReadFile(filename)

	s, err = vorbis.Decode(audioContext, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	//https://titanwolf.org/Network/Articles/Article?AID=edc8b4a9-5d00-4326-838a-35e2234bd11c#gsc.tab=0
	//thanks
	l := audio.NewInfiniteLoop(s, s.Length())

	p, err := audio.NewPlayer(audioContext, l)
	if err != nil {
		return nil, err
	}
	player := &AudioPlayer{
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		volume128:    128,
		//seCh:         make(chan []byte),
		musicType: typeOgg,
	}

	if player.total == 0 {
		player.total = 1
	}
	player.audioPlayer.SetVolume(SettingConfigInstance.DefaultBGMVolume)
	player.audioPlayer.Play()

	return player, nil
}

func (p *AudioPlayer) Close() error {
	return p.audioPlayer.Close()
}

func (p *AudioPlayer) update() error {
	//select {
	//case p.seBytes = <-p.seCh:
	//	close(p.seCh)
	//	p.seCh = nil
	//default:
	//}

	//if p.audioPlayer.IsPlaying() {
	//	p.current = p.audioPlayer.Current()
	//}
	return nil
}

type AudioManager struct {
	audioContext              *audio.Context
	audioPlayer               *AudioPlayer
	musicPlayerCh             chan *AudioPlayer
	errCh                     chan error
	playNum                   int
	seDict                    map[string][]byte
	lastPlayLoopMusicFilename string
}

func (m *AudioManager) Init() {
	if m.audioContext != nil {
		return
	}
	m.audioContext = audio.NewContext(sampleRate)
	m.seDict = map[string][]byte{}

	//m.PlayNum(0)
	//newPlayer, _ := NewPlayer(audioContext, musicList[0])
	//m.audioPlayer = newPlayer
	//m.musicPlayerCh = make(chan *AudioPlayer)
	//m.errCh = make(chan error)
}

func (m *AudioManager) PlayLoopMusic(filename string) {
	if filename == m.lastPlayLoopMusicFilename {
		return
	}

	m.lastPlayLoopMusicFilename = filename

	if m.audioPlayer != nil {
		m.audioPlayer.Close()
	}

	newPlayer, _ := NewPlayer(m.audioContext, SettingConfigInstance.WorkFolder+filename)
	m.audioPlayer = newPlayer
}

func (m *AudioManager) Update() error {
	//select {
	//case p := <-m.musicPlayerCh:
	//	m.audioPlayer = p
	//case err := <-m.errCh:
	//	return err
	//default:
	//}
	//
	if m.audioPlayer != nil {
		if err := m.audioPlayer.update(); err != nil {
			return err
		}
	}

	return nil
}

func (m *AudioManager) PlayWave(filename string) {
	//
	//go func() {
	//	b, err := ReadFile(filename)
	//
	//	s, err := wav.Decode(m.audioContext, bytes.NewReader(b))
	//	if err != nil {
	//		log.Fatal(err)
	//		return
	//	}
	//
	//	b, err = ioutil.ReadAll(s)
	//	if err != nil {
	//		log.Fatal(err)
	//		return
	//	}
	//	m.audioPlayer.seCh <- b
	//}()

	v, ok := m.seDict[filename]

	if !ok {
		b, err := ReadFile(filename)

		s, err := wav.Decode(m.audioContext, bytes.NewReader(b))
		if err != nil {
			log.Fatal(err)
			return
		}

		v, err = ioutil.ReadAll(s)
		if err != nil {
			log.Fatal(err)
			return
		}
		m.seDict[filename] = v
	}

	sePlayer := audio.NewPlayerFromBytes(m.audioContext, v)
	sePlayer.SetVolume(SettingConfigInstance.DefaultSFXVolume)
	sePlayer.Play()
}
