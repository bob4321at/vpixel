package utils

import (
	"errors"
	"fmt"
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

func RemoveArrayElement[T any](index_to_remove int, slice *[]T) {
	*slice = append((*slice)[:index_to_remove], (*slice)[index_to_remove+1:]...)
}

func MoveElement[T any](slice []T, from, to int) error {
	n := len(slice)

	if n == 0 {
		return errors.New("cannot move element in empty slice")
	}
	if from < 0 || from >= n {
		return fmt.Errorf("from index %d out of bounds [0, %d]", from, n-1)
	}
	if to < 0 || to >= n {
		return fmt.Errorf("to index %d out of bounds [0, %d]", to, n-1)
	}
	if from == to {
		return nil
	}

	value := slice[from]

	if from < to {
		copy(slice[from:], slice[from+1:to+1])
		slice[to] = value
	} else { // from > to
		copy(slice[to+1:], slice[to:from])
		slice[to] = value
	}

	return nil
}

var MousePos Vec2
var ICS debugui.InputCapturingState
