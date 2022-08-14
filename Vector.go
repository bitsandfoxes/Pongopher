package main

import "math"

type Vector struct {
	x float64
	y float64
}

func (v *Vector) Add(v1 Vector) Vector {
	return Vector{v.x + v1.x, v.y + v1.y}
}

func (v *Vector) Subtract(v1 Vector) Vector {
	return Vector{v.x - v1.x, v.y - v1.y}
}

func (v *Vector) Multiply(m float64) Vector {
	return Vector{v.x * m, v.y * m}
}

//func Subtract(v1, v2 Vector) Vector {
//	return Vector{v1.x - v1.x, v1.y - v1.y}
//}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v *Vector) Normalized() Vector {
	length := math.Abs(v.Length())
	return Vector{v.x / length, v.y / length}
}
