package main

// 创建游戏对象
type GameObject struct {
	width  int
	height int
	x      float64
	y      float64
}

func (gameObj *GameObject) Width() int {
	return gameObj.width
}

func (gameObj *GameObject) Height() int {
	return gameObj.height
}

func (gameObj *GameObject) X() float64 {
	return gameObj.x
}

func (gameObj *GameObject) Y() float64 {
	return gameObj.y
}

// 创建接口
type Entity interface {
	Width() int
	Height() int
	X() float64
	Y() float64
}
