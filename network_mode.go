package main

import (
	"image"
	"image/color"
	"main/models"
	"main/networking"
	"main/utils"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

func NetworkModeUpdate(ui *debugui.DebugUI, model *models.Model) {
	var err error
	if utils.ICS, err = ui.Update(func(ctx *debugui.Context) error {
		ctx.Window("Networking", image.Rect(0, 0, 426*3, 240*3), func(layout debugui.ContainerLayout) {
			if networking.InitialNetworkStarted {
				ctx.Button("Send New Model").On(func() {
					networking.TellServerToChangeModel(model)
				})
				return
			}

			ctx.Button("Start server").On(func() {
				go networking.StartServer()
			})
			ctx.TextField(&networking.JoinServerAddress)
			ctx.Button("Join server").On(func() {
				go networking.JoinServer()
			})
		})
		return nil
	}); err != nil {
		panic(err)
	}
}

func NetworkModeDraw(screen *ebiten.Image, ui *debugui.DebugUI) {
	screen.Fill(color.RGBA{0, 125, 0, 255})
	ui.Draw(screen)
}
