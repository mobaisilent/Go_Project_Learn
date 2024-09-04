package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Ship struct {
	image  *ebiten.Image
	width  int
	height int
	x      float64 // x坐标
	y      float64 // y坐标
}

// 在game.go那里实现传入这两个参数
func NewShip(screenWidth int, screenHeight int) *Ship {
	// 用ebiten自带的就能解析png图片
	img, _, err := ebitenutil.NewImageFromFile("resource/ship.png")
	if err != nil {
		log.Fatal(err)
	}

	width, height := img.Size()
	fmt.Println("width:", width, "height:", height)
	// 能够获取到尺寸 60 48 呢 那么应该不是飞船对象问题导致的图片加载失败
	ship := &Ship{
		image:  img,
		width:  width,
		height: height,
		x:      float64(screenWidth-width) / 2,
		y:      float64(screenHeight - height),
	}

	return ship
}

func (ship *Ship) Draw(screen *ebiten.Image, cfg *Config) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(cfg.ScreenWidth-ship.width)/2, float64(cfg.ScreenHeight-ship.height))
	// 这里是做了减法运算获取到飞船的左上角位置，是的飞船位于中心位置
	screen.DrawImage(ship.image, op)
	// fmt.Println("draw ship") 是一直在绘制中
}
