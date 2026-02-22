package networking

import (
	"bytes"
	"encoding/json"
	"image/color"
	"main/models"
	"main/triangle"
	"main/utils"
	"net/http"
)

type NetworkedModelJson struct {
	Triangles []NetworkedTriangleJson
	UserID    uint8
}

type NetworkedTriangleJson struct {
	Points  [3]triangle.Point
	Texture [][]utils.Vec4
	Color   utils.Vec3
}

func TurnTriangleToNetworkedTriangle(tri *triangle.Triangle) NetworkedTriangleJson {
	new_triangle := NetworkedTriangleJson{}

	new_triangle.Points = tri.Points

	new_triangle.Texture = make([][]utils.Vec4, int(tri.TextureSize.Y))

	for i := range int(tri.TextureSize.Y) {
		new_triangle.Texture[i] = make([]utils.Vec4, int(tri.TextureSize.X))
	}

	for y := range len(new_triangle.Texture) {
		for x := range len(new_triangle.Texture[y]) {
			color := tri.Texture.At(x, y)
			color_x, color_y, color_z, color_w := color.RGBA()
			new_triangle.Texture[y][x] = utils.Vec4{X: float64(color_x), Y: float64(color_y), Z: float64(color_z), W: float64(color_w)}
		}
	}

	new_triangle.Color = tri.Color

	return new_triangle
}

func TurnNetworkedTriangleToTriangle(tri *NetworkedTriangleJson) triangle.Triangle {
	new_triangle := triangle.NewTriangle(320, 240)

	new_triangle.Points = tri.Points

	for y := range tri.Texture {
		for x := range tri.Texture[y] {
			pixel_color := tri.Texture[y][x]
			new_triangle.Texture.Set(x, y, color.RGBA{uint8(pixel_color.X), uint8(pixel_color.Y), uint8(pixel_color.Z), uint8(pixel_color.W)})
		}
	}

	new_triangle.Color = tri.Color

	return new_triangle
}

func TurnNetworkedModelToModel(model NetworkedModelJson) models.Model {
	new_model := models.Model{}

	for i := range model.Triangles {
		triangle := &model.Triangles[i]
		new_model.Triangles = append(new_model.Triangles, TurnNetworkedTriangleToTriangle(triangle))
	}

	return new_model
}

func TurnModelToNetworkedModel(model *models.Model) NetworkedModelJson {
	new_model := NetworkedModelJson{}

	for i := range model.Triangles {
		triangle := &model.Triangles[i]
		new_model.Triangles = append(new_model.Triangles, TurnTriangleToNetworkedTriangle(triangle))
	}

	return new_model
}

func TellServerToChangeModel(model *models.Model) {
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
