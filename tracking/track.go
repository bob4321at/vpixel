package tracking

import (
	"encoding/json"
	"io"
	"main/utils"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
)

type MouthPoints struct {
	LeftCorner  []float64 `json:"left_corner"`
	RightCorner []float64 `json:"right_corner"`
	UpperLip    []float64 `json:"upper_lip"`
	LowerLip    []float64 `json:"lower_lip"`
}

type TrackingData struct {
	Left  []float64   `json:"left"`
	Right []float64   `json:"right"`
	Mouth MouthPoints `json:"mouth"`
	Error string      `json:"error,omitempty"`

	HowOpenMouthIs float64
	HowWideMouthIs float64
	AverageHeadPos utils.Vec2
}

var scaleX = 426.0 / 320.0
var scaleY = 240.0 / 240.0

var WeightOptions = map[string]float64{
	"mouth_open":  0,
	"mouth_width": 0,
}

var BiggestMouth float64 = 1
var WidestMouth float64 = 1

var AverageHeadPos utils.Vec2
var HeadAngle float64

func (face *TrackingData) Update() {
	HttpRequest, err := http.Get("http://localhost:8080")
	if err != nil {
		panic(err)
	}

	FaceDataBytes, err := io.ReadAll(HttpRequest.Body)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(FaceDataBytes, &face)

	face.Left[0] *= scaleX
	face.Left[1] *= scaleY
	face.Right[0] *= scaleX
	face.Right[1] *= scaleY

	face.Mouth.LeftCorner[0] *= scaleX
	face.Mouth.LeftCorner[1] *= scaleY
	face.Mouth.RightCorner[0] *= scaleX
	face.Mouth.RightCorner[1] *= scaleY

	face.Mouth.LowerLip[0] *= scaleX
	face.Mouth.LowerLip[1] *= scaleY
	face.Mouth.UpperLip[0] *= scaleX
	face.Mouth.UpperLip[1] *= scaleY

	face.HowOpenMouthIs = utils.GetDistance(face.Mouth.UpperLip[0], face.Mouth.UpperLip[1], face.Mouth.LowerLip[0], face.Mouth.LowerLip[1])
	face.HowWideMouthIs = utils.GetDistance(face.Mouth.LeftCorner[0], face.Mouth.LeftCorner[1], face.Mouth.RightCorner[0], face.Mouth.RightCorner[1])

	if ebiten.IsKeyPressed(ebiten.KeyC) {
		BiggestMouth = face.HowOpenMouthIs
		WidestMouth = face.HowWideMouthIs
	}

	WeightOptions["mouth_open"] += face.HowOpenMouthIs / BiggestMouth
	WeightOptions["mouth_width"] += face.HowWideMouthIs / WidestMouth
	WeightOptions["mouth_open"] /= 2
	WeightOptions["mouth_width"] /= 2

	face.AverageHeadPos.X = (face.Left[0] + face.Right[0] + (face.Mouth.LeftCorner[0]+face.Mouth.RightCorner[0]+face.Mouth.UpperLip[0]+face.Mouth.LowerLip[0])/4) / 3
	face.AverageHeadPos.Y = (face.Left[1] + face.Right[1] + (face.Mouth.LeftCorner[1]+face.Mouth.RightCorner[1]+face.Mouth.UpperLip[1]+face.Mouth.LowerLip[1])/4) / 3

	HeadAngle += ((utils.Deg2Rad(utils.GetAngle(utils.Vec2{X: face.Mouth.UpperLip[0], Y: face.Mouth.UpperLip[1]}, utils.Vec2{X: face.Left[0], Y: face.Left[1]})) + utils.GetAngle(utils.Vec2{X: face.Mouth.UpperLip[0], Y: face.Mouth.UpperLip[1]}, utils.Vec2{X: face.Right[0], Y: face.Right[1]})) / 2) + 14
	HeadAngle /= 2

	AverageHeadPos = face.AverageHeadPos
}
