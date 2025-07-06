package anamorph

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

type Point struct {
	X float64
	Y float64
}

func (p Point) Equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

type Anamorpher struct {
	Angle  float64
	Radius float64
	pos    Point
	Img    image.Image
	Mod    ImageMod
}

func New(img image.Image, mod ImageMod, angle, radius float64) *Anamorpher {
	return &Anamorpher{
		Angle:  angle,
		Radius: radius,
		Img:    img,
		Mod:    mod,
		pos:    Point{0, 0},
	}
}

type ImageMod interface {
	image.Image
	Set(x, y int, c color.Color)
}

func SetAll(i ImageMod, c color.Color) {
	for y := i.Bounds().Min.Y; y < i.Bounds().Max.Y; y++ {
		for x := i.Bounds().Min.X; x < i.Bounds().Max.X; x++ {
			i.Set(x, y, c)
		}
	}
}

var ErrImageOutOfBounds = fmt.Errorf("image out of bounds")

func (a *Anamorpher) Anamorph() error {
	if a.pos.Equals(Point{0, 0}) {
		SetAll(a.Mod, color.White)
	}

	imgratio := a.Radius / float64(a.Img.Bounds().Size().X)

	for y := a.Img.Bounds().Min.Y; y < a.Img.Bounds().Max.Y; y++ {
		for x := a.Img.Bounds().Min.X; x < a.Img.Bounds().Max.X; x++ {
			a.pos = Point{float64(x), float64(y)}
			p := GetNewPoint(Point{float64(x), float64(y)}, a.Angle, float64(a.Img.Bounds().Size().X)/2)
			nX := float64(int((p.X * imgratio) + (float64(a.Mod.Bounds().Size().X) / 2)))
			nY := float64(int((p.Y * imgratio)))
			if nX < float64(a.Mod.Bounds().Min.X) || nX > float64(a.Mod.Bounds().Max.X) || nY < float64(a.Mod.Bounds().Min.Y) || nY > float64(a.Mod.Bounds().Max.Y) {
				return ErrImageOutOfBounds
			}
			p.X = nX
			p.Y = nY

			a.Mod.Set(int(math.Round(p.X)), int(math.Round(p.Y)), a.Img.At(x, (a.Img.Bounds().Max.Y-1)-y))

		}
	}

	return nil
}

func (a *Anamorpher) MaximumRequiredBounds() image.Rectangle {
	ratio := a.Radius / float64(a.Img.Bounds().Size().X)
	yPoint := GetNewPoint(Point{math.Floor(float64(a.Img.Bounds().Max.X+a.Img.Bounds().Min.X) / 2), float64(a.Img.Bounds().Max.Y)}, a.Angle, float64(a.Img.Bounds().Size().X)/2)
	xPoint := GetNewPoint(Point{float64(a.Img.Bounds().Max.X), float64(a.Img.Bounds().Max.Y)}, a.Angle, float64(a.Img.Bounds().Size().X)/2)
	return image.Rect(0, 0, int((2*xPoint.X)*ratio), int(yPoint.Y*ratio))
}

func GetNewPoint(p Point, angle float64, radius float64) Point {
	s := math.Sin(angle)
	t := math.Tan(angle)
	d := s * (t + 1) * p.Y
	l := (p.X - radius) / radius
	a := math.Asin(l)
	w := Point{
		X: radius * math.Sin(a),
		Y: radius * math.Cos(a),
	}

	return Point{
		X: w.X + d*math.Sin(a),
		Y: w.Y + d*math.Cos(a),
	}
}
