package main

import (
	"image"
	"image/color"
	"main/models"
	"main/utils"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

func EditUpdate(ui *debugui.DebugUI) {
	var err error
	if utils.ICS, err = ui.Update(func(ctx *debugui.Context) error {
		models.Ctx = ctx
		ctx.Window("Edit", image.Rect(0, 0, 66*5, 240*5), Model.TriangleEditWindow)
		return nil
	}); err != nil {
		panic(err)
	}

	Model.Update()
}

var UpscaleImg = ebiten.NewImage(360, 240)

func EditDraw(screen *ebiten.Image, game Game, ui *debugui.DebugUI) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	UpscaleImg.Fill(color.RGBA{0, 0, 125, 255})

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(66*5, 0)

	Model.Draw(UpscaleImg)

	screen.DrawImage(UpscaleImg, &op)

	ui.Draw(screen)
}
