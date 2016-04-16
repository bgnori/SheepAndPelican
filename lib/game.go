package SheepAndPelican

import (
	"errors"
	"fmt"
)

type SquareState int

const (
	Empty SquareState = iota
	Pelican
	Sheep
	Wall
)

func (ss SquareState) String() string {
	switch ss {
	case Empty:
		return "."
	case Pelican:
		return "P"
	case Sheep:
		return "S"
	case Wall:
		return "#"
	}
	panic("unkown SquareState")
}

type GameState int

const (
	PelicantoPlay GameState = iota
	SheeptoPlay
	PelicanHasWon
	SheepHasWon
	PelicanStaleMate
	SheepStaleMate
)

func (gs GameState) String() string {
	switch gs {
	case PelicantoPlay:
		return "Pelican to play." //FistTo Play
	case SheeptoPlay:
		return "Sheep to play."
	case PelicanHasWon:
		return "Pelican has won."
	case SheepHasWon:
		return "Sheep has won"
	case PelicanStaleMate:
		return "Pelican has no leagal moves, stale mate."
	case SheepStaleMate:
		return "Sheep has no leagal moves, stale mate."
	}
	panic("unkown GameState")
}

type Row int

const RowMin Row = 0
const RowMax Row = 8

type Col int

const ColMin Col = 0
const ColMax Col = 8

type Coor struct {
	Row Row
	Col Col
}

func (c Coor) Clone() Coor {
	return Coor{Row: c.Row, Col: c.Col}
}

type Move struct {
	Src Coor
	Dst Coor
}

type Game struct {
	GameState GameState
	Board     [RowMax][ColMax]SquareState
}

func NewGame() *Game {
	g := &Game{
		GameState: PelicantoPlay,
	}
	for i := RowMin; i < RowMax; i++ {
		g.Board[i][int(ColMin)] = Wall
		g.Board[i][ColMax-1] = Wall
	}
	for i := ColMin; i < ColMax; i++ {
		g.Board[RowMin][i] = Wall
		g.Board[RowMax-1][i] = Wall
	}

	g.Set(Coor{1, 1}, Pelican)
	g.Set(Coor{1, 2}, Pelican)
	g.Set(Coor{1, 3}, Sheep)
	g.Set(Coor{1, 4}, Sheep)
	g.Set(Coor{1, 5}, Pelican)
	g.Set(Coor{1, 6}, Pelican)

	g.Set(Coor{6, 1}, Sheep)
	g.Set(Coor{6, 2}, Sheep)
	g.Set(Coor{6, 3}, Pelican)
	g.Set(Coor{6, 4}, Pelican)
	g.Set(Coor{6, 5}, Sheep)
	g.Set(Coor{6, 6}, Sheep)

	return g
}

func (g *Game) ShowTextArt() {
	fmt.Printf("GameState: %v\n", g.GameState)
	for i := 0; i < 8; i++ {
		fmt.Printf("%+v\n", g.Board[i])
	}

}

func (g *Game) BoundaryCheck(c Coor) {
	msg := ""
	if !(RowMin <= c.Row && c.Row < RowMax) {
		msg += "Row is out of bound!"
	}
	if !(ColMin <= c.Col && c.Col < ColMax) {
		msg += "COl is out of bound!"
	}
	if msg != "" {
		panic(msg)
	}
}

func (g *Game) Set(c Coor, v SquareState) {
	g.BoundaryCheck(c)
	g.Board[c.Row][c.Col] = v
}

func (g *Game) Get(c Coor) SquareState {
	g.BoundaryCheck(c)
	return g.Board[c.Row][c.Col]
}

func (here Coor) RowStep(goal Coor) Coor {
	if goal.Row == here.Row {
		return here.Clone()
	}
	if goal.Row-here.Row > 0 {
		return Coor{Row: here.Row + 1, Col: here.Col}
	} else {
		return Coor{Row: here.Row - 1, Col: here.Col}
	}
	panic("never reach here")
}

func (here Coor) ColStep(goal Coor) Coor {
	if goal.Col == here.Col {
		return here.Clone()
	}
	if goal.Col-here.Col > 0 {
		return Coor{Col: here.Col + 1, Row: here.Row}
	} else {
		return Coor{Col: here.Col - 1, Row: here.Row}
	}
	panic("never reach here")
}

