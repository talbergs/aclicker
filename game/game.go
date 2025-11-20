package game

// Rock represents the entity that is clicked.
type Rock struct {
	Health int
}

// Player represents the user's state.
type Player struct {
	Dust int
}

// Game holds the overall game state.
type Game struct {
	TheRock   *Rock
	ThePlayer *Player
}

// NewGame creates a new game state with initial values.
func NewGame() *Game {
	return &Game{
		TheRock: &Rock{
			Health: 10000000,
		},
		ThePlayer: &Player{
			Dust: 0,
		},
	}
}

// Click handles the logic for a single click on the rock.
func (g *Game) Click() {
	g.TheRock.Health--
	g.ThePlayer.Dust++
}
