package visualize

import "github.com/hajimehoshi/ebiten/v2"

type MockGame struct{}

func (g *MockGame) Update() error {
	return nil
}

func (g *MockGame) Draw(_ *ebiten.Image) {}

func (g *MockGame) Layout(_, _ int) (_, _ int) {
	return 0, 0
}
