package shaders

import (
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed grayscale.kage
var grayscaleSrc []byte

//go:embed invert.kage
var invertSrc []byte

//go:embed warp.kage
var warpSrc []byte

//go:embed clickheat.kage
var clickHeatSrc []byte

//go:embed timeclick.kage
var timeClickSrc []byte

var (
	GrayscaleShader *ebiten.Shader
	InvertShader    *ebiten.Shader
	WarpShader      *ebiten.Shader
	ClickHeatShader *ebiten.Shader
	TimeClickShader *ebiten.Shader
)

func init() {
	var err error
	GrayscaleShader, err = ebiten.NewShader(grayscaleSrc)
	if err != nil {
		log.Fatal(err)
	}

	InvertShader, err = ebiten.NewShader(invertSrc)
	if err != nil {
		log.Fatal(err)
	}

	WarpShader, err = ebiten.NewShader(warpSrc)
	if err != nil {
		log.Fatal(err)
	}

	ClickHeatShader, err = ebiten.NewShader(clickHeatSrc)
	if err != nil {
		log.Fatal(err)
	}

	TimeClickShader, err = ebiten.NewShader(timeClickSrc)
	if err != nil {
		log.Fatal(err)
	}
}

// Grayscale returns a ShaderEffect for the grayscale shader.
func Grayscale() ShaderEffect {
	return ShaderEffect{Shader: GrayscaleShader}
}

// Invert returns a ShaderEffect for the invert shader.
func Invert() ShaderEffect {
	return ShaderEffect{Shader: InvertShader}
}

// Warp returns a ShaderEffect for the warp shader.
func Warp(time float32) ShaderEffect {
	return ShaderEffect{
		Shader: WarpShader,
		Uniforms: map[string]interface{}{
			"Time": time,
		},
	}
}

// ClickHeat returns a ShaderEffect for the click heat shader.
func ClickHeat() ShaderEffect {
	return ShaderEffect{Shader: ClickHeatShader}
}

// TimeClick returns a ShaderEffect for the time-click shader.
func TimeClick(time float32) ShaderEffect {
	return ShaderEffect{
		Shader: TimeClickShader,
		Uniforms: map[string]interface{}{
			"Time":   time,
		},
	}
}
