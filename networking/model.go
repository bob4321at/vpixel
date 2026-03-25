package networking

import (
	"bytes"
	"encoding/json"
	"main/models"
	"main/triangle"
	"main/utils"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
)

type NetworkedModelJson struct {
	Triangles []NetworkedTriangleJson
	UserID    int
}

type NetworkedTriangleJson struct {
	Points      [3]triangle.Point
	Pixels      []byte
	Color       utils.Vec3
	TexturePath string
}

func TurnTriangleToNetworkedTriangle(tri triangle.Triangle) NetworkedTriangleJson {
	new_triangle := NetworkedTriangleJson{}
	new_triangle.TexturePath = tri.TexturePath

	new_triangle.Points = [3]triangle.Point{}
	for index, point := range tri.Points {
		new_point := triangle.Point{}
		new_point.UvPos = point.UvPos
		new_point.VecPos = point.VecPos
		new_point.Weight = []triangle.Weight{}
		for _, weight := range tri.Points[index].Weight {
			new_weight := triangle.Weight{}
			new_weight.Invert = weight.Invert
			new_weight.Maximum = weight.Maximum
			new_weight.Minimum = weight.Minimum
			new_weight.TestValue = weight.TestValue
			new_weight.Name = weight.Name
			new_weight.Posistion = weight.Posistion
			new_weight.Posistion_Changing = weight.Posistion_Changing
			new_weight.RealValue = weight.RealValue

			new_point.Weight = append(new_point.Weight, new_weight)
		}
		new_triangle.Points[index] = new_point
	}

	new_triangle.Pixels = make([]byte, 360*240*4)
	tri.Texture.ReadPixels(new_triangle.Pixels)

	new_triangle.Color = tri.Color

	return new_triangle
}

func TurnNetworkedTriangleToTriangle(tri NetworkedTriangleJson) triangle.Triangle {
	new_triangle := triangle.NewTriangle(360, 240)
	new_triangle.TexturePath = tri.TexturePath

	new_triangle.Points = [3]triangle.Point{}
	for index, point := range tri.Points {
		new_point := triangle.Point{}
		new_point.UvPos = point.UvPos
		new_point.VecPos = point.VecPos
		new_point.Weight = []triangle.Weight{}
		for _, weight := range tri.Points[index].Weight {
			new_weight := triangle.Weight{}
			new_weight.Invert = weight.Invert
			new_weight.Maximum = weight.Maximum
			new_weight.Minimum = weight.Minimum
			new_weight.TestValue = weight.TestValue
			new_weight.Name = weight.Name
			new_weight.Posistion = weight.Posistion
			new_weight.Posistion_Changing = weight.Posistion_Changing
			new_weight.RealValue = weight.RealValue

			new_point.Weight = append(new_point.Weight, new_weight)
		}
		new_triangle.Points[index] = new_point
	}

	new_triangle.Texture = ebiten.NewImage(360, 240)
	new_triangle.Texture.WritePixels(tri.Pixels)

	new_triangle.Color = tri.Color

	return new_triangle
}

func TurnNetworkedModelToModel(model NetworkedModelJson) models.Model {
	new_model := models.NewModel()

	for i := range model.Triangles {
		triangle := model.Triangles[i]
		new_model.Triangles = append(new_model.Triangles, TurnNetworkedTriangleToTriangle(triangle))
	}

	return new_model
}

func TurnModelToNetworkedModel(model models.Model) NetworkedModelJson {
	new_model := NetworkedModelJson{}

	for i := range model.Triangles {
		triangle := model.Triangles[i]
		new_model.Triangles = append(new_model.Triangles, TurnTriangleToNetworkedTriangle(triangle))
	}

	return new_model
}

func TellServerToChangeModel(model models.Model) {
	networked_model_to_send := TurnModelToNetworkedModel(model)
	networked_model_to_send.UserID = ThisUsersId

	networked_model_to_send_bytes, err := json.Marshal(networked_model_to_send)
	if err != nil {
		panic(err)
	}

	_, err = http.Post("http://"+JoinServerAddress+"/SetUsersModel", "application/json", bytes.NewBuffer(networked_model_to_send_bytes))
	if err != nil {
		panic(err)
	}
}
