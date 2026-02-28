package main

import (
	"image/color"
	"main/models"
	"main/networking"
	"main/tracking"

	"github.com/hajimehoshi/ebiten/v2"
)

func ViewUpdate() {
}

func ViewDraw(screen *ebiten.Image, face tracking.TrackingData, model models.Model) {
	screen.Fill(color.RGBA{0, 255, 0, 255})
	UpscaleImg.Fill(color.RGBA{0, 0, 0, 0})

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(-66*2.5, -60*2.5)
	op.GeoM.Rotate(tracking.HeadAngle * (3.14159 / 180))
	op.GeoM.Scale(tracking.DistToEyes/tracking.AverageDistToEyes, tracking.DistToEyes/tracking.AverageDistToEyes)
	op.GeoM.Translate(66*2.5, 60*2.5)
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(-66*2.5+tracking.AverageHeadPos.X, -60*2.5+tracking.AverageHeadPos.Y)

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

	networking.UsersModel.Range(func(key, value any) bool {
		if networking.ThisUsersId == key {
			return true
		}

		var model_to_draw models.Model
		var err bool

		if model_to_draw, err = value.(models.Model); !err {
			panic(model_to_draw)
		}

		tracking_data_for_user, ok := networking.UsersFaceTrackingData.Load(key)
		data, ok := tracking_data_for_user.(networking.FaceTrackingNetworked)
		data.HeadAngle = 0
		if ok {
			op.GeoM.Reset()
			op.GeoM.Translate(-66*2.5, -60*2.5)
			op.GeoM.Rotate(data.HeadAngle * (3.14159 / 180))
			op.GeoM.Scale(data.DistToEyes/data.AverageDistToEyes, data.DistToEyes/data.AverageDistToEyes)
			op.GeoM.Translate(66*2.5, 60*2.5)
			op.GeoM.Scale(3, 3)
			op.GeoM.Translate(-66*2.5+data.AverageHeadPos.X, -60*2.5+data.AverageHeadPos.Y)

			for _, triangle := range model_to_draw.Triangles {
				for i := range triangle.Points {
					for w := range triangle.Points[i].Weight {
						weight := &triangle.Points[i].Weight[w]

						if ok {
							weight.RealValue = float64(int(float64(100*data.FaceTrackingData[weight.Name]))) / 100
						}
					}
				}
				triangle_op := ebiten.DrawImageOptions{}
				triangle_op.GeoM.Translate(tracking.AverageHeadPos.X, tracking.AverageHeadPos.Y)
				triangle.Draw(UpscaleImg, true)
			}
		}
		return true
	})

	screen.DrawImage(UpscaleImg, &op)
}
