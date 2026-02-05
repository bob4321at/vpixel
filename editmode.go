package main

import (
	"encoding/json"
	"image"
	"image/color"
	"main/models"
	"main/utils"
	"os"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

func EditUpdate(ui *debugui.DebugUI, model *models.Model) {
	var err error
	if utils.ICS, err = ui.Update(func(ctx *debugui.Context) error {
		models.Ctx = ctx
		ctx.Window("Edit", image.Rect(0, 0, 66*5, 240*2.6), model.TriangleEditWindow)
		ctx.Window("Base", image.Rect(0, 240*2.6, 66*5, 240*5), func(layout debugui.ContainerLayout) {
			ctx.TextField(&models.SaveName)
			ctx.Button("Save").On(func() {
				f, err := os.Create("./model_files/" + models.SaveName)
				if err != nil {
					panic(err)
				}

				save_model_json := model.Encode()

				bytes, err := json.Marshal(&save_model_json)
				if err != nil {
					panic(err)
				}

				f.Write(bytes)

				f.Close()
			})
			ctx.Button("Load").On(func() {
				f, err := os.ReadFile("./model_files/" + models.SaveName)
				if err != nil {
					panic(err)
				}

				load_model_json := models.ModelJson{}

				json.Unmarshal(f, &load_model_json)

				load_model_json.Decode(model)
			})
		})
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
