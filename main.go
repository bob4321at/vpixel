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
	EditMode      int = 0
	ViewMode          = 1
	NetworkConfig     = 2
	TrackingMode      = 3
)

var Mode = EditMode

var Model models.Model

var FaceData tracking.TrackingData

func (g *Game) Update() error {
	utils.Tick += 1
	mx, my := ebiten.CursorPosition()
	utils.MousePos.X = float64(mx-66*5) / 3
	utils.MousePos.Y = float64(my) / 3

	if inpututil.IsKeyJustPressed(ebiten.KeyV) && !ebiten.IsKeyPressed(ebiten.KeyShift) && !ebiten.IsKeyPressed(ebiten.KeyControl) && !ebiten.IsKeyPressed(ebiten.KeyAlt) {
		switch Mode {
		case EditMode:
			Mode = ViewMode
		case ViewMode:
			Mode = NetworkConfig
		case NetworkConfig:
			Mode = TrackingMode
		case TrackingMode:
			Mode = EditMode
		}
	}

	FaceData.Update()

	switch Mode {
	case EditMode:
		// ebiten.SetWindowMousePassthrough(false)
		EditUpdate(&g.debugui, &Model)
	case ViewMode:
		// ebiten.SetWindowMousePassthrough(true)
		ViewUpdate()
	case NetworkConfig:
		// ebiten.SetWindowMousePassthrough(false)
		NetworkModeUpdate(&g.debugui, &Model)
	case TrackingMode:
		// ebiten.SetWindowMousePassthrough(false)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch Mode {
	case EditMode:
		EditDraw(screen, *g, &g.debugui)
	case ViewMode:
		ViewDraw(screen, FaceData, Model)
	case NetworkConfig:
		NetworkModeDraw(screen, &g.debugui)
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
	ebiten.SetWindowDecorated(true)

	Model = models.NewModel()

	_, err := os.ReadDir("./model_files")
	if err != nil {
		os.Mkdir("model_files", os.ModePerm)
	}

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
