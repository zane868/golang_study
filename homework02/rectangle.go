package main

// Rectangle 表示矩形，包含长和宽两个维度。
type Rectangle struct {
	length int
	width  int
}

// Area 返回矩形面积。
func (r *Rectangle) Area() int {
	return r.length * r.width
}

// Perimeter 返回矩形周长。
func (r *Rectangle) Perimeter() int {
	return 2 * (r.length + r.width)
}
