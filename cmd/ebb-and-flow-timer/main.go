package main

import (
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Screen struct {
	time time.Time

	audioPlayer *audio.Player
}

//go:embed files
var fs embed.FS

func (s *Screen) Update() error {
	if time.Now().After(s.time) {
		s.time = s.time.Add(time.Minute * 2)
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
	}

	seconds := time.Now().Sub(s.time).Round(time.Second).Abs()
	if seconds == 15*time.Second && !s.audioPlayer.IsPlaying() { //Todo user variable for time until sound
		s.audioPlayer.Play()
		s.audioPlayer.Rewind()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		s.time = time.Now()
	}

	if isPressed(ebiten.KeyUp) {
		s.audioPlayer.SetVolume(s.audioPlayer.Volume() + .01)
	}
	if isPressed(ebiten.KeyDown) {
		s.audioPlayer.SetVolume(s.audioPlayer.Volume() - .01)
	}

	if isPressed(ebiten.KeyLeft) {
		s.time = s.time.Add(-time.Second)
	}
	if isPressed(ebiten.KeyRight) {
		s.time = s.time.Add(time.Second)
	}
	return nil
}

func isPressed(key ebiten.Key) bool {
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		return ebiten.IsKeyPressed(key)
	}
	return inpututil.IsKeyJustPressed(key)
}

func (s *Screen) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf(`Time left: %s
Volume: %.0f%%
Backspace: Reset time
<-: Remove one second
->: Add one second
^: increase volume
V: lower volume
Shift: Increase action
`, s.FormatTime(), s.audioPlayer.Volume()*100))
}

func (s *Screen) FormatTime() string {
	seconds := time.Now().Sub(s.time).Round(time.Second).Abs()

	minute := seconds / time.Minute
	seconds -= minute * time.Minute
	seconds = seconds / time.Second

	return fmt.Sprintf("%02d:%02d", minute, seconds)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't time have to adjust the screen size with the outside size, just return a fixed size.
func (s *Screen) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	audioContext := audio.NewContext(32000)

	audioFile, err := fs.Open("files/quack.mp3")
	if err != nil {
		log.Fatal(err)
	}
	stream, err := mp3.DecodeWithoutResampling(audioFile)
	if err != nil {
		log.Fatal(err)
	}

	audioPlayer, err := audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatal(err)
	}
	audioPlayer.SetVolume(0.2)

	s := &Screen{
		time:        time.Now().Add(time.Minute * 2),
		audioPlayer: audioPlayer,
	}

	ebiten.SetWindowSize(200, 200)
	ebiten.SetWindowTitle("Ebb and Flow Timer")

	if err := ebiten.RunGame(s); err != nil {
		log.Fatal(err)
	}

}
