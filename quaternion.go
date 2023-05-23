package main

import (
	"math"
)

func multiplyQuaternion(q1 Quaternion, q0 Quaternion) Quaternion {
	x0, y0, z0, w0 := q0.X, q0.Y, q0.Z, q0.W
	x1, y1, z1, w1 := q1.X, q1.Y, q1.Z, q1.W
	return Quaternion{
		X: x1*w0 + y1*z0 - z1*y0 + w1*x0,
		Y: -x1*z0 + y1*w0 + z1*x0 + w1*y0,
		Z: x1*y0 - y1*x0 + z1*w0 + w1*z0,
		W: -x1*x0 - y1*y0 - z1*z0 + w1*w0,
	}
}

func threeAxisRotote(r11 float32, r12 float32, r21 float32, r31 float32, r32 float32) [3]float32 {
	return [3]float32{
		float32(math.Atan2(float64(r11), float64(r12))),
		float32(math.Asin(float64(r21))),
		float32(math.Atan2(float64(r31), float64(r32))),
	}
}

func quaternionToEuler(q Quaternion, order string) Euler {
	switch order {
	case "ZYX":
		ary := threeAxisRotote(
			2*(q.X*q.Y+q.W*q.Z),
			q.W*q.W+q.X*q.X-q.Y*q.Y-q.Z*q.Z,
			-2*(q.X*q.Z-q.W*q.Y),
			2*(q.Y*q.Z+q.W*q.X),
			q.W*q.W-q.X*q.X-q.Y*q.Y+q.Z*q.Z,
		)
		return Euler{ary[2], ary[1], ary[0]}
	case "ZXY":
		ary := threeAxisRotote(
			-2*(q.X*q.Y-q.W*q.Z),
			q.W*q.W-q.X*q.X+q.Y*q.Y-q.Z*q.Z,
			2*(q.Y*q.Z+q.W*q.X),
			-2*(q.X*q.Z-q.W*q.Y),
			q.W*q.W-q.X*q.X-q.Y*q.Y+q.Z*q.Z,
		)
		return Euler{ary[1], ary[2], ary[0]}
	case "YXZ":
		ary := threeAxisRotote(
			2*(q.X*q.Z+q.W*q.Y),
			q.W*q.W-q.X*q.X-q.Y*q.Y+q.Z*q.Z,
			-2*(q.Y*q.Z-q.W*q.X),
			2*(q.X*q.Y+q.W*q.Z),
			q.W*q.W-q.X*q.X+q.Y*q.Y-q.Z*q.Z,
		)
		return Euler{ary[1], ary[0], ary[2]}
	case "YZX":
		ary := threeAxisRotote(
			-2*(q.X*q.Z-q.W*q.Y),
			q.W*q.W+q.X*q.X-q.Y*q.Y-q.Z*q.Z,
			2*(q.X*q.Y+q.W*q.Z),
			-2*(q.Y*q.Z-q.W*q.X),
			q.W*q.W-q.X*q.X+q.Y*q.Y-q.Z*q.Z,
		)
		return Euler{ary[2], ary[0], ary[1]}
	case "XYZ":
		ary := threeAxisRotote(
			-2*(q.Y*q.Z-q.W*q.X),
			q.W*q.W-q.X*q.X-q.Y*q.Y+q.Z*q.Z,
			2*(q.X*q.Z+q.W*q.Y),
			-2*(q.X*q.Y-q.W*q.Z),
			q.W*q.W+q.X*q.X-q.Y*q.Y-q.Z*q.Z,
		)
		return Euler{ary[0], ary[1], ary[2]}
	case "XZY":
		ary := threeAxisRotote(
			2*(q.Y*q.Z+q.W*q.X),
			q.W*q.W-q.X*q.X+q.Y*q.Y-q.Z*q.Z,
			-2*(q.X*q.Y-q.W*q.Z),
			2*(q.X*q.Z+q.W*q.Y),
			q.W*q.W+q.X*q.X-q.Y*q.Y-q.Z*q.Z,
		)
		return Euler{ary[0], ary[2], ary[1]}
	default:
		return Euler{0., 0., 0.}
	}
}

func inverseQuaternion(q Quaternion) Quaternion {
	return Quaternion{q.W, -q.X, -q.Y, -q.Z}
}
