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

func AreaOfTriangle(a vec2, b vec2, c vec2) float {
    return (b.x - a.x) * (c.y - a.y) - (b.y - a.y) * (c.x - a.x)
}

func Fragment(targetCoords vec4, srcPos vec2, _ vec4) vec4 {
    total_area := AreaOfTriangle(PointVecOne, PointVecTwo, PointVecThree)

    if abs(total_area) < 0.0001 {
        return vec4(0, 0, 0, 0)
    }

    alpha := AreaOfTriangle(targetCoords.xy, PointVecTwo, PointVecThree) / total_area
	beta  := AreaOfTriangle(PointVecOne, targetCoords.xy, PointVecThree) / total_area
    gamma := 1.0 - alpha - beta  // More stable than computing 3rd area

    if alpha < 0.0 || beta < 0.0 || gamma < 0.0 {
        return vec4(0, 0, 0, 0)
    }

    UvSource := alpha * PointUvOne + beta * PointUvTwo + gamma * PointUvThree

    tex := imageSrc1At(UvSource+imageSrc0Origin())

	dir_to_pixel:= vec2(-1*(targetCoords.y- PointVecOne.y), targetCoords.x - PointVecOne.x)
	DirFromPointOneToTwo := vec2(PointVecOne.x - PointVecTwo.x,PointVecOne.y-PointVecTwo.y)
	SideOne:= dot(dir_to_pixel, DirFromPointOneToTwo)

	DirFromPointOneToThree:= vec2(PointVecOne.x-PointVecThree.x,PointVecOne.y - PointVecThree.y)
	SideTwo := -dot(dir_to_pixel, DirFromPointOneToThree)

	dir_to_pixel= vec2(-1*(targetCoords.y - PointVecThree.y), targetCoords.x - PointVecThree.x)
	DirFromPointTwoToThree:= vec2(PointVecThree.x-PointVecTwo.x,PointVecThree.y - PointVecTwo.y)
	SideThree:= -dot(dir_to_pixel, DirFromPointTwoToThree)

	powa := 1.0

	if SideOne < 0    {
		if SideTwo > 0   {
			powa = 0
		}
		if SideThree > 0   {
			powa = 0
		}
	} else {
		if SideTwo < 0   {
			powa = 0
		}
		if SideThree < 0   {
			powa = 0
		}
	}

	if imageSrc1Size().x > 0 {
    	return tex*vec4(Color/255, powa*tex.w)
	} else {
    	return vec4(Color/255, powa)
	}
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
	Points      [3]Point
	Image       *textures.Texture
	Texture     *ebiten.Image
	TexturePath string
	Color       utils.Vec3
}

type TriangleJson struct {
	Points      [3]Point
	TexturePath string
	Color       utils.Vec3
}

func (triangle *TriangleJson) Decode() Triangle {
	new_triangle := NewTriangle(360, 240)

	new_triangle.Points = triangle.Points

	new_triangle.SetTexture(triangle.TexturePath)
	new_triangle.SetPointsUvPos(
		triangle.Points[0].UvPos,
		triangle.Points[1].UvPos,
		triangle.Points[2].UvPos,
	)

	new_triangle.Color = triangle.Color

	return new_triangle
}

func (triangle *Triangle) Encode() TriangleJson {
	new_triangle := TriangleJson{}

	new_triangle.Points = triangle.Points
	new_triangle.TexturePath = triangle.TexturePath
	new_triangle.Color = triangle.Color

	return new_triangle
}

