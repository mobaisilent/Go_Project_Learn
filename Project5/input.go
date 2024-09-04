package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Input struct {
	msg string
}

// NewGame()函数用于初始化游戏：为Game结构体的信息赋值

func (i *Input) Update(ship *Ship) {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		fmt.Println("left")
		ship.x -= 1
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		fmt.Println("right")
		ship.x += 1
	}
}
// 能够检测到按下的按键，但是飞船的位置没有发生变化？

