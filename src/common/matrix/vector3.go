package matrix

import (
	"fmt"
	"math"
	"strconv"
)

type RotateAxis uint

const (
	ROTATE_AXIS_X RotateAxis = 0
	ROTATE_AXIS_Y RotateAxis = 1
	ROTATE_AXIS_Z RotateAxis = 2
)

// 三维向量：(x,y,z)
type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// 返回：新向量
func NewVector3(x, y, z float64) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

func (this *Vector3) IsZero() bool {
	return this.X == 0.0 && this.Y == 0.0 && this.Z == 0.0
}

func (this *Vector3) Equal(v Vector3) bool {
	return this.X == v.X && this.Y == v.Y && this.Z == v.Z
}

// 三维向量：设值
func (this *Vector3) Set(x, y, z float64) {
	this.X = x
	this.Y = y
	this.Z = z
}

// 三维向量：拷贝
func (this *Vector3) Clone() Vector3 {
	return NewVector3(this.X, this.Y, this.Z)
}

// 三维向量：长度
func (this *Vector3) Length() float64 {
	return math.Sqrt(this.X*this.X + this.Y*this.Y + this.Z*this.Z)
}

// 三维向量：长度平方
func (this *Vector3) LengthSq() float64 {
	return this.X*this.X + this.Y*this.Y + this.Z*this.Z
}

// 三维向量：加上 this = this + v
func (this *Vector3) Add(v Vector3) {
	this.X += v.X
	this.Y += v.Y
	this.Z += v.Z
}

// 三维向量：减去  this = this - v
func (this *Vector3) Sub(v Vector3) {
	this.X -= v.X
	this.Y -= v.Y
	this.Z -= v.Z
}

// 三维向量：数乘
func (this *Vector3) Multiply(scalar float64) {
	this.X *= scalar
	this.Y *= scalar
	this.Z *= scalar
}

func (this *Vector3) Divide(scalar float64) {
	if scalar == 0 {
		panic("分母不能为零！")
	}
	this.Multiply(1 / scalar)
}

// 三维向量：单位化
func (this *Vector3) Normalize() {
	this.Divide(this.Length())
}

// 三维向量：点积
func (this *Vector3) Dot(v Vector3) float64 {
	return this.X*v.X + this.Y*v.Y + this.Z*v.Z
}

// 三维向量：叉积
func (this *Vector3) Cross(v Vector3) {
	x, y, z := this.X, this.Y, this.Z
	this.X = y*v.Z - z*v.Y
	this.Y = z*v.X - x*v.Z
	this.Z = x*v.Y - y*v.X
}

// 沿 axis 轴顺时针旋转 alpha 弧度
func (this *Vector3) Rotate(axis RotateAxis, alpha float64) (v Vector3) {
	switch axis {
	case ROTATE_AXIS_X:
		v = Vector3{
			X: this.X,
			Y: this.Y*math.Cos(alpha) - this.Z*math.Sin(alpha),
			Z: -this.Y*math.Sin(alpha) + this.Z*math.Cos(alpha),
		}
	case ROTATE_AXIS_Y:
		v = Vector3{
			X: this.X*math.Cos(alpha) + this.Z*math.Sin(alpha),
			Y: this.Y,
			Z: -this.X*math.Sin(alpha) + this.Z*math.Cos(alpha),
		}
	default:
		v = Vector3{
			X: this.Y*math.Cos(alpha) + this.X*math.Sin(alpha),
			Y: -this.Y*math.Sin(alpha) + this.X*math.Cos(alpha),
			Z: this.Z,
		}
	}

	return v
}

// 沿 axis 轴顺时针旋转 angle 欧拉角度
func (this *Vector3) RotateAngle(axis RotateAxis, angle float64) Vector3 {
	return this.Rotate(axis, angle*math.Pi/180.0)
}

// 格式化数据（取小数点后decimal位有效四舍五入后的数据）
func (this *Vector3) FormatFloatFloor(decimal int) {
	formatFloat := func(num float64, decimal int) (float64, error) {
		if decimal < 1 {
			decimal = 1
		}
		formatStr := fmt.Sprintf("%."+fmt.Sprint(decimal)+"f", num)
		return strconv.ParseFloat(formatStr, 64)
	}

	if decimal < 1 {
		decimal = 1
	}

	this.X, _ = formatFloat(this.X, decimal)
	this.Y, _ = formatFloat(this.Y, decimal)
	this.Z, _ = formatFloat(this.Z, decimal)
}

// 返回：零向量(0,0,0)
func Zero3() Vector3 {
	return Vector3{X: 0, Y: 0, Z: 0}
}

// X 轴 单位向量
func XAxis3() Vector3 {
	return Vector3{X: 1, Y: 0, Z: 0}
}

// Y 轴 单位向量
func YAxis3() Vector3 {
	return Vector3{X: 0, Y: 1, Z: 0}
}

// Z 轴 单位向量
func ZAxis3() Vector3 {
	return Vector3{X: 0, Y: 0, Z: 1}
}

func XYAxis3() Vector3 {
	return Vector3{X: 1, Y: 1, Z: 0}
}
func XZAxis3() Vector3 {
	return Vector3{X: 1, Y: 0, Z: 1}
}
func YZAxis3() Vector3 {
	return Vector3{X: 0, Y: 1, Z: 1}
}
func XYZAxis3() Vector3 {
	return Vector3{X: 1, Y: 1, Z: 1}
}

// 返回：a + b 向量
func Add3(a, b Vector3) Vector3 {
	return Vector3{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

// 返回：a - b 向量
func Sub3(a, b Vector3) Vector3 {
	return Vector3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

// 返回：a X b 向量 (X 叉乘)
func Cross3(a, b Vector3) Vector3 {
	return Vector3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func AddArray3(vs []Vector3, dv Vector3) []Vector3 {
	for i, _ := range vs {
		vs[i].Add(dv)
	}
	return vs
}

func Multiply(v Vector3, scalars float64) Vector3 {
	vector := v.Clone()
	vector.Multiply(scalars)
	return vector
}

func Multiply3(v Vector3, scalars []float64) []Vector3 {
	vs := []Vector3{}
	for _, value := range scalars {
		vector := v.Clone()
		vector.Multiply(value)
		vs = append(vs, vector)
	}
	return vs
}

// 返回：单位化向量
func Normalize3(a Vector3) Vector3 {
	b := a.Clone()
	b.Normalize()
	return b
}

// 求两点间距离平方
func DistanceSq(a, b Vector3) float64 {
	return math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2) + math.Pow(a.Z-b.Z, 2)
}

// 求两点间距离
func Distance(a, b Vector3) float64 {
	return math.Sqrt(DistanceSq(a, b))
}

// 求两向量弧度
func GetRad(v1, v2 Vector3) (angel float64) {
	a := v1.Dot(v2)
	b := math.Sqrt(math.Pow(v1.X, 2)+math.Pow(v1.Y, 2)+math.Pow(v1.Z, 2)) *
		math.Sqrt(math.Pow(v2.X, 2)+math.Pow(v2.Y, 2)+math.Pow(v2.Z, 2))

	angel = math.Acos(a / b)
	return
}

// 求两向量夹角
func GetAngle(v1, v2 Vector3) (angel float64) {
	return GetRad(v1, v2) * 180 / math.Pi
}
