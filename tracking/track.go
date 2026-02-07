package tracking

import (
	"encoding/json"
	"fmt"
	"io"
	"main/utils"
	"net/http"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type EyePoints struct {
	Top    []float64 `json:"top"`
	Bottom []float64 `json:"bottom"`
}

type MouthPoints struct {
	LeftCorner  []float64 `json:"left_corner"`
	RightCorner []float64 `json:"right_corner"`
	UpperLip    []float64 `json:"upper_lip"`
	LowerLip    []float64 `json:"lower_lip"`
}

type TrackingData struct {
	LeftEye   EyePoints   `json:"left_eye"`
	RightEye  EyePoints   `json:"right_eye"`
	Mouth     MouthPoints `json:"mouth"`
	Error     string      `json:"error,omitempty"`
	Timestamp float64     `json:"timestamp"`

	// Client-side computed fields (not from server JSON)
	HowOpenMouthIs float64
	HowWideMouthIs float64

	LeftEyeOpenness  float64
	RightEyeOpenness float64

	AverageHeadPos     utils.Vec2
	AverageLeftEyePos  utils.Vec2
	AverageRightEyePos utils.Vec2
}

var scaleX = 426.0 / 320.0
var scaleY = 240.0 / 240.0

var WeightOptions = map[string]float64{
	"mouth_open":    0,
	"mouth_width":   0,
	"lefteye_open":  0,
	"righteye_open": 0,
}

var BiggestMouth float64 = 1
var WidestMouth float64 = 1
var BiggestLeftEye float64 = 1
var BiggestRightEye float64 = 1

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

	if strings.Contains(string(FaceDataBytes), "null") {
		return
	}

	json.Unmarshal(FaceDataBytes, &face)

	face.LeftEye.Bottom[0] *= scaleX
	face.LeftEye.Bottom[1] *= scaleY
	face.LeftEye.Top[0] *= scaleX
	face.LeftEye.Top[1] *= scaleY

	face.RightEye.Bottom[0] *= scaleX
	face.RightEye.Bottom[1] *= scaleY
	face.RightEye.Top[0] *= scaleX
	face.RightEye.Top[1] *= scaleY

	face.LeftEyeOpenness = utils.GetDistance(face.LeftEye.Top[0], face.LeftEye.Top[1], face.LeftEye.Bottom[0], face.LeftEye.Bottom[1])
	face.RightEyeOpenness = utils.GetDistance(face.RightEye.Top[0], face.RightEye.Top[1], face.RightEye.Bottom[0], face.RightEye.Bottom[1])

	face.AverageLeftEyePos = utils.Vec2{X: (face.LeftEye.Bottom[0] + face.LeftEye.Top[0]) / 2, Y: (face.LeftEye.Bottom[1] + face.LeftEye.Top[1]) / 2}
	face.AverageRightEyePos = utils.Vec2{X: (face.RightEye.Bottom[0] + face.RightEye.Top[0]) / 2, Y: (face.RightEye.Bottom[1] + face.RightEye.Top[1]) / 2}

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
		BiggestLeftEye = face.LeftEyeOpenness
		BiggestRightEye = face.RightEyeOpenness
	}

	WeightOptions["mouth_open"] += face.HowOpenMouthIs / BiggestMouth
	WeightOptions["mouth_width"] += face.HowWideMouthIs / WidestMouth
	WeightOptions["mouth_open"] /= 2
	WeightOptions["mouth_width"] /= 2

	WeightOptions["lefteye_open"] += face.LeftEyeOpenness / BiggestLeftEye
	WeightOptions["righteye_open"] += face.RightEyeOpenness / BiggestRightEye
	WeightOptions["lefteye_open"] /= 2
	WeightOptions["righteye_open"] /= 2

	face.AverageHeadPos.X = (face.AverageLeftEyePos.X + face.AverageRightEyePos.X + (face.Mouth.LeftCorner[0]+face.Mouth.RightCorner[0]+face.Mouth.UpperLip[0]+face.Mouth.LowerLip[0])/4) / 3
	face.AverageHeadPos.Y = (face.AverageLeftEyePos.Y + face.AverageRightEyePos.Y + (face.Mouth.LeftCorner[1]+face.Mouth.RightCorner[1]+face.Mouth.UpperLip[1]+face.Mouth.LowerLip[1])/4) / 3

	HeadAngle += ((utils.Deg2Rad(utils.GetAngle(utils.Vec2{X: face.Mouth.UpperLip[0], Y: face.Mouth.UpperLip[1]}, face.AverageLeftEyePos)) + utils.GetAngle(utils.Vec2{X: face.Mouth.UpperLip[0], Y: face.Mouth.UpperLip[1]}, face.AverageRightEyePos)) / 2) - 19
	HeadAngle /= 2
	fmt.Println(HeadAngle)

	AverageHeadPos = face.AverageHeadPos
}
