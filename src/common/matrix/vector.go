package matrix

import (
	"math"
)

type Vector2D struct {
	X float32
	Y float32
}

var NormRightX *Vector2D = &Vector2D{X: 1.0}
var NormLeftX *Vector2D = &Vector2D{X: -1.0}
var NormRightY *Vector2D = &Vector2D{Y: 1.0}
var NormLeftY *Vector2D = &Vector2D{Y: -1.0}

func (this *Vector2D) Clear() {
	this.X = 0
	this.Y = 0
}

func (this *Vector2D) Copy() *Vector2D {
	newB := new(Vector2D)
	*newB = (*this)
	return newB
}

func (this *Vector2D) IsZero() bool {
	return 0 == this.X && 0 == this.Y
}

func (this *Vector2D) ToAngle() float64 {
	y := &Vector2D{X: 0, Y: 1}
	return this.Angle(y)
}

func (this *Vector2D) ToRadian() float64 {
	y := &Vector2D{X: 0, Y: 1}
	return this.Radian(y)
}

func (this *Vector2D) ContainInRectRange(v *Vector2D, r float32) bool {
	if v.X >= this.X-r && v.X <= this.X+r && v.Y >= this.Y-r && v.Y <= this.Y+r {
		return true
	}
	return false
}

func (this *Vector2D) ContainInRange(v *Vector2D, r float32) bool {
	if this.DistanceSq(v) <= float64(r)*float64(r) {
		return true
	}
	return false
}

func (this *Vector2D) LengthSq() float64 {
	return float64(this.X)*float64(this.X) + float64(this.Y)*float64(this.Y)
}

func (this *Vector2D) Length() float32 {
	return float32(math.Sqrt(this.LengthSq()))
}

func (this *Vector2D) DistanceSq(v *Vector2D) float64 {
	deltX := float64(this.X - v.X)
	deltY := float64(this.Y - v.Y)
	return deltX*deltX + deltY*deltY
}

func (this *Vector2D) Distance(v *Vector2D) float64 {
	return math.Sqrt(this.DistanceSq(v))
}

func (this *Vector2D) Equal(v *Vector2D) bool {
	return 0 == math.Round(float64(this.X-v.X)) && 0 == math.Round(float64(this.Y-v.Y))
	// return math.Abs(this.X - v.X) <= ERROR_MARGIN && math.Abs(this.Y - v.Y) <= ERROR_MARGIN
}

func (this *Vector2D) Orth() *Vector2D {
	return &Vector2D{
		X: -this.Y,
		Y: this.X,
	}
}

func (this *Vector2D) Reverse() *Vector2D {
	return &Vector2D{
		X: -this.X,
		Y: -this.Y,
	}
}

func (this *Vector2D) Add(v *Vector2D) *Vector2D {
	if v == nil {
		return &Vector2D{
			X: this.X,
			Y: this.Y,
		}
	}
	return &Vector2D{
		X: this.X + v.X,
		Y: this.Y + v.Y,
	}
}

func (this *Vector2D) Sub(v *Vector2D) *Vector2D {
	return &Vector2D{
		X: this.X - v.X,
		Y: this.Y - v.Y,
	}
}

func (this *Vector2D) Mul(v float32) *Vector2D {
	return &Vector2D{
		X: this.X * v,
		Y: this.Y * v,
	}
}

func (this *Vector2D) Div(v float32) *Vector2D {
	return &Vector2D{
		X: this.X / v,
		Y: this.Y / v,
	}
}

func (this *Vector2D) Angle(v *Vector2D) float64 {
	return this.Radian(v) * 180 / math.Pi
}

func (this *Vector2D) Radian(v *Vector2D) float64 {
	sin := float64(this.X)*float64(v.Y) - float64(v.X)*float64(this.Y)
	cos := float64(this.X)*float64(v.X) + float64(this.Y)*float64(v.Y)
	return math.Atan2(sin, cos)
}

func (this *Vector2D) Rotate(alpha float64) *Vector2D {
	return &Vector2D{
		X: float32(float64(this.X)*math.Cos(alpha) - float64(this.Y)*math.Sin(alpha)),
		Y: float32(float64(this.X)*math.Sin(alpha) + float64(this.Y)*math.Cos(alpha)),
	}
}

func (this *Vector2D) RotateAngle(angle float64) *Vector2D {
	return this.Rotate(float64(math.Pi) * float64(angle) / 180)
}

func (this *Vector2D) Norm() *Vector2D {
	len := this.Length()
	return &Vector2D{
		X: this.X / len,
		Y: this.Y / len,
	}
}

// copy from r2.point
func (this *Vector2D) Cross(op *Vector2D) float64 {
	return float64(this.X)*float64(op.Y) - float64(this.Y)*float64(op.X)
}
func (this *Vector2D) Dot(op *Vector2D) float64 {
	return float64(this.X)*float64(op.X) + float64(this.Y)*float64(op.Y)
}
