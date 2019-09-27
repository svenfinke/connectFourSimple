package main

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/jroimartin/gocui"
	"log"
)

type GameData struct {
	currentPlayer int
	gameData [][]int
}

const (
	VIEW_MENU = "menu"
	VIEW_GAME_GRID = "gamegrid"
	VIEW_MESSAGES = "messages"
)

var (
	gui       *gocui.Gui
	game      GameData
	menuIndex int
)

func main(){
	gui.SetManagerFunc(layoutFunc)

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gui.SetKeybinding(VIEW_MENU, gocui.KeyArrowDown, gocui.ModNone, nextItem); err != nil {
		log.Panicln(err)
	}

	if err := gui.SetKeybinding(VIEW_MENU, gocui.KeyArrowUp, gocui.ModNone, prevItem); err != nil {
		log.Panicln(err)
	}

	if err := gui.SetKeybinding(VIEW_MENU, gocui.KeyEnter, gocui.ModNone, dropToken); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func init(){
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panic(err)
	}
	defer g.Close()

	game.currentPlayer = 0
	for x := 0 ; x < 7 ; x++ {
		game.gameData = append(game.gameData, []int{})
		for y := 0 ; y < 6 ; y++ {
			game.gameData[x] = append(game.gameData[x], -1)
		}
	}
}

func layoutFunc(g *gocui.Gui) error{
	e, done := renderHeader()
	e, done = renderMenu()
	e, done = renderMessages()
	e, done = renderGameGrid()

	g.SetCurrentView(VIEW_MENU)
	if done {
		return e
	}

	return nil
}

func renderMessages() (error, bool) {
	maxX, _ := gui.Size()
	if v, err := gui.SetView(VIEW_MESSAGES, 7, 3, maxX-1, 12); err != nil {
		if err != gocui.ErrUnknownView {
			return err, true
		}

		v.Wrap = true
		v.Autoscroll = true
	}
	return nil, false
}

func renderMenu() (error, bool) {
	if _, err := gui.SetView(VIEW_MENU, 0, 3, 6, 12); err != nil {
		if err != gocui.ErrUnknownView {
			return err, true
		}

		printMenu()
	}
	return nil, false
}

func renderHeader() (error, bool) {
	maxX, _ := gui.Size()
	if v, err := gui.SetView("header", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err, true
		}
		fmt.Fprintln(v, " connectFour")
	}
	return nil, false
}

func renderGameGrid() (error, bool) {
	maxX, maxY := gui.Size()
	if _, err := gui.SetView(VIEW_GAME_GRID, maxX/2-8, maxY/2-4, maxX/2+8, maxY/2+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err, true
		}

		printGame()
	}
	return nil, false
}

func printGame() {
	v, err := gui.View(VIEW_GAME_GRID)
	if err != nil {
		log.Panic(err)
	}
	v.Clear()
	for y := 5; y >= 0; y-- {
		row := ""
		for x := 0; x < 7; x++ {
			char := " "
			switch game.gameData[x][y] {
			case 0:
				char = "X"
				break
			case 1:
				char = "O"
				break
			}
			row = row + " " + char
		}
		fmt.Fprintln(v, row)
	}
}

func printMenu() {
	v, err := gui.View(VIEW_MENU)
	if err != nil {
		log.Panic(err)
	}
	v.Clear()
	for x := 0 ; x < 7 ; x++ {
		if menuIndex == x {
			fmt.Fprintln(v, " ", color.Red.Sprint(x))
		} else {
			fmt.Fprintln(v, " ", color.Yellow.Sprint(x))
		}
	}

}

func nextItem(g *gocui.Gui, v *gocui.View) error {
	menuIndex++
	if v != nil {
		printMenu()
	}
	return nil
}


func prevItem(g *gocui.Gui, v *gocui.View) error {
	menuIndex--
	if v != nil {
		printMenu()
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func dropToken(g *gocui.Gui, v *gocui.View) error {
	row := game.gameData[menuIndex]

	if !moveValid() {
		return nil
	}

	for idx, item := range row {
		if item == -1 {
			game.gameData[menuIndex][idx] = game.currentPlayer
			game.currentPlayer++
			if game.currentPlayer > 1 {
				game.currentPlayer = 0
			}
			break
		}
	}

	printGame()

	return nil
}

func moveValid() bool {
	if game.gameData[menuIndex][5] != -1 {
		printMessage(fmt.Sprintf("Row %v is already full", menuIndex))
		return false
	}

	return true
}

func gameFinished() bool {
	return false
}

func printMessage(msg string) {
	v, err := gui.View(VIEW_MESSAGES)
	if err != nil {
		log.Panic(err)
	}

	fmt.Fprintln(v, msg)
}