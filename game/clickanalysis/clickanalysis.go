package clickanalysis

import (
	"image"
	"image/color"
	"time"
)

const (
	defaultGridWidth  = 100
	defaultGridHeight = 100
	defaultDecayRate  = 0.05 // How much heat decays per second
	defaultMaxHeat    = 1.0  // Max heat value for a cell
	clickHeatAmount   = 0.5  // How much heat a single click adds
)

// ClickGrid represents a 2D grid that tracks click "heat" over time.
type ClickGrid struct {
	Width     int
	Height    int
	Grid      [][]float64
	DecayRate float64
	MaxHeat   float64
	lastUpdate time.Time
}

// NewClickGrid creates and initializes a new ClickGrid.
func NewClickGrid(width, height int) *ClickGrid {
	grid := make([][]float64, height)
	for i := range grid {
		grid[i] = make([]float64, width)
	}
	return &ClickGrid{
		Width:     width,
		Height:    height,
		Grid:      grid,
		DecayRate: defaultDecayRate,
		MaxHeat:   defaultMaxHeat,
		lastUpdate: time.Now(),
	}
}

// AddClick adds "heat" to the grid at the specified (x, y) coordinates.
// Coordinates are expected to be in screen space and will be mapped to grid space.
func (cg *ClickGrid) AddClick(screenX, screenY, screenWidth, screenHeight int) {
	gridX := int(float64(screenX) / float64(screenWidth) * float64(cg.Width))
	gridY := int(float64(screenY) / float64(screenHeight) * float64(cg.Height))

	if gridX >= 0 && gridX < cg.Width && gridY >= 0 && gridY < cg.Height {
		cg.Grid[gridY][gridX] += clickHeatAmount
		if cg.Grid[gridY][gridX] > cg.MaxHeat {
			cg.Grid[gridY][gridX] = cg.MaxHeat
		}
	}
}

// Update decays the heat in the grid over time.
func (cg *ClickGrid) Update() {
	now := time.Now()
	deltaTime := now.Sub(cg.lastUpdate).Seconds()
	cg.lastUpdate = now

	if deltaTime <= 0 {
		return
	}

	for y := 0; y < cg.Height; y++ {
		for x := 0; x < cg.Width; x++ {
			cg.Grid[y][x] -= cg.DecayRate * deltaTime
			if cg.Grid[y][x] < 0 {
				cg.Grid[y][x] = 0
			}
		}
	}
}

// ToRGBA converts the current heat grid into an image.RGBA for shader consumption.
func (cg *ClickGrid) ToRGBA() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, cg.Width, cg.Height))
	for y := 0; y < cg.Height; y++ {
		for x := 0; x < cg.Width; x++ {
			heat := cg.Grid[y][x]
			// Normalize heat to 0-255 for RGB.
			// For now, let's just use red channel to represent heat.
			// We can make this more sophisticated later (e.g., heatmap colors).
			val := uint8(heat / cg.MaxHeat * 255)
			img.SetRGBA(x, y, color.RGBA{R: val, A: 255})
		}
	}
	return img
}
