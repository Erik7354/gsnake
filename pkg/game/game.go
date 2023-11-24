package game

import (
	"math/rand"
	"reflect"
	"slices"
)

type direction = int
type fieldType = string

const (
	Up direction = iota
	Left
	Down
	Right
)

const (
	normalField fieldType = "normal"
	headField   fieldType = "head"
	tailField   fieldType = "tail"
	appleField  fieldType = "apple"
)

type coordinate struct {
	x int
	y int
}

type Game struct {
	N        int
	Fields   [][]fieldType
	Score    int
	GameOver bool
	Dir      direction
	head     coordinate
	tail     []coordinate
	apple    *coordinate
}

// New creates a new Snake Game.
// n should be >4
func New(n int) *Game {
	g := Game{
		N:      n,
		Fields: make([][]fieldType, n),
		Score:  0,
		Dir:    Right,
		head:   coordinate{1, 0},
		tail:   []coordinate{{0, 0}},
		apple:  &coordinate{n - 2, 0},
	}

	for i := 0; i < len(g.Fields); i++ {
		g.Fields[i] = make([]fieldType, g.N)
		for j := range g.Fields[i] {
			g.Fields[i][j] = normalField
		}
	}
	g.Fields[0][1] = headField
	g.Fields[0][0] = tailField
	g.Fields[0][n-2] = appleField

	return &g
}

func (g *Game) Render() {
	if g.GameOver {
		return
	}

	g.renderSnake()
	g.renderApple()
	g.updateFields()
}

func (g *Game) SetDirection(newDir direction) bool {
	if (g.Dir^newDir)%2 == 0 {
		return false
	}
	g.Dir = newDir
	return true
}

func (g *Game) updateFields() {
	for y := 0; y < g.N; y++ {
		for x := 0; x < g.N; x++ {
			if reflect.DeepEqual(g.head, coordinate{x, y}) {
				g.Fields[y][x] = headField
				continue
			}
			if slices.Contains(g.tail, coordinate{x, y}) {
				g.Fields[y][x] = tailField
				continue
			}
			if g.apple != nil && g.apple.x == x && g.apple.y == y {
				g.Fields[y][x] = appleField
				continue
			}
			g.Fields[y][x] = normalField
		}
	}
}

func (g *Game) renderSnake() {
	// move body
	for i := 0; i < len(g.tail)-1; i++ {
		g.tail[i] = g.tail[i+1]
	}
	g.tail[len(g.tail)-1] = g.head

	// move head
	switch g.Dir {
	case Up:
		g.head.y--
	case Left:
		g.head.x--
	case Down:
		g.head.y++
	case Right:
		g.head.x++
	}

	if g.head.x < 0 || g.head.x >= g.N || g.head.y < 0 || g.head.y >= g.N || slices.Contains(g.tail, coordinate{g.head.x, g.head.y}) {
		g.GameOver = true
	}
}

func (g *Game) renderApple() {
	if g.apple == nil {
		// fields used for head and tail
		used := make(map[coordinate]struct{})
		used[g.head] = struct{}{}
		for _, c := range g.tail {
			used[c] = struct{}{}
		}

		// generate random position for new apple
		var rX, rY int
		for {
			rX, rY = rand.Intn(g.N), rand.Intn(g.N)
			if _, ok := used[coordinate{rX, rY}]; !ok {
				break
			}
		}

		g.apple = &coordinate{rX, rY}
	}

	if g.apple.x == g.head.x && g.apple.y == g.head.y {
		g.Score++
		g.tail = append([]coordinate{g.tail[0]}, g.tail...)
		g.apple = nil
	}
}
