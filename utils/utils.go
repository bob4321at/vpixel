package utils

import (
	"math"

	"github.com/ebitengine/debugui"
)

var Tick int

type Vec2 struct {
	X, Y float64
}

type Vec3 struct {
	X, Y, Z float64
}

func GetDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
}

func GetAngle(point_1, point_2 Vec2) float64 {
	offset_x := point_1.X - point_2.X
	offset_y := point_1.Y - point_2.Y

	return math.Atan2(offset_x, offset_y)
}

func Deg2Rad(num float64) float64 {
	return num * (180 / 3.14159)
}

var MousePos Vec2
var ICS debugui.InputCapturingState
