package main

import (
	"main/models"
	"main/tracking"
	"main/utils"
	"os"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	debugui debugui.DebugUI
}

const (
	EditMode     int = 0
	ViewMode         = 1
	TrackingMode     = 2
)

var Mode = EditMode

var Model models.Model

var FaceData tracking.TrackingData

func (g *Game) Update() error {
	utils.Tick += 1
	mx, my := ebiten.CursorPosition()
	utils.MousePos.X = float64(mx-66*5) / 3
	utils.MousePos.Y = float64(my) / 3

	if inpututil.IsKeyJustPressed(ebiten.KeyV) && !ebiten.IsKeyPressed(ebiten.KeyShift) && !ebiten.IsKeyPressed(ebiten.KeyControl) {
		if Mode == EditMode {
			Mode = ViewMode
		} else if Mode == ViewMode {
			Mode = TrackingMode
		} else if Mode == TrackingMode {
			Mode = EditMode
		}
	}

	FaceData.Update()

	ebiten.SetWindowMousePassthrough(false)

	switch Mode {
	case EditMode:
		EditUpdate(&g.debugui, &Model)
	case ViewMode:
		ebiten.SetWindowMousePassthrough(true)
		ViewUpdate()
	case TrackingMode:
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch Mode {
	case EditMode:
		EditDraw(screen, *g, &g.debugui)
	case ViewMode:
		ViewDraw(screen, FaceData, Model)
	case TrackingMode:
		TrackDraw(screen, FaceData)
	}
}

func (g *Game) Layout(ow, oh int) (sw, sh int) {
	return 426 * 3, 240 * 3
}

func main() {
	ebiten.SetWindowSize(426*5, 240*5)

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowDecorated(false)

	Model = models.NewModel()

	_, err := os.ReadDir("./model_files")
	if err != nil {
		os.Mkdir("model_files", os.ModePerm)
	}

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
