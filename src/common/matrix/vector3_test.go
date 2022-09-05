package matrix

import (
	"encoding/json"
	"testing"
)

func Test_Vector3(t *testing.T) {
	t.Log()
	t.Run("newVector3", func(t *testing.T) {
		v := NewVector3(10.99, 1, 3.0)
		t.Log(v)
		vJson, err := json.Marshal(v)
		t.Log(err)
		t.Log(string(vJson))

		t.Log(v.Length())
		t.Log(v.LengthSq())
		v.Normalize()
		t.Log(v)
	})

	t.Run("Vector3_XYZAxis3", func(t *testing.T) {
		v := XYZAxis3()
		t.Log(v)
		v.Normalize()
		t.Log(v)
	})

	t.Run("ZeroVector3", func(t *testing.T) {
		v := Zero3()
		t.Log(v)
		vJson, err := json.Marshal(v)
		t.Log(err)
		t.Log(string(vJson))
		t.Log(v.Length())
		t.Log(v.LengthSq())
	})

	t.Run("Vector3Add", func(t *testing.T) {
		v1 := NewVector3(10.99, 1, 3.0)
		t.Log(v1)
		v2 := NewVector3(0.11, 10, 8.0)
		t.Log(v2)
		v1.Add(v2)
		t.Log(v1)

		v1 = NewVector3(3, 9, 3.0)
		t.Log(v1)
		v2 = NewVector3(-5, 9, 4)
		t.Log(v2)
		v1.Add(v2)
		t.Log(v1)
	})

	t.Run("Vector3Sub", func(t *testing.T) {
		v1 := NewVector3(10.99, 1, 3.0)
		t.Log(v1)
		v2 := NewVector3(0.11, 10, 8.0)
		t.Log(v2)
		v1.Sub(v2)
		t.Log(v1)
		v1.Normalize()
		t.Log(v1)

		t.Log("-----------------------------")

		v1 = NewVector3(3, 9, 3.0)
		t.Log(v1)
		v2 = NewVector3(-5, 9, 4)
		t.Log(v2)
		v1.Sub(v2)
		t.Log(v1)
		v1.Normalize()
		t.Log(v1)
	})

	t.Run("Vector3Sub", func(t *testing.T) {
		v1 := NewVector3(10.99, 1, 3.0)
		t.Log(v1)
		v2 := NewVector3(0.11, 10, 8.0)
		t.Log(v2)
		v1.Sub(v2)
		t.Log(v1)
		v1.Normalize()
		t.Log(v1)

		t.Log("-----------------------------")

		v1 = NewVector3(3, 9, 3.0)
		t.Log(v1)
		v2 = NewVector3(-5, 9, 4)
		t.Log(v2)
		v1.Sub(v2)
		t.Log(v1)
		v1.Normalize()
		t.Log(v1)
	})

	t.Run("Vector3_rotate_X", func(t *testing.T) {
		v := NewVector3(1, 1, 1)
		t.Log(v)
		// angleArr := []float64{0.0, 45.0, 90.0, 135.0, 180.0, 225.0, 270, 315.0, 360.0}
		angleArr := []float64{0, 90, 180, 270, 360}
		for _, angle := range angleArr {
			t.Log("--------------------------------------------------")
			vX := v.RotateAngle(ROTATE_AXIS_X, angle)
			t.Log("vX = ", vX)
			vY := v.RotateAngle(ROTATE_AXIS_Y, angle)
			t.Log("vY = ", vY)
			vZ := v.RotateAngle(ROTATE_AXIS_Z, angle)
			t.Log("vZ = ", vZ)
		}

	})

	t.Run("Vector3_Move", func(t *testing.T) {
		moveDir := NewVector3(1, 0.1, 1)
		t.Log("init Dir", moveDir)
		moveDir.Normalize()
		t.Log("morDir", moveDir)
		curPos := NewVector3(100, 0, 100)
		realSpeedMs := float64(1) / 1000.0
		t.Log(" realSpeedMs = ", realSpeedMs)
		for i := 0; i < 10; i++ {
			prePos := curPos.Clone()
			moved := moveDir.Clone()
			moved.Multiply(realSpeedMs * 100)
			curPos.Add(moved)
			curPos.FormatFloatFloor(6)
			t.Log("moved  = ", moved, "begin cur pos = ", prePos, "---------- end cur pos = ", curPos)
			// curPos.FormatFloatFloor(5)
			// t.Log("FormatFloatFloor", curPos)
		}

	})

	t.Run("Vector3_Rad_Angle", func(t *testing.T) {
		defDir := ZAxis3()
		t.Log("defDir = ", defDir)

		dir := Vector3{X: 0, Y: 0, Z: 1}
		t.Log("dir=", dir, ", angle=", GetAngle(defDir, dir), ",  rad =", GetRad(defDir, dir))

		dir = Vector3{X: 1, Y: 0, Z: 1}
		t.Log("dir=", dir, ", angle=", GetAngle(defDir, dir), ",  rad =", GetRad(defDir, dir))

		dir = Vector3{X: 1, Y: 0, Z: 0}
		t.Log("dir=", dir, ", angle=", GetAngle(defDir, dir), ",  rad =", GetRad(defDir, dir))

		dir = Vector3{X: 1, Y: 0, Z: -1}
		t.Log("dir=", dir, ", angle=", GetAngle(defDir, dir), ",  rad =", GetRad(defDir, dir))

		dir = Vector3{X: 0, Y: 0, Z: -1}
		t.Log("dir=", dir, ", angle=", GetAngle(defDir, dir), ",  rad =", GetRad(defDir, dir))

		dir = Vector3{X: -1, Y: 0, Z: -1}
		t.Log("dir=", dir, ", angle=", GetAngle(defDir, dir), ",  rad =", GetRad(defDir, dir))

		dir = Vector3{X: -1, Y: 0, Z: 0}
		t.Log("dir=", dir, ", angle=", GetAngle(defDir, dir), ",  rad =", GetRad(defDir, dir))

		dir = Vector3{X: -1, Y: 0, Z: 1}
		t.Log("dir=", dir, ", angle=", GetAngle(defDir, dir), ",  rad =", GetRad(defDir, dir))

		dir = Vector3{X: 0, Y: 0, Z: 1}
		t.Log("dir=", dir, ", angle=", GetAngle(defDir, dir), ",  rad =", GetRad(defDir, dir))

	})

}
