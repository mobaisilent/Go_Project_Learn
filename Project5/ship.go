package main

import (
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
	// fmt.Println("width:", width, "height:", height)
	// 能够获取到尺寸 60 48 呢 那么应该不是飞船对象问题导致的图片加载失败
	ship := &Ship{
		image:  img,
		width:  width,
		height: height,
		x:      float64(screenWidth-width) / 2,
		y:      float64(screenHeight - height),
	}
	// 这里传值的时候就已经界定好了x和y，所以下面的Draw函数直接填写就行

	return ship
}

func (ship *Ship) Draw(screen *ebiten.Image, cfg *Config) {
	op := &ebiten.DrawImageOptions{}
	// ship.x = float64(cfg.ScreenWidth-ship.width) / 2
	ship.y = float64(cfg.ScreenHeight - ship.height)
	op.GeoM.Translate(ship.x, ship.y)
	// 这里是做了减法运算获取到飞船的左上角位置，是的飞船位于中心位置
	screen.DrawImage(ship.image, op)
	// fmt.Println("draw ship") 是一直在绘制中
}

// 未实现移动时因为这里Draw的op.GeoM.Translate是固定的，所以飞船不会移动
// emm
