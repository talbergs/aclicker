package shaders

import "github.com/hajimehoshi/ebiten/v2"

// ShaderEffect represents a shader and its uniforms.
type ShaderEffect struct {
	Shader   *ebiten.Shader
	Uniforms map[string]interface{}
}

// Apply applies a chain of shader effects to a source image.
func Apply(source *ebiten.Image, clickGridTexture *ebiten.Image, effects ...ShaderEffect) *ebiten.Image {
	if len(effects) == 0 {
		return source
	}

	// Create two offscreen images for ping-ponging
	offscreen1 := ebiten.NewImage(source.Size())
	offscreen2 := ebiten.NewImage(source.Size())

	// Initial draw
	op := &ebiten.DrawImageOptions{}
	offscreen1.DrawImage(source, op)

	for i, effect := range effects {
		src := offscreen1
		dst := offscreen2
		if i%2 != 0 {
			src = offscreen2
			dst = offscreen1
		}

		dst.Clear()
		dst.DrawRectShader(source.Bounds().Dx(), source.Bounds().Dy(), effect.Shader, &ebiten.DrawRectShaderOptions{
			Images:   [4]*ebiten.Image{src, clickGridTexture}, // Pass clickGridTexture as the second image
			Uniforms: effect.Uniforms,
		})
	}

	if len(effects)%2 == 0 {
		return offscreen1
	}

	return offscreen2
}
