package game

import (
	"math/rand"
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
	snakeField  fieldType = "snake"
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
	dir      direction
	snake    []coordinate
	apple    *coordinate
}

func New(n int) *Game {
	g := Game{
		N:      n,
		Fields: make([][]fieldType, n),
		Score:  0,
		dir:    Right,
		snake:  []coordinate{{0, 0}},
	}

	for i := 0; i < len(g.Fields); i++ {
		g.Fields[i] = make([]fieldType, g.N)
	}

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
	if (g.dir^newDir)%2 == 0 {
		return false
	}
	g.dir = newDir
	return true
	// TODO: fix bug, doenst set new dir but still renders
}

func (g *Game) updateFields() {
	for y := 0; y < g.N; y++ {
		for x := 0; x < g.N; x++ {
			if slices.Contains(g.snake, coordinate{x, y}) {
				g.Fields[y][x] = snakeField
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
	tail := g.snake[0 : len(g.snake)-1]
	for i := 0; i < len(tail); i++ {
		g.snake[i] = g.snake[i+1]
	}
	// move head
	head := &g.snake[len(g.snake)-1]
	switch g.dir {
	case Up:
		head.y--
	case Left:
		head.x--
	case Down:
		head.y++
	case Right:
		head.x++
	}

	if head.x < 0 || head.x >= g.N || head.y < 0 || head.y >= g.N || slices.Contains(tail, coordinate{head.x, head.y}) {
		g.GameOver = true
	}
}

func (g *Game) renderApple() {
	if g.apple == nil {
		used := make(map[coordinate]struct{})
		for _, c := range g.snake {
			used[c] = struct{}{}
		}

		var rX, rY int
		for {
			rX, rY = rand.Intn(g.N), rand.Intn(g.N)
			if _, ok := used[coordinate{rX, rY}]; !ok {
				break
			}
		}

		g.apple = &coordinate{rX, rY}
	}

	head := g.snake[len(g.snake)-1]
	if g.apple.x == head.x && g.apple.y == head.y {
		g.Score++
		g.snake = append([]coordinate{g.snake[0]}, g.snake...)
		g.apple = nil
	}
}