func (triangle *Triangle) Draw(screen *ebiten.Image, test_or_real bool) {
	op := ebiten.DrawImageOptions{}
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = triangle.Image.Img
	opts.Images[1] = triangle.Texture

	ModifiedPointPoses := []utils.Vec2{
		triangle.Points[0].VecPos,
		triangle.Points[1].VecPos,
		triangle.Points[2].VecPos,
	}

	for i := range triangle.Points {
		ModifiedPoint := &triangle.Points[i]
		ModifiedPointPos := &ModifiedPointPoses[i]
		for _, modification := range ModifiedPoint.Weight {
			if modification.RealValue > modification.Maximum {
				modification.RealValue = modification.Maximum
			} else if modification.RealValue < modification.Minimum {
				modification.RealValue = modification.Minimum
			}
			if modification.Invert {
				angle := utils.GetAngle(*ModifiedPointPos, utils.Vec2{X: ModifiedPointPos.X + modification.Posistion.X, Y: ModifiedPointPos.Y + modification.Posistion.Y})
				dist_between_points := utils.GetDistance(ModifiedPointPos.X, ModifiedPointPos.Y, ModifiedPointPos.X+modification.Posistion.X, ModifiedPointPos.Y+modification.Posistion.Y)

				if test_or_real {
					ModifiedPointPos.X = ModifiedPointPos.X - math.Sin(angle)*(dist_between_points) + math.Sin(angle)*(dist_between_points*modification.RealValue)
					ModifiedPointPos.Y = ModifiedPointPos.Y - math.Cos(angle)*(dist_between_points) + math.Cos(angle)*(dist_between_points*modification.RealValue)
				} else {
					ModifiedPointPos.X = ModifiedPointPos.X - math.Sin(angle)*(dist_between_points) + math.Sin(angle)*(dist_between_points*modification.TestValue)
					ModifiedPointPos.Y = ModifiedPointPos.Y - math.Cos(angle)*(dist_between_points) + math.Cos(angle)*(dist_between_points*modification.TestValue)
				}
			} else {
				angle := utils.GetAngle(*ModifiedPointPos, utils.Vec2{X: ModifiedPointPos.X + modification.Posistion.X, Y: ModifiedPointPos.Y + modification.Posistion.Y})
				dist_between_points := utils.GetDistance(ModifiedPointPos.X, ModifiedPointPos.Y, ModifiedPointPos.X+modification.Posistion.X, ModifiedPointPos.Y+modification.Posistion.Y)

				if test_or_real {
					ModifiedPointPos.X -= math.Sin(angle) * (dist_between_points * modification.RealValue)
					ModifiedPointPos.Y -= math.Cos(angle) * (dist_between_points * modification.RealValue)
				} else {
					ModifiedPointPos.X -= math.Sin(angle) * (dist_between_points * modification.TestValue)
					ModifiedPointPos.Y -= math.Cos(angle) * (dist_between_points * modification.TestValue)
				}
			}
		}
	}

	if triangle.Texture != nil {
		triangle.Image.SetUniforms(map[string]any{
			"PointVecOne":   []float32{float32(ModifiedPointPoses[0].X), float32(ModifiedPointPoses[0].Y)},
			"PointVecTwo":   []float32{float32(ModifiedPointPoses[1].X), float32(ModifiedPointPoses[1].Y)},
			"PointVecThree": []float32{float32(ModifiedPointPoses[2].X), float32(ModifiedPointPoses[2].Y)},

			"PointUvOne":   []float32{float32(triangle.Points[0].UvPos.X), float32(triangle.Points[0].UvPos.Y)},
			"PointUvTwo":   []float32{float32(triangle.Points[1].UvPos.X), float32(triangle.Points[1].UvPos.Y)},
			"PointUvThree": []float32{float32(triangle.Points[2].UvPos.X), float32(triangle.Points[2].UvPos.Y)},

			"Color": []float32{float32(triangle.Color.X), float32(triangle.Color.Y), float32(triangle.Color.Z)},

			"UvImgWidth":  triangle.Texture.Bounds().Dx(),
			"UvImgHeight": triangle.Texture.Bounds().Dy(),
		})
	} else {
		triangle.Image.SetUniforms(map[string]any{
			"PointVecOne":   []float32{float32(ModifiedPointPoses[0].X), float32(ModifiedPointPoses[0].Y)},
			"PointVecTwo":   []float32{float32(ModifiedPointPoses[1].X), float32(ModifiedPointPoses[1].Y)},
			"PointVecThree": []float32{float32(ModifiedPointPoses[2].X), float32(ModifiedPointPoses[2].Y)},

			"PointUvOne":   []float32{float32(triangle.Points[0].UvPos.X), float32(triangle.Points[0].UvPos.Y)},
			"PointUvTwo":   []float32{float32(triangle.Points[1].UvPos.X), float32(triangle.Points[1].UvPos.Y)},
			"PointUvThree": []float32{float32(triangle.Points[2].UvPos.X), float32(triangle.Points[2].UvPos.Y)},

			"Color": []float32{float32(triangle.Color.X), float32(triangle.Color.Y), float32(triangle.Color.Z)},

			"UvImgWidth":  0,
			"UvImgHeight": 0,
		})
		opts.Images[1] = triangle.Texture
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
	triangle.TexturePath = path
	triangle.Texture.Clear()
	triangle.Texture.DrawImage(temp_img, &ebiten.DrawImageOptions{})
}

func NewTriangle(screen_res_width, screen_res_height int) (triangle Triangle) {
	triangle.Image = textures.NewTexture("./empty.png", Triangle_Shader)
	triangle.Image.Img = ebiten.NewImage(screen_res_width, screen_res_height)
	triangle.SetTexture("./art/default_tex.png")
	triangle.SetPointsUvPos(
		utils.Vec2{X: 0, Y: 3},
		utils.Vec2{X: 1, Y: 0},
		utils.Vec2{X: 3, Y: 3},
	)

	triangle.SetPointsVectorPos(utils.Vec2{X: float64(screen_res_width) / 4 * 3, Y: float64(screen_res_height) / 4}, utils.Vec2{X: float64(screen_res_width) / 2, Y: float64(screen_res_height) / 4 * 3}, utils.Vec2{X: float64(screen_res_width) / 4, Y: float64(screen_res_height) / 4})

	return triangle
}
