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

var NetworkedUserUpscaleImg = ebiten.NewImage(360, 240)

func ViewDraw(screen *ebiten.Image, face tracking.TrackingData, model models.Model) {
	screen.Fill(color.RGBA{0, 255, 0, 255})
	UpscaleImg.Fill(color.RGBA{0, 0, 0, 0})

	networking.UsersModel.Range(func(key, value any) bool {
		op := ebiten.DrawImageOptions{}
		NetworkedUserUpscaleImg.Fill(color.RGBA{0, 0, 0, 0})
		if networking.ThisUsersId == key {
			return true
		}

		var model_to_draw_init models.Model
		var model_to_draw models.Model
		var err bool

		if model_to_draw_init, err = value.(models.Model); !err {
			panic(model_to_draw)
		}

		model_to_draw = *&model_to_draw_init

		tracking_data_for_user, ok := networking.UsersFaceTrackingData.Load(key.(int))
		if !ok {
			return true
		}

		data, ok := tracking_data_for_user.(networking.FaceTrackingNetworked)
		if !ok {
			return true
		}

		NetworkedUserUpscaleImg.Fill(color.White)
		op.GeoM.Reset()
		// op.GeoM.Translate(-66*2.5, -60*2.5)
		// op.GeoM.Rotate(data.HeadAngle * (3.14159 / 180))
		// op.GeoM.Scale(data.DistToEyes/data.AverageDistToEyes, data.DistToEyes/data.AverageDistToEyes)
		// op.GeoM.Translate(66*2.5, 60*2.5)
		// op.GeoM.Scale(3, 3)
		// op.GeoM.Translate(-66*2.5+data.AverageHeadPos.X, -60*2.5+data.AverageHeadPos.Y)

		for _, triangle := range model_to_draw.Triangles {
			for i := range triangle.Points {
				for w := range triangle.Points[i].Weight {
					weight := &triangle.Points[i].Weight[w]

					if !ok {
						break
					}
					weight.RealValue = float64(int(float64(100*data.FaceTrackingData[weight.Name]))) / 100
				}
			}
			triangle.Draw(NetworkedUserUpscaleImg, true)
		}
		screen.DrawImage(NetworkedUserUpscaleImg, &op)

		return true
	})

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
		triangle.Draw(UpscaleImg, true)
	}

	screen.DrawImage(UpscaleImg, &op)
}
