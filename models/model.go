package models

import (
	"fmt"
	"image"
	"image/color"
	"main/triangle"
	"main/utils"
	"strconv"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Model struct {
	Triangles []triangle.Triangle
}

type ModelJson struct {
	Triangles []triangle.TriangleJson
}

func (model *Model) Encode() ModelJson {
	new_model := ModelJson{}

	for _, triangle := range model.Triangles {
		new_model.Triangles = append(new_model.Triangles, triangle.Encode())
	}

	return new_model
}

func (model *ModelJson) Decode(new_model *Model) {
	new_model.Triangles = nil
	fmt.Println(model)
	for _, triangle := range model.Triangles {
		new_model.Triangles = append(new_model.Triangles, triangle.Decode())
	}
}

var SelectedVertex int
var SelectedTriangle int = 0

var LastSelectedVertex int

var Ctx *debugui.Context

var WeightConfigOpen bool

var SaveName string

var CopyPasteColorIndex int

var TexturePathEdit string

func (model *Model) TriangleEditWindow(layout debugui.ContainerLayout) {
	if SelectedTriangle == -1 {
		return
	}

	if len(model.Triangles) != 0 {
		Ctx.SliderF(&model.Triangles[SelectedTriangle].Color.X, 0, 255, 1, 3)
		Ctx.SliderF(&model.Triangles[SelectedTriangle].Color.Y, 0, 255, 1, 3)
		Ctx.SliderF(&model.Triangles[SelectedTriangle].Color.Z, 0, 255, 1, 3)

		if WeightConfigOpen {
			Weights := model.Triangles[SelectedTriangle].Points[LastSelectedVertex].Weight

			Ctx.Window("Weight Config", image.Rect(66*5, 0, 66*5+300, 200), func(layout debugui.ContainerLayout) {
				Ctx.Loop(len(Weights), func(i int) {
					weight := &Weights[i]
					Ctx.TextField(&weight.Name)
					Ctx.Checkbox(&weight.Invert, "Invert")
					Ctx.SliderF(&weight.Minimum, 0, 1, 0.01, 3)
					Ctx.SliderF(&weight.Maximum, 0, 1, 0.01, 3)
					Ctx.SliderF(&weight.TestValue, 0, 1, 0.01, 3)
					weight_pos_x := strconv.Itoa(int(weight.Posistion.X))
					weight_pos_y := strconv.Itoa(int(weight.Posistion.X))
					Ctx.Text(weight_pos_x)
					Ctx.Text(weight_pos_y)
					Ctx.Button("Set Weight Pos").On(func() {
						weight.Posistion_Changing = true
					})

					if weight.Posistion_Changing {
						if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
							weight.Posistion = utils.Vec2{X: utils.MousePos.X - model.Triangles[SelectedTriangle].Points[LastSelectedVertex].VecPos.X, Y: utils.MousePos.Y - model.Triangles[SelectedTriangle].Points[LastSelectedVertex].VecPos.Y}
							weight.Posistion_Changing = false
						}
					}
				})
				Ctx.Button("Add Empty Weight").On(func() {
					model.Triangles[SelectedTriangle].Points[LastSelectedVertex].Weight = append(Weights, triangle.Weight{})
				})
			})
		}

		UvEditPopupId := Ctx.Popup(func(layout debugui.ContainerLayout, popupID debugui.PopupID) {
			Ctx.NumberFieldF(&model.Triangles[SelectedTriangle].Points[LastSelectedVertex].UvPos.X, 1, 0)
			Ctx.NumberFieldF(&model.Triangles[SelectedTriangle].Points[LastSelectedVertex].UvPos.Y, 1, 0)
			Ctx.Button("Close").On(func() {
				Ctx.ClosePopup(popupID)
			})
		})

		Ctx.Button("Open Weight Config").On(func() {
			WeightConfigOpen = !WeightConfigOpen
		})

		Ctx.Button("Open Uv Config").On(func() {
			Ctx.OpenPopup(UvEditPopupId)
		})

		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			if ebiten.IsKeyPressed(ebiten.KeyControl) {
				CopyPasteColorIndex = SelectedTriangle
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyV) {
			if ebiten.IsKeyPressed(ebiten.KeyControl) {
				model.Triangles[SelectedTriangle].SetColors(model.Triangles[CopyPasteColorIndex].Color)
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			closest_point := utils.Vec2{X: -1, Y: -1}
			closest_dist := 10000000000000.0

			for i, other_traingles := range model.Triangles {
				if &model.Triangles[i] != &model.Triangles[SelectedTriangle] {
					for _, vertex := range other_traingles.Points {
						dist := utils.GetDistance(
							model.Triangles[SelectedTriangle].Points[SelectedVertex].VecPos.X,
							model.Triangles[SelectedTriangle].Points[SelectedVertex].VecPos.Y,
							vertex.VecPos.X,
							vertex.VecPos.Y,
						)

						if dist < float64(closest_dist) {
							closest_point = vertex.VecPos
							closest_dist = dist
						}
					}
				}
			}

			nothing_found := utils.Vec2{X: -1, Y: -1}

			if closest_point != nothing_found {
				model.Triangles[SelectedTriangle].Points[SelectedVertex].VecPos = closest_point
			}
		}
	}

	Ctx.TextField(&TexturePathEdit)
	Ctx.Button("Set Texture").On(func() {
		model.Triangles[SelectedTriangle].SetTexture("./textures/" + TexturePathEdit + ".png")
		model.Triangles[SelectedTriangle].SetPointsUvPos(utils.Vec2{0, 32}, utils.Vec2{16, 0}, utils.Vec2{32, 32})
		fmt.Println(model.Triangles[SelectedTriangle])
	})

	Ctx.Button("New Triangle").On(func() {
		model.Triangles = append(model.Triangles, triangle.NewTriangle(360, 240))
	})
}

func (model *Model) Update() {
	if utils.ICS == 0 {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			LastSelectedVertex = SelectedVertex
		}
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			SelectedTriangle = -1
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyDelete) && SelectedTriangle != -1 {
			utils.RemoveArrayElement(SelectedTriangle, &model.Triangles)
			SelectedTriangle = -1
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyN) {
			model.Triangles = append(model.Triangles, triangle.NewTriangle(360, 240))
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyV) && ebiten.IsKeyPressed(ebiten.KeyShift) {
			LastSelectedVertex += 1
			if LastSelectedVertex > 2 {
				LastSelectedVertex = 0
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				if SelectedTriangle-1 >= 0 {
					SelectedTriangle -= 1
				}
			} else {
				if SelectedTriangle+1 < len(model.Triangles) {
					SelectedTriangle += 1
				}
			}
		}
		if len(model.Triangles) != 0 {
			if ebiten.IsKeyPressed(ebiten.KeyR) {
				if ebiten.IsKeyPressed(ebiten.KeyShift) {
					if model.Triangles[SelectedTriangle].Color.X-1 >= 0 {
						model.Triangles[SelectedTriangle].Color.X -= 1
					}
				} else {
					if model.Triangles[SelectedTriangle].Color.X+1 <= 255 {
						model.Triangles[SelectedTriangle].Color.X += 1
					}
				}
			}
		}
		if len(model.Triangles) != 0 {
			if ebiten.IsKeyPressed(ebiten.KeyG) {
				if ebiten.IsKeyPressed(ebiten.KeyShift) {
					if model.Triangles[SelectedTriangle].Color.Y-1 >= 0 {
						model.Triangles[SelectedTriangle].Color.Y -= 1
					}
				} else {
					if model.Triangles[SelectedTriangle].Color.Y+1 <= 255 {
						model.Triangles[SelectedTriangle].Color.Y += 1
					}
				}
			}
		}
		if len(model.Triangles) != 0 {
			if ebiten.IsKeyPressed(ebiten.KeyB) {
				if ebiten.IsKeyPressed(ebiten.KeyShift) {
					if model.Triangles[SelectedTriangle].Color.Z-1 >= 0 {
						model.Triangles[SelectedTriangle].Color.Z -= 1
					}
				} else {
					if model.Triangles[SelectedTriangle].Color.Z+1 <= 255 {
						model.Triangles[SelectedTriangle].Color.Z += 1
					}
				}
			}
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			model.Triangles[SelectedTriangle].Points[LastSelectedVertex].VecPos = utils.Vec2{X: utils.MousePos.X, Y: utils.MousePos.Y}
		} else {
			DistFromMouse := 1000000000.0

			if len(model.Triangles) != 0 && SelectedTriangle != -1 {
				triangle := model.Triangles[SelectedTriangle]
				for p, Point := range triangle.Points {
					dist := utils.GetDistance(utils.MousePos.X, utils.MousePos.Y, Point.VecPos.X, Point.VecPos.Y)
					if dist < float64(DistFromMouse) {
						SelectedVertex = p
						DistFromMouse = dist
					}
				}
			}
		}
	}
}

