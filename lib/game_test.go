package SheepAndPelican

import (
	"github.com/bgnori/SheepAndPelican/lib"
	"testing"
)

func TestNewGame(t *testing.T) {
	g := SheepAndPelican.NewGame()
	g.ShowTextArt()
}
