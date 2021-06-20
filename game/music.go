package game

import (
	"bytes"
	_ "image/png"
	"io"
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

type MusicPlayer struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seCh         chan []byte
	volume128    int
	musicType    musicType
}

func NewPlayer(audioContext *audio.Context, filename string) (*MusicPlayer, error) {
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
	player := &MusicPlayer{
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		volume128:    128,
		seCh:         make(chan []byte),
		musicType:    typeOgg,
	}

	if player.total == 0 {
		player.total = 1
	}
	player.audioPlayer.Play()

	return player, nil
}

func (p *MusicPlayer) Close() error {
	return p.audioPlayer.Close()
}

func (p *MusicPlayer) update() error {
	if p.audioPlayer.IsPlaying() {
		p.current = p.audioPlayer.Current()
	}
	return nil
}

type MusicManager struct {
	audioContext  *audio.Context
	musicPlayer   *MusicPlayer
	musicPlayerCh chan *MusicPlayer
	errCh         chan error
	playNum       int
}

func (m *MusicManager) Init() {
	if m.audioContext != nil {
		return
	}
	m.audioContext = audio.NewContext(sampleRate)

	//m.PlayNum(0)
	//newPlayer, _ := NewPlayer(audioContext, musicList[0])
	//m.musicPlayer = newPlayer
	//m.musicPlayerCh = make(chan *MusicPlayer)
	//m.errCh = make(chan error)
}

func (m *MusicManager) Play(filename string) {
	if m.musicPlayer != nil {
		m.musicPlayer.Close()
	}

	newPlayer, _ := NewPlayer(m.audioContext, SettingConfigInstance.WorkFolder+filename)
	m.musicPlayer = newPlayer
}

func (m *MusicManager) Update() error {
	//select {
	//case p := <-m.musicPlayerCh:
	//	m.musicPlayer = p
	//case err := <-m.errCh:
	//	return err
	//default:
	//}
	//
	if m.musicPlayer != nil {
		if err := m.musicPlayer.update(); err != nil {
			return err
		}
	}

	return nil
}