func (model *Model) Draw(screen *ebiten.Image) {
	for _, triangle := range model.Triangles {
		triangle.Draw(screen, false)
	}

	if len(model.Triangles) != 0 && SelectedTriangle != -1 {
		selected_tri := &model.Triangles[SelectedTriangle]
		selected_tri.Draw(screen, false)
		vector.StrokeLine(screen, float32(selected_tri.Points[0].VecPos.X), float32(selected_tri.Points[0].VecPos.Y), float32(selected_tri.Points[1].VecPos.X), float32(selected_tri.Points[1].VecPos.Y), 2, color.RGBA{125, 125, 125, 255}, false)
		vector.StrokeLine(screen, float32(selected_tri.Points[0].VecPos.X), float32(selected_tri.Points[0].VecPos.Y), float32(selected_tri.Points[2].VecPos.X), float32(selected_tri.Points[2].VecPos.Y), 2, color.RGBA{125, 125, 125, 255}, false)
		vector.StrokeLine(screen, float32(selected_tri.Points[2].VecPos.X), float32(selected_tri.Points[2].VecPos.Y), float32(selected_tri.Points[1].VecPos.X), float32(selected_tri.Points[1].VecPos.Y), 2, color.RGBA{125, 125, 125, 255}, false)

		vertex := &selected_tri.Points[LastSelectedVertex]
		vector.StrokeCircle(screen, float32(vertex.VecPos.X), float32(vertex.VecPos.Y), 2, 2, color.RGBA{100, 100, 100, 255}, false)

		for _, weight := range vertex.Weight {
			vector.StrokeCircle(screen, float32(vertex.VecPos.X+weight.Posistion.X), float32(vertex.VecPos.Y+weight.Posistion.Y), 1, 1, color.RGBA{255, 100, 100, 255}, false)
		}

		if selected_tri.Texture != nil {
			screen.DrawImage(selected_tri.Texture, &ebiten.DrawImageOptions{})
		}
	}
}

func NewModel() (model Model) {
	return model
}
