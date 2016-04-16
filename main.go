package main

import (
	"fmt"
	. "github.com/bgnori/SheepAndPelican/lib"
)

func main() {
	g := NewGame()
	fmt.Printf("%+v\n", g)
	g.ShowTextArt()

	{
		m := Move{
			Src: Coor{Row: 1, Col: 1},
			Dst: Coor{Row: 1, Col: 4},
		}
		b, xs := g.IsLeagal(m)
		fmt.Printf("%+v is %b, %v\n", m, b, xs)
	}

	{
		m := Move{
			Src: Coor{Row: 1, Col: 1},
			Dst: Coor{Row: 3, Col: 4},
		}
		b, xs := g.IsLeagal(m)
		fmt.Printf("%+v is %b, %v\n", m, b, xs)
		g.MakeMove(m)
	}
	fmt.Printf("%+v\n", g)
	g.ShowTextArt()
}
