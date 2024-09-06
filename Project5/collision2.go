package main

import ()

// CheckCollision 检查飞机和外星人之间是否有碰撞  -- 代码基本上是一致的，直接复刻基本上就可以了
func CheckCollision2(ship *Ship, alien *Alien) bool {
	alienTop, alienLeft := alien.y, alien.x
	alienBottom, alienRight := alien.y+float64(alien.height), alien.x+float64(alien.width)
	// 左上角
	x, y := ship.x, ship.y
	if y > alienTop && y < alienBottom && x > alienLeft && x < alienRight {
		return true
	}

	// 右上角
	x, y = ship.x+float64(ship.width), ship.y
	if y > alienTop && y < alienBottom && x > alienLeft && x < alienRight {
		return true
	}

	// 左下角
	x, y = ship.x, ship.y+float64(ship.height)
	if y > alienTop && y < alienBottom && x > alienLeft && x < alienRight {
		return true
	}

	// 右下角
	x, y = ship.x+float64(ship.width), ship.y+float64(ship.height)
	if y > alienTop && y < alienBottom && x > alienLeft && x < alienRight {
		return true
	}

	return false
}
