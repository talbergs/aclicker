package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

//go:embed sprite.png
var sprite_png []byte

//go:embed Forgotten_Planet.mp3
var forgotten_planet_mp3 []byte

//go:embed melancholic_track.mp3
var melancholic_track_mp3 []byte

//go:embed click.mp3
var click_mp3 []byte

//go:embed upgrade.mp3
var upgrade_mp3 []byte

//go:embed error.mp3
var error_mp3 []byte

var SpriteSheet *ebiten.Image

var RockSpriteFull *ebiten.Image
var RockSpriteCracked1 *ebiten.Image
var RockSpriteCracked2 *ebiten.Image
var RockSpriteShattered *ebiten.Image
var MarketplaceSprite *ebiten.Image

var AudioContext *audio.Context
var HealthyMusicPlayer *audio.Player
var MelancholicMusicPlayer *audio.Player
var ClickSFXPlayer *audio.Player
var UpgradeSFXPlayer *audio.Player
var ErrorSFXPlayer *audio.Player

func init() {
	// Image assets
	img, _, err := image.Decode(bytes.NewReader(sprite_png))
	if err != nil {
		log.Fatal(err)
	}
	SpriteSheet = ebiten.NewImageFromImage(img)

	// Placeholder coordinates - MUST BE VERIFIED AGAINST ACTUAL sprite.png
	RockSpriteFull = SpriteSheet.SubImage(image.Rect(0, 0, 128, 128)).(*ebiten.Image)
	RockSpriteCracked1 = SpriteSheet.SubImage(image.Rect(128, 0, 256, 128)).(*ebiten.Image)
	RockSpriteCracked2 = SpriteSheet.SubImage(image.Rect(256, 0, 384, 128)).(*ebiten.Image)
	RockSpriteShattered = SpriteSheet.SubImage(image.Rect(384, 0, 512, 128)).(*ebiten.Image)

	marketplaceRect := image.Rect(20, 20, 120, 120) // Assuming this is correct
	MarketplaceSprite = SpriteSheet.SubImage(marketplaceRect).(*ebiten.Image)

	// Audio assets
	AudioContext = audio.NewContext(44100) // Standard sample rate

	// Load Healthy Music
	sHealthy, err := mp3.DecodeWithSampleRate(AudioContext.SampleRate(), bytes.NewReader(forgotten_planet_mp3))
	if err != nil {
		log.Fatal(err)
	}
	HealthyMusicPlayer, err = audio.NewPlayer(AudioContext, audio.NewInfiniteLoop(sHealthy, sHealthy.Length()))
	if err != nil {
		log.Fatal(err)
	}
	HealthyMusicPlayer.SetVolume(0.5) // Default volume

	// Load Melancholic Music
	sMelancholic, err := mp3.DecodeWithSampleRate(AudioContext.SampleRate(), bytes.NewReader(melancholic_track_mp3))
	if err != nil {
		log.Fatal(err)
	}
	MelancholicMusicPlayer, err = audio.NewPlayer(AudioContext, audio.NewInfiniteLoop(sMelancholic, sMelancholic.Length()))
	if err != nil {
		log.Fatal(err)
	}
	MelancholicMusicPlayer.SetVolume(0.0) // Start muted

	// Load Click SFX
	sClick, err := mp3.DecodeWithSampleRate(AudioContext.SampleRate(), bytes.NewReader(click_mp3))
	if err != nil {
		log.Fatal(err)
	}
	ClickSFXPlayer, err = audio.NewPlayer(AudioContext, sClick)
	if err != nil {
		log.Fatal(err)
	}

	// Load Upgrade SFX
	sUpgrade, err := mp3.DecodeWithSampleRate(AudioContext.SampleRate(), bytes.NewReader(upgrade_mp3))
	if err != nil {
		log.Fatal(err)
	}
	UpgradeSFXPlayer, err = audio.NewPlayer(AudioContext, sUpgrade)
	if err != nil {
		log.Fatal(err)
	}

	// Load Error SFX
	sError, err := mp3.DecodeWithSampleRate(AudioContext.SampleRate(), bytes.NewReader(error_mp3))
	if err != nil {
		log.Fatal(err)
	}
	ErrorSFXPlayer, err = audio.NewPlayer(AudioContext, sError)
	if err != nil {
		log.Fatal(err)
	}
}