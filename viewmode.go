package main

import (
	"fmt"
	"image/color"
	"main/models"
	"main/networking"
	"main/tracking"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var MovingPlayerMouse bool
var BeforeMoveOffset utils.Vec2
var BeforeMousePos utils.Vec2

func ViewUpdate() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		tracking.PosOffset.Y -= 5
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		tracking.PosOffset.Y += 5
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		tracking.PosOffset.X -= 5
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		tracking.PosOffset.X += 5
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		MovingPlayerMouse = true
		BeforeMoveOffset = tracking.PosOffset
		BeforeMousePos = utils.MousePos
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		MovingPlayerMouse = false
		BeforeMoveOffset = utils.Vec2{}
		BeforeMousePos = utils.Vec2{}
	}

	if MovingPlayerMouse {
		tracking.PosOffset.X = BeforeMoveOffset.X - (utils.MousePos.X-BeforeMousePos.X)*5
		tracking.PosOffset.Y = BeforeMoveOffset.Y - (utils.MousePos.Y-BeforeMousePos.Y)*5
	}

	_, wheel_y := ebiten.Wheel()

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		tracking.ScaleOffset -= 0.03
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			tracking.ScaleOffset = 1
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyE) {
		tracking.ScaleOffset += 0.03
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			tracking.ScaleOffset = 1
		}
	}

	if tracking.ScaleOffset < 0.1 {
		tracking.ScaleOffset = 0.1
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		tracking.RotationOffset -= 0.03
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			tracking.RotationOffset = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyX) {
		tracking.RotationOffset += 0.03
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			tracking.RotationOffset = 0
		}
	}

	if MovingPlayerMouse {
		tracking.RotationOffset += wheel_y / 10
	} else {
		tracking.ScaleOffset += wheel_y / 10
	}
}

var NetworkedLayer *ebiten.Image = ebiten.NewImage(360, 240)

func ViewDraw(screen *ebiten.Image, face tracking.TrackingData, model models.Model) {
	screen.Fill(color.RGBA{0, 255, 0, 255})
	model.RenderLayer.Fill(color.RGBA{0, 0, 0, 0})

	networking.UsersModel.Range(func(key, value any) bool {
		op := ebiten.DrawImageOptions{}
		if networking.ThisUsersId == key {
			return true
		}

		var model_to_draw_init models.Model
		var model_to_draw models.Model

		model_to_draw_init, ok := value.(models.Model)
		if !ok {
			fmt.Printf("Unexpected type in UsersModel[%d]: %T\n", key, value)
			return true
		}

		model_to_draw = model_to_draw_init

		tracking_data_for_user, ok := networking.UsersFaceTrackingData.Load(key.(int))
		if !ok {
			return true
		}

		data, ok := tracking_data_for_user.(networking.FaceTrackingNetworked)
		if !ok {
			return true
		}

		if data.ID == networking.ThisUsersId {
			return true
		}

		op.GeoM.Reset()
		op.GeoM.Translate(-66*2.5, -60*2.5)
		op.GeoM.Rotate(data.HeadAngle*(3.14159/180) + data.RotationOffset)
		op.GeoM.Scale((data.DistToEyes/data.AverageDistToEyes)*data.ScaleOffset, (data.DistToEyes/data.AverageDistToEyes)*data.ScaleOffset)
		op.GeoM.Translate(66*2.5, 60*2.5)
		op.GeoM.Scale(3, 3)
		op.GeoM.Translate(-66*2.5+data.AverageHeadPos.X+data.PosOffset.X, -60*2.5+data.AverageHeadPos.Y+data.PosOffset.Y)
		NetworkedLayer.Clear()

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
			triangle.Draw(NetworkedLayer, true)
		}
		screen.DrawImage(NetworkedLayer, &op)

		return true
	})

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(-66*2.5, -60*2.5)
	op.GeoM.Rotate(tracking.HeadAngle*(3.14159/180) + tracking.RotationOffset)
	op.GeoM.Scale((tracking.DistToEyes/tracking.AverageDistToEyes)*tracking.ScaleOffset, (tracking.DistToEyes/tracking.AverageDistToEyes)*tracking.ScaleOffset)
	op.GeoM.Translate(66*2.5, 60*2.5)
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(-66*2.5+tracking.AverageHeadPos.X+tracking.PosOffset.X, -60*2.5+tracking.AverageHeadPos.Y+tracking.PosOffset.Y)

	for _, triangle := range model.Triangles {
		for i := range triangle.Points {
			for w := range triangle.Points[i].Weight {
				weight := &triangle.Points[i].Weight[w]

				weight.RealValue = float64(int(float64(100*tracking.WeightOptions[weight.Name]))) / 100
			}
		}
		triangle.Draw(model.RenderLayer, true)
	}

	screen.DrawImage(model.RenderLayer, &op)
}
