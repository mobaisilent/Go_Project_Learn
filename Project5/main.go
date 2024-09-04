package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// 更新游戏内容显示

func (g *Game) Update() error {
	g.input.Update(g.ship, g.cfg)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.cfg.BgColor) // 绘制背景
	g.ship.Draw(screen, g.cfg) // 绘制飞船
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

// 这里不用除以2 直接返回config.json中的尺寸大小即可， /2 之后飞船绘制不出来

// 关于Game结构体写三个接口 实现游戏的初始化、更新、渲染

func main() {

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