func (g *Game) isEmpty(c Coor) bool {
	return g.Get(c) == Empty
}

func (g *Game) isConsistentWith(c Coor) bool {
	sq := g.Get(c)
	return (sq == Pelican && g.GameState == PelicantoPlay) ||
		(sq == Sheep && g.GameState == SheeptoPlay)
}

func (g *Game) hasRowFirstPath(m Move) (bool, []Coor) {
	sq := m.Src.RowStep(m.Dst)
	xs := make([]Coor, 0)
	for ; sq.Row != m.Dst.Row; sq = sq.RowStep(m.Dst) {
		if !g.isEmpty(sq) {
			return false, nil
		}
		xs = append(xs, sq)
	}
	if !g.isEmpty(sq) {
		return false, nil
		xs = append(xs, sq)
	}
	for ; sq.Col != m.Dst.Col; sq = sq.ColStep(m.Dst) {
		if !g.isEmpty(sq) {
			return false, nil
		}
		xs = append(xs, sq)
	}
	return true, xs
}

func (g *Game) hasColFirstPath(m Move) (bool, []Coor) {
	sq := m.Src.ColStep(m.Dst)
	xs := make([]Coor, 0)
	for ; sq.Col != m.Dst.Col; sq = sq.ColStep(m.Dst) {
		if !g.isEmpty(sq) {
			return false, nil
		}
		xs = append(xs, sq)
	}
	if !g.isEmpty(sq) {
		return false, nil
	}
	for ; sq.Row != m.Dst.Row; sq = sq.RowStep(m.Dst) {
		if !g.isEmpty(sq) {
			return false, nil
		}
		xs = append(xs, sq)
	}
	return true, xs
}

func (g *Game) IsLeagal(m Move) (bool, []Coor) {
	r := false
	xs := make([]Coor, 0)
	if !g.isEmpty(m.Dst) {
		return false, nil
	}
	if !g.isConsistentWith(m.Src) {
		return false, nil
	}
	if b, ys := g.hasRowFirstPath(m); b {
		r = r || b
		xs = append(xs, ys...)
	}
	if b, ys := g.hasColFirstPath(m); b {
		r = r || b
		xs = append(xs, ys...)
	}
	return r, xs
}

func (g *Game) GameStateToSqaureState() SquareState {
	switch g.GameState {
	case PelicantoPlay:
		return Pelican
	case SheeptoPlay:
		return Sheep
	}
	panic("Bad GameState")
}

func (g *Game) HorizotalCheck(player SquareState) (bool, []Coor) {
	var xs []Coor

	for r := 1; r < 6; r++ {
	Checking:
		for start := 1; start < 4; start += 1 {
			xs = make([]Coor, 0, 6)
			for checked := start; checked < start+4; checked += 1 {
				c := Coor{Row: Row(r), Col: Col(checked)}
				if g.Get(c) == player {
					xs = append(xs, c)
				} else {
					continue Checking
				}
			}
			return true, xs
		}
	}
	return false, xs
}

func (g *Game) CalcGameState() GameState {
	// called post-move process, to figure out outcome of the game.
	if !(g.GameState == PelicantoPlay || g.GameState == SheeptoPlay) {
		panic("Bad GameState")
	}
	//x, xs := g.HorizotalCheck()

	return PelicantoPlay
}

func (g *Game) NextTurn() error {
	switch g.GameState {
	case PelicantoPlay:
		g.GameState = SheeptoPlay
		return nil
	case SheeptoPlay:
		g.GameState = PelicantoPlay
		return nil
	case PelicanHasWon:
		fallthrough
	case SheepHasWon:
		fallthrough
	case PelicanStaleMate:
		fallthrough
	case SheepStaleMate:
		return errors.New(fmt.Sprintf("Game has ended, %v", g.GameState))
	}
	return errors.New("Unkown state")
}

func (g *Game) MakeMove(m Move) {
	//done by swapping content, since destination must be empty.
	//Leagality is not scope of this function.
	src := g.Get(m.Src)
	dst := g.Get(m.Dst)

	g.Set(m.Dst, src)
	g.Set(m.Src, dst)

	g.GameState = g.CalcGameState()
	g.NextTurn()
}
