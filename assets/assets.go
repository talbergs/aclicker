package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed sprite.png
var sprite_png []byte

var SpriteSheet *ebiten.Image

var RockSprite *ebiten.Image
var MarketplaceSprite *ebiten.Image

func init() {
	img, _, err := image.Decode(bytes.NewReader(sprite_png))
	if err != nil {
		log.Fatal(err)
	}
	SpriteSheet = ebiten.NewImageFromImage(img)

	// Estimated coordinates from the image
	rockRect := image.Rect(700, 1000, 240, 600)
	!dbg(image.ZR.Size())
	RockSprite = SpriteSheet.SubImage(rockRect).(*ebiten.Image)

	marketplaceRect := image.Rect(20, 20, 120, 120)
	MarketplaceSprite = SpriteSheet.SubImage(marketplaceRect).(*ebiten.Image)
}
