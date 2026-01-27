package main

import (
	"image/color"
	"main/models"
	"main/tracking"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func ViewUpdate() {
}

func ViewDraw(screen *ebiten.Image, face tracking.TrackingData, model models.Model) {
	screen.Fill(color.RGBA{0, 0, 0, 0})
	UpscaleImg.Fill(color.RGBA{0, 0, 0, 0})

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(-66*2.5+tracking.AverageHeadPos.X, -60*2.5+tracking.AverageHeadPos.Y)

	mouth_openness_string := "MouthOpenness: " + strconv.FormatFloat(tracking.WeightOptions["mouth_open"], 'f', -1, 64)
	ebitenutil.DebugPrintAt(screen, mouth_openness_string, 1, 1)
	mouth_wideness_string := "MouthWideness: " + strconv.FormatFloat(tracking.WeightOptions["mouth_width"], 'f', -1, 64)
	ebitenutil.DebugPrintAt(screen, mouth_wideness_string, 1, 16)

	for _, triangle := range model.Triangles {
		for i := range triangle.Points {
			for w := range triangle.Points[i].Weight {
				weight := &triangle.Points[i].Weight[w]

				weight.RealValue = float64(int(float64(100*tracking.WeightOptions[weight.Name]))) / 100
			}
		}
		triangle_op := ebiten.DrawImageOptions{}
		triangle_op.GeoM.Translate(tracking.AverageHeadPos.X, tracking.AverageHeadPos.Y)
		triangle.Draw(UpscaleImg, true)
	}

	screen.DrawImage(UpscaleImg, &op)
}
