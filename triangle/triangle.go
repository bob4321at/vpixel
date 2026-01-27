package triangle

import (
	"main/utils"
	"math"

	"github.com/bob4321at/textures"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var Triangle_Shader = `//kage:unit pixels
			package main

			var PointVecOne vec2
			var PointVecTwo vec2
			var PointVecThree vec2

			var PointUvOne vec2
			var PointUvTwo vec2
			var PointUvThree vec2

			var Color vec3

			var UvImgWidth int
			var UvImgHeight int

			func angleBetween(a vec2, b vec2) float {
				Offset := a - b

				return atan2(Offset.y, Offset.x)
			}

			func DistanceBetween(a vec2, b vec2) float {
				return sqrt((a.x - b.x)*(a.x - b.x) + (a.y - b.y)*(a.y - b.y));
			}

			func Fragment(targetCoords vec4, srcPos vec2, _ vec4) vec4 {
				TriangleVertsBase := DistanceBetween(PointVecTwo, PointVecOne)
				TriangleVertsAngleToOtherBase := angleBetween(PointVecTwo, PointVecThree)
				TriangleVertsCenterOfBase := vec2(PointVecTwo.x+cos(TriangleVertsAngleToOtherBase)/2, PointVecTwo.y-sin(TriangleVertsAngleToOtherBase))
				TriangleVertsHeight := DistanceBetween(TriangleVertsCenterOfBase, PointVecOne)

				TriangleUvBase := DistanceBetween(PointUvTwo, PointUvOne)
				TriangleUvAngleToOtherBase := angleBetween(PointUvTwo, PointUvThree)
				TriangleUvCenterOfBase := vec2(PointUvTwo.x+cos(TriangleUvAngleToOtherBase)/2, PointUvTwo.y-sin(TriangleUvAngleToOtherBase))
				TriangleUvHeight := DistanceBetween(TriangleUvCenterOfBase, vec2(TriangleUvCenterOfBase.x, PointVecOne.y))

				scale_offset_x := float(TriangleUvBase/TriangleVertsBase)
				scale_offset_y := float(TriangleUvHeight/TriangleVertsHeight)

				dir_to_pixel:= vec2(-1*(targetCoords.y- PointVecOne.y), targetCoords.x - PointVecOne.x)
				DirFromPointOneToTwo := vec2(PointVecOne.x - PointVecTwo.x,PointVecOne.y-PointVecTwo.y)
				SideOne:= dot(dir_to_pixel, DirFromPointOneToTwo)

				DirFromPointOneToThree:= vec2(PointVecOne.x-PointVecThree.x,PointVecOne.y - PointVecThree.y)
				SideTwo := -dot(dir_to_pixel, DirFromPointOneToThree)

				dir_to_pixel= vec2(-1*(targetCoords.y - PointVecThree.y), targetCoords.x - PointVecThree.x)
				DirFromPointTwoToThree:= vec2(PointVecThree.x-PointVecTwo.x,PointVecThree.y - PointVecTwo.y)
				SideThree:= -dot(dir_to_pixel, DirFromPointTwoToThree)

				powa := 1.0

				if SideOne< 0    {
					if SideOne >0 {
					powa = 0
					}
					if SideTwo> 0   {
						powa = 0
					}
					if SideThree> 0   {
						powa = 0
					}
				} else {
					if SideTwo< 0   {
						powa = 0
					}
					if SideThree< 0   {
						powa = 0
					}
				}

				DistFromVertexOne := sqrt(abs((srcPos.x - PointVecOne.x))*abs((srcPos.x - PointVecOne.x)) + abs((srcPos.y - PointVecOne.y))*abs((srcPos.y - PointVecOne.y)));
				DistFromVertexTwo := sqrt(abs((srcPos.x - PointVecTwo.x))*abs((srcPos.x - PointVecTwo.x)) + abs((srcPos.y - PointVecTwo.y))*abs((srcPos.y - PointVecTwo.y)));
				DistFromVertexThree := sqrt(abs((srcPos.x - PointVecThree.x))*abs((srcPos.x - PointVecThree.x)) + abs((srcPos.y - PointVecThree.y))*abs((srcPos.y - PointVecThree.y)));

				AngleFromVertexOne := angleBetween(PointVecOne, srcPos)
				AngleFromVertexTwo := angleBetween(PointVecTwo, srcPos)
				AngleFromVertexThree := angleBetween(PointVecThree, srcPos)

				UvSourceOne := vec2((PointUvOne.x)+(cos(AngleFromVertexOne)*(DistFromVertexOne)*scale_offset_x), (PointUvOne.y)-(sin(AngleFromVertexOne)*(DistFromVertexOne)*scale_offset_y))
				UvSourceTwo := vec2((PointUvTwo.x)+(cos(AngleFromVertexTwo)*(DistFromVertexTwo)*scale_offset_x), (PointUvTwo.y)-(sin(AngleFromVertexTwo)*(DistFromVertexTwo)*scale_offset_y))
				UvSourceThree := vec2((PointUvThree.x)+(cos(AngleFromVertexThree)*(DistFromVertexThree)*scale_offset_x), (PointUvThree.y)-(sin(AngleFromVertexThree)*(DistFromVertexThree)*scale_offset_y))

				TextureSourceCoords := (UvSourceOne/3)+(UvSourceTwo/3)+(UvSourceThree/3)
				TextureFromSource := imageSrc1At(TextureSourceCoords)
				TextureFromSource.w = 1

				return vec4(Color/255, 1)*vec4(powa)
				return TextureFromSource*vec4(powa)
				SideOne -= SideTwo*powa
				SideOne -= SideThree
			}
`

type Weight struct {
	Name               string
	Invert             bool
	Minimum            float64
	Maximum            float64
	TestValue          float64
	RealValue          float64
	Posistion          utils.Vec2
	Posistion_Changing bool
}

type Point struct {
	VecPos utils.Vec2
	UvPos  utils.Vec2
	Weight []Weight
}

type Triangle struct {
	Points  [3]Point
	Image   *textures.Texture
	Texture *ebiten.Image
	Color   utils.Vec3
}

func (triangle *Triangle) Draw(screen *ebiten.Image, test_or_real bool) {
	op := ebiten.DrawImageOptions{}
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = triangle.Image.Img

	ModifiedVecOnePos := triangle.Points[0].VecPos
	ModifiedVecTwoPos := triangle.Points[1].VecPos
	ModifiedVecThreePos := triangle.Points[2].VecPos

	for _, modification := range triangle.Points[0].Weight {
		if modification.Invert {
			angle := utils.GetAngle(ModifiedVecOnePos, utils.Vec2{X: ModifiedVecOnePos.X + modification.Posistion.X, Y: ModifiedVecOnePos.Y + modification.Posistion.Y})
			dist_between_points := utils.GetDistance(ModifiedVecOnePos.X, ModifiedVecOnePos.Y, ModifiedVecOnePos.X+modification.Posistion.X, ModifiedVecOnePos.Y+modification.Posistion.Y)

			if test_or_real {
				ModifiedVecOnePos.X = ModifiedVecOnePos.X - math.Sin(angle)*(dist_between_points) + math.Sin(angle)*(dist_between_points*modification.RealValue)
				ModifiedVecOnePos.Y = ModifiedVecOnePos.Y - math.Cos(angle)*(dist_between_points) + math.Cos(angle)*(dist_between_points*modification.RealValue)
			} else {
				ModifiedVecOnePos.X = ModifiedVecOnePos.X - math.Sin(angle)*(dist_between_points) + math.Sin(angle)*(dist_between_points*modification.TestValue)
				ModifiedVecOnePos.Y = ModifiedVecOnePos.Y - math.Cos(angle)*(dist_between_points) + math.Cos(angle)*(dist_between_points*modification.TestValue)
			}
		} else {
			angle := utils.GetAngle(ModifiedVecOnePos, utils.Vec2{X: ModifiedVecOnePos.X + modification.Posistion.X, Y: ModifiedVecOnePos.Y + modification.Posistion.Y})
			dist_between_points := utils.GetDistance(ModifiedVecOnePos.X, ModifiedVecOnePos.Y, ModifiedVecOnePos.X+modification.Posistion.X, ModifiedVecOnePos.Y+modification.Posistion.Y)

			if test_or_real {
				ModifiedVecOnePos.X -= math.Sin(angle) * (dist_between_points * modification.RealValue)
				ModifiedVecOnePos.Y -= math.Cos(angle) * (dist_between_points * modification.RealValue)
			} else {
				ModifiedVecOnePos.X -= math.Sin(angle) * (dist_between_points * modification.TestValue)
				ModifiedVecOnePos.Y -= math.Cos(angle) * (dist_between_points * modification.TestValue)
			}
		}
	}
	for _, modification := range triangle.Points[1].Weight {
		if modification.Invert {
			angle := utils.GetAngle(ModifiedVecTwoPos, utils.Vec2{X: ModifiedVecTwoPos.X + modification.Posistion.X, Y: ModifiedVecTwoPos.Y + modification.Posistion.Y})
			dist_between_points := utils.GetDistance(ModifiedVecTwoPos.X, ModifiedVecTwoPos.Y, ModifiedVecTwoPos.X+modification.Posistion.X, ModifiedVecTwoPos.Y+modification.Posistion.Y)

			if test_or_real {
				ModifiedVecTwoPos.X = ModifiedVecTwoPos.X - math.Sin(angle)*(dist_between_points) + math.Sin(angle)*(dist_between_points*modification.RealValue)
				ModifiedVecTwoPos.Y = ModifiedVecTwoPos.Y - math.Cos(angle)*(dist_between_points) + math.Cos(angle)*(dist_between_points*modification.RealValue)
			} else {
				ModifiedVecTwoPos.X = ModifiedVecTwoPos.X - math.Sin(angle)*(dist_between_points) + math.Sin(angle)*(dist_between_points*modification.TestValue)
				ModifiedVecTwoPos.Y = ModifiedVecTwoPos.Y - math.Cos(angle)*(dist_between_points) + math.Cos(angle)*(dist_between_points*modification.TestValue)
			}
		} else {
			angle := utils.GetAngle(ModifiedVecTwoPos, utils.Vec2{X: ModifiedVecTwoPos.X + modification.Posistion.X, Y: ModifiedVecTwoPos.Y + modification.Posistion.Y})
			dist_between_points := utils.GetDistance(ModifiedVecTwoPos.X, ModifiedVecTwoPos.Y, ModifiedVecTwoPos.X+modification.Posistion.X, ModifiedVecTwoPos.Y+modification.Posistion.Y)

			if test_or_real {
				ModifiedVecTwoPos.X -= math.Sin(angle) * (dist_between_points * modification.RealValue)
				ModifiedVecTwoPos.Y -= math.Cos(angle) * (dist_between_points * modification.RealValue)
			} else {
				ModifiedVecTwoPos.X -= math.Sin(angle) * (dist_between_points * modification.TestValue)
				ModifiedVecTwoPos.Y -= math.Cos(angle) * (dist_between_points * modification.TestValue)
			}
		}
	}
	for _, modification := range triangle.Points[2].Weight {
		if modification.Invert {
			angle := utils.GetAngle(ModifiedVecThreePos, utils.Vec2{X: ModifiedVecThreePos.X + modification.Posistion.X, Y: ModifiedVecThreePos.Y + modification.Posistion.Y})
			dist_between_points := utils.GetDistance(ModifiedVecThreePos.X, ModifiedVecThreePos.Y, ModifiedVecThreePos.X+modification.Posistion.X, ModifiedVecThreePos.Y+modification.Posistion.Y)

			if test_or_real {
				ModifiedVecThreePos.X = ModifiedVecThreePos.X - math.Sin(angle)*(dist_between_points) + math.Sin(angle)*(dist_between_points*modification.RealValue)
				ModifiedVecThreePos.Y = ModifiedVecThreePos.Y - math.Cos(angle)*(dist_between_points) + math.Cos(angle)*(dist_between_points*modification.RealValue)
			} else {
				ModifiedVecThreePos.X = ModifiedVecThreePos.X - math.Sin(angle)*(dist_between_points) + math.Sin(angle)*(dist_between_points*modification.TestValue)
				ModifiedVecThreePos.Y = ModifiedVecThreePos.Y - math.Cos(angle)*(dist_between_points) + math.Cos(angle)*(dist_between_points*modification.TestValue)
			}
		} else {
			angle := utils.GetAngle(ModifiedVecOnePos, utils.Vec2{X: ModifiedVecOnePos.X + modification.Posistion.X, Y: ModifiedVecOnePos.Y + modification.Posistion.Y})
			dist_between_points := utils.GetDistance(ModifiedVecOnePos.X, ModifiedVecOnePos.Y, ModifiedVecOnePos.X+modification.Posistion.X, ModifiedVecOnePos.Y+modification.Posistion.Y)

			if test_or_real {
				ModifiedVecThreePos.X -= math.Sin(angle) * (dist_between_points * modification.RealValue)
				ModifiedVecThreePos.Y -= math.Cos(angle) * (dist_between_points * modification.RealValue)
			} else {
				ModifiedVecThreePos.X -= math.Sin(angle) * (dist_between_points * modification.TestValue)
				ModifiedVecThreePos.Y -= math.Cos(angle) * (dist_between_points * modification.TestValue)
			}
		}
	}

	if triangle.Texture != nil {
		triangle.Image.SetUniforms(map[string]any{
			"PointVecOne":   []float32{float32(ModifiedVecOnePos.X), float32(ModifiedVecOnePos.Y)},
			"PointVecTwo":   []float32{float32(ModifiedVecTwoPos.X), float32(ModifiedVecTwoPos.Y)},
			"PointVecThree": []float32{float32(ModifiedVecThreePos.X), float32(ModifiedVecThreePos.Y)},

			"PointUvOne":   []float32{float32(triangle.Points[0].UvPos.X), float32(triangle.Points[0].UvPos.Y)},
			"PointUvTwo":   []float32{float32(triangle.Points[1].UvPos.X), float32(triangle.Points[1].UvPos.Y)},
			"PointUvThree": []float32{float32(triangle.Points[2].UvPos.X), float32(triangle.Points[2].UvPos.Y)},

			"Color": []float32{float32(triangle.Color.X), float32(triangle.Color.Y), float32(triangle.Color.Z)},

			"UvImgWidth":  triangle.Texture.Bounds().Dx(),
			"UvImgHeight": triangle.Texture.Bounds().Dy(),
		})
		opts.Images[1] = triangle.Texture
	} else {
		triangle.Image.SetUniforms(map[string]any{
			"PointVecOne":   []float32{float32(ModifiedVecOnePos.X), float32(ModifiedVecOnePos.Y)},
			"PointVecTwo":   []float32{float32(ModifiedVecTwoPos.X), float32(ModifiedVecTwoPos.Y)},
			"PointVecThree": []float32{float32(ModifiedVecThreePos.X), float32(ModifiedVecThreePos.Y)},

			"PointUvOne":   []float32{float32(triangle.Points[0].UvPos.X), float32(triangle.Points[0].UvPos.Y)},
			"PointUvTwo":   []float32{float32(triangle.Points[1].UvPos.X), float32(triangle.Points[1].UvPos.Y)},
			"PointUvThree": []float32{float32(triangle.Points[2].UvPos.X), float32(triangle.Points[2].UvPos.Y)},

			"Color": []float32{float32(triangle.Color.X), float32(triangle.Color.Y), float32(triangle.Color.Z)},

			"UvImgWidth":  0,
			"UvImgHeight": 0,
		})
	}
	opts.Uniforms = triangle.Image.Uniforms
	opts.GeoM = op.GeoM
	screen.DrawRectShader(triangle.Image.Img.Bounds().Dx(), triangle.Image.Img.Bounds().Dy(), triangle.Image.Shader, opts)
}

func (triangle *Triangle) SetPointsVectorPos(point_1, point_2, point_3 utils.Vec2) {
	triangle.Points[0].VecPos = point_1
	triangle.Points[1].VecPos = point_2
	triangle.Points[2].VecPos = point_3
}

func (triangle *Triangle) SetPointsUvPos(point_1, point_2, point_3 utils.Vec2) {
	triangle.Points[0].UvPos = point_1
	triangle.Points[1].UvPos = point_2
	triangle.Points[2].UvPos = point_3
}

func (triangle *Triangle) SetColors(color utils.Vec3) {
	triangle.Color = color
}

func (triangle *Triangle) SetTexture(path string) {
	temp_img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}

	triangle.Texture = triangle.Image.Img
	triangle.Texture.DrawImage(temp_img, &ebiten.DrawImageOptions{})
}

func NewTriangle(screen_res_width, screen_res_height int) (triangle Triangle) {
	triangle.Image = textures.NewTexture("./empty.png", Triangle_Shader)
	triangle.Image.Img = ebiten.NewImage(screen_res_width, screen_res_height)

	triangle.SetPointsVectorPos(utils.Vec2{X: float64(screen_res_width) / 4 * 3, Y: float64(screen_res_height) / 4}, utils.Vec2{X: float64(screen_res_width) / 2, Y: float64(screen_res_height) / 4 * 3}, utils.Vec2{X: float64(screen_res_width) / 4, Y: float64(screen_res_height) / 4})

	return triangle
}
