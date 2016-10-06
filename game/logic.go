package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const EasyMode = 0
const MediumMode = 1
const HardMode = 2
const EasyLvlDimension = 9
const EasyLvlBombsNumber = 10
const MediumLvlDimension = 16
const MediumLvlBombsNumber = 40
const HardLvlDimension = 30
const HardLvlBombsNumber = 116

type point struct {
	touched     bool
	isBomb      bool
	bombsNumber int
	hasFlag     bool
}

// Board represents minesweeper board.
type Board struct {
	bombsNumber int
	flagsLeft   int
	dimension   int
	field       [][]*point
	gameOver    bool
	gameWin     bool
}

func (p *point) toString() string {
	return " " + strconv.FormatBool(p.isBomb) + " neighbours " + strconv.Itoa(p.bombsNumber)
}

func (b *Board) setBoard() {
	for i := 0; i < b.dimension; i++ {
		row := []*point{}
		for j := 0; j < b.dimension; j++ {
			row = append(row, new(point))
		}
		b.field = append(b.field, row)
	}
	b.setBombs()
	b.setBombsNeighbours()
}

func (b *Board) setBombs() {
	rand.Seed(time.Now().UTC().UnixNano())
	count := b.bombsNumber
	for count > 0 {
		x := rand.Intn(b.dimension)
		y := rand.Intn(b.dimension)
		if b.field[x][y].isBomb == false {
			b.field[x][y].isBomb = true
			count--
		}
	}
}

func (b *Board) setBombsNeighbours() {
	for i := 0; i < b.dimension; i++ {
		for j := 0; j < b.dimension; j++ {
			coords := []int{-1, 0, 1}
			for _, ki := range coords {
				for _, kj := range coords {
					if ki == 0 && kj == 0 {
						continue
					} else if ((ki+i >= 0) && (ki+i < b.dimension)) &&
						((kj+j >= 0) && (kj+j < b.dimension)) {
						if b.field[ki+i][kj+j].isBomb == true {
							b.field[i][j].bombsNumber++
						}
					}
				}
			}
		}
	}
}

func (b *Board) performRightClick(col int, row int) {
	newBoardState := b.field
	if newBoardState[row][col].touched == true {
	} else {
		newBoardState[row][col].hasFlag = true
	}
	b.updateState(newBoardState)
}

func (b *Board) performLeftClick(rowCoord int, colCoord int) {
	newBoardState := b.field
	if newBoardState[rowCoord][colCoord].touched == true {
	} else {
		if newBoardState[rowCoord][colCoord].isBomb == true {
			//TODO:make visible all bomb points
			newBoardState[rowCoord][colCoord].touched = true
			b.updateState(newBoardState)
			b.gameOver = true
			b.showAllBombs()
		} else {
			bombs := newBoardState[rowCoord][colCoord].bombsNumber
			if bombs > 0 {
				newBoardState[rowCoord][colCoord].touched = true
				b.updateState(newBoardState)
				if b.isWin() == true {
					b.gameWin = true
				}
			} else {
				//empty amount of bomb neighbours
				newBoardState[rowCoord][colCoord].touched = true
				coords := []int{-1, 0, 1}
				for _, ki := range coords {
					for _, kj := range coords {
						if ki == 0 && kj == 0 {
							continue
						} else if ((ki+rowCoord >= 0) && (ki+rowCoord < b.dimension)) &&
							((kj+colCoord >= 0) && (kj+colCoord < b.dimension)) {
							b.performLeftClick(ki+rowCoord, kj+colCoord)
						}
					}
				}
			}
		}

	}
}

func (b *Board) choose(col int, row int) {
	newBoardState := b.field
	if newBoardState[row][col].touched == true {
	} else {
		if newBoardState[row][col].isBomb == true {
			newBoardState[row][col].touched = true
			b.updateState(newBoardState)
			b.gameOver = true
			b.showAllBombs()
		} else {
			bombs := newBoardState[row][col].bombsNumber
			if bombs > 0 {
				newBoardState[row][col].touched = true
				b.updateState(newBoardState)
				if b.isWin() == true {
					b.gameWin = true
				}
			} else {
				//empty amount of bomb neighbours
				newBoardState[row][col].touched = true
				coords := []int{-1, 0, 1}
				for _, ki := range coords {
					for _, kj := range coords {
						if ki == 0 && kj == 0 {
							continue
						} else if ((ki+row >= 0) && (ki+row < b.dimension)) &&
							((kj+col >= 0) && (kj+col < b.dimension)) {
							b.choose(ki+row, kj+col)
						}
					}
				}
			}
		}
	}
}

func (b *Board) updateState(newBoard [][]*point) {
	b.field = newBoard
}

func (b *Board) isWin() bool {
	for i := 0; i < b.dimension; i++ {
		for j := 0; j < b.dimension; j++ {
			if (b.field[i][j].isBomb == false) &&
				(b.field[i][j].touched == false) {
				return false
			}
		}
	}
	return true
}

// continuePlaying tells whether we should keep playing.
func (b *Board) continuePlaying() bool {
	return !b.gameOver
}

func (b *Board) showBoard() {
	for i := 0; i < b.dimension; i++ {
		for j := 0; j < b.dimension; j++ {
			if b.field[i][j].touched {
				if b.field[i][j].isBomb {
					fmt.Print("x" + " ")

				} else {
					fmt.Print(strconv.Itoa(b.field[i][j].bombsNumber) + " ")

				}
			} else {
				fmt.Print("*" + " ")
			}
		}
		fmt.Println()
	}
}

func (b *Board) showAllBombs() {
	for i := 0; i < b.dimension; i++ {
		for j := 0; j < b.dimension; j++ {
			if b.field[i][j].isBomb {
				b.field[i][j].touched = true
			}
		}
	}
}

func (b *Board) initGame(mode int) {
	dimension := -1
	bombsNumber := -1
	switch mode {
	case EasyMode:
		dimension = EasyLvlDimension
		bombsNumber = EasyLvlBombsNumber
	case MediumMode:
		dimension = MediumLvlDimension
		bombsNumber = MediumLvlBombsNumber
	case HardMode:
		dimension = HardLvlDimension
		bombsNumber = HardLvlBombsNumber
	}
	b.dimension = dimension
	b.bombsNumber = bombsNumber
	b.flagsLeft = bombsNumber
	b.setBoard()
}

func (b *Board) resetGame() {
	b.field = [][]*point{}
	b.gameWin = false
	b.gameOver = false
}

//func main() {
//	b := Board {dimension:EasyLvlDimension,
//		   bombsNumber: EasyLvlBombsNumber, }
//	b.setBoard(b.dimension)
//	fmt.Println(strconv.FormatBool(b.gameOver))
//	b.showBoard()
//	fmt.Println()
//	fmt.Println()
//	count := 0
//	for b.continuePlaying() {
//		b.performLeftClick(rand.Intn(b.dimension),rand.Intn(b.dimension))
//		b.showBoard()
//		fmt.Println()
//		fmt.Println()
//		count++
//	}
//	fmt.Println(count)
//}
