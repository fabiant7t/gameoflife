package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Cell struct {
	isAlive bool
}

func (c *Cell) Iterate(neighbourCount uint8) {
	if c.isAlive {
		if neighbourCount < 2 {
			c.isAlive = false
		} else if neighbourCount > 3 {
			c.isAlive = false
		}
	} else {
		if neighbourCount > 3 {
			c.isAlive = true
		}
	}
}

type Position struct {
	Row    int
	Column int
}

type board struct {
	rows    int
	columns int
	matrix  [][]Cell
}

func (b *board) Initialize() {
	for i := range b.matrix {
		for j := range len(b.matrix[i]) {
			cell := b.matrix[i][j]
			cell.isAlive = rand.Intn(100) < 20
			b.matrix[i][j] = cell
		}
	}
}

func (b *board) Neighbours() [][]uint8 {
	neighbours := make([][]uint8, b.rows)
	for i := range neighbours {
		neighbours[i] = make([]uint8, b.columns)
	}

	for row := range b.matrix {
		for col := range len(b.matrix[row]) {
			var neighbourCount uint8
			for _, pos := range []Position{
				{row - 1, col - 1},
				{row - 1, col},
				{row - 1, col + 1},
				{row, col - 1},
				{row, col + 1},
				{row + 1, col - 1},
				{row + 1, col},
				{row + 1, col + 1},
			} {
				cell, err := b.Cell(pos.Row, pos.Column)
				if err == nil {
					if cell.isAlive {
						neighbourCount++
					}
				}
			}
			neighbours[row][col] = neighbourCount
		}
	}
	return neighbours
}

func (b *board) Iterate() error {
	neighbours := b.Neighbours()

	for i := range b.matrix {
		for j := range len(b.matrix[i]) {
			cell, err := b.Cell(i, j)
			if err != nil {
				return err
			}
			cell.Iterate(neighbours[i][j])
		}
	}
	return nil
}

func (b *board) Cell(row, column int) (*Cell, error) {
	if row < 0 {
		return nil, errors.New("Row must not be negative")
	} else if row > b.rows-1 {
		return nil, errors.New("Row exceeds dimension")
	}
	if column < 0 {
		return nil, errors.New("Column must not be negative")
	} else if column > b.columns-1 {
		return nil, errors.New("Column exceeds dimension")
	}
	return &b.matrix[row][column], nil
}

func (b *board) String() string {
	var sb strings.Builder
	for i := range b.matrix {
		for j := range len(b.matrix[i]) {
			cell := b.matrix[i][j]
			if cell.isAlive {
				sb.WriteString(" * ")
			} else {
				sb.WriteString("   ")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func NewBoard(rows, columns int) *board {
	b := &board{
		rows:    rows,
		columns: columns,
		matrix:  make([][]Cell, rows),
	}
	for i := range b.matrix {
		b.matrix[i] = make([]Cell, columns)
	}
	b.Initialize()
	return b
}

func mainOld() {
	b := NewBoard(40, 40)
	for i := 0; i < 20; i++ {
		fmt.Println(b.String())
		time.Sleep(1 * time.Second)
		b.Iterate()
	}
}

// ----------- Bubble Tea integration -----------

type model struct {
	board *board
	tick  int
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		m.board.Iterate()
		m.tick++
		return m, tick()
	case tea.KeyMsg:
		return m, tea.Quit // exit on any key
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Conway's Game of Life — generation %d\n\n%s\n(press any key to quit)", m.tick, m.board.String())
}

func main() {
	b := NewBoard(40, 40)

	p := tea.NewProgram(model{board: b})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
