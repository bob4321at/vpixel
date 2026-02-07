package main

import (
	"image/color"
	"main/tracking"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var CameraOutput *ebiten.Image = ebiten.NewImage(426*3, 240*3)

var AverageFacePos utils.Vec2

func TrackDraw(screen *ebiten.Image, face tracking.TrackingData) {
	screen.Fill(color.RGBA{200, 200, 200, 255})
	CameraOutput.Clear()

	vector.StrokeLine(CameraOutput, float32(FaceData.LeftEye.Top[0]), float32(FaceData.LeftEye.Top[1]), float32(FaceData.LeftEye.Bottom[0]), float32(FaceData.LeftEye.Bottom[1]), 2, color.Black, false)
	vector.StrokeLine(CameraOutput, float32(FaceData.RightEye.Top[0]), float32(FaceData.RightEye.Top[1]), float32(FaceData.RightEye.Bottom[0]), float32(FaceData.RightEye.Bottom[1]), 2, color.Black, false)

	vector.StrokeLine(CameraOutput, float32(FaceData.Mouth.LeftCorner[0]), float32(FaceData.Mouth.LeftCorner[1]), float32(FaceData.Mouth.UpperLip[0]), float32(FaceData.Mouth.UpperLip[1]), 2, color.Black, false)
	vector.StrokeLine(CameraOutput, float32(FaceData.Mouth.LeftCorner[0]), float32(FaceData.Mouth.LeftCorner[1]), float32(FaceData.Mouth.LowerLip[0]), float32(FaceData.Mouth.LowerLip[1]), 2, color.Black, false)

	vector.StrokeLine(CameraOutput, float32(FaceData.Mouth.RightCorner[0]), float32(FaceData.Mouth.RightCorner[1]), float32(FaceData.Mouth.UpperLip[0]), float32(FaceData.Mouth.UpperLip[1]), 2, color.Black, false)
	vector.StrokeLine(CameraOutput, float32(FaceData.Mouth.RightCorner[0]), float32(FaceData.Mouth.RightCorner[1]), float32(FaceData.Mouth.LowerLip[0]), float32(FaceData.Mouth.LowerLip[1]), 2, color.Black, false)

	screen.DrawImage(CameraOutput, &ebiten.DrawImageOptions{})
}
