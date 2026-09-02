package main

// Shape 表示几何图形接口。
type Shape interface {
	Area() int
	Perimeter() int
}
