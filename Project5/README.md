> 该项目主要是利用已知go语言包创建一个小游戏
>
> 官方网站：
>
> https://ebitengine.org/
>
> 仓库地址 ：
> https://github.com/hajimehoshi/ebiten/

# Part1

## 安装

ebitengine 要求Go版本 >= 1.15。使用go module下载这个包：

```shell
go get -u github.com/hajimehoshi/ebiten/v2
```

## 初始化界面

main.go

```go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

// 关于Game结构体写三个接口 实现游戏的初始化、更新、渲染

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("外星人入侵")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

```

![image-20240903112758823](images/image-20240903112758823.png)

首先，ebiten引擎运行时要求传入一个游戏对象，该对象的必须实现`ebiten.Game`这个接口：

```go
// Game defines necessary functions for a game.
type Game interface {
  Update() error
  Draw(screen *Image)
  Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}
```

`ebiten.Game`接口定义了ebiten游戏需要的3个方法：`Update`,`Draw`和`Layout`。

- `Update`：每个tick都会被调用。tick是引擎更新的一个时间单位，默认为1/60s。tick的倒数我们一般称为帧，即游戏的更新频率。默认ebiten游戏是60帧，即每秒更新60次。该方法主要用来更新游戏的逻辑状态，例如子弹位置更新。上面的例子中，游戏对象没有任何状态，故`Update`方法为空。注意到`Update`方法的返回值为`error`类型，当`Update`方法返回一个非空的`error`值时，游戏停止。在上面的例子中，我们一直返回nil，故只有关闭窗口时游戏才停止。
- `Draw`：每帧（frame）调用。帧是渲染使用的一个时间单位，依赖显示器的刷新率。如果显示器的刷新率为60Hz，`Draw`将会每秒被调用60次。`Draw`接受一个类型为`*ebiten.Image`的screen对象。ebiten引擎每帧会渲染这个screen。在上面的例子中，我们调用`ebitenutil.DebugPrint`函数在screen上渲染一条调试信息。由于调用`Draw`方法前，screen会被重置，故`DebugPrint`每次都需要调用。
- `Layout`：该方法接收游戏窗口的尺寸作为参数，返回游戏的逻辑屏幕大小。我们实际上计算坐标是对应这个逻辑屏幕的，`Draw`将逻辑屏幕渲染到实际窗口上。这个时候可能会出现伸缩。在上面的例子中游戏窗口大小为(640, 480)，`Layout`返回的逻辑大小为(320, 240)，所以显示会放大1倍。

> Layout是设置逻辑分辨率，而main函数中设置的是实际分辨率
>
> 不想区分这里或者因为尺寸问题造成元素不显示就直接传递数值就好，根据效果对应修改就行了
>
> 比如后面的飞船那里：如果不删去 /2 的话飞船不显示

## 处理输入

没有交互的游戏不是真的游戏！下面我们来监听键盘的输入，当前只处理3个键：左方向←，右方向→和空格。

ebiten提供函数`IsKeyPressed`来判断某个键是否按下，同时内置了100多个键的常量定义，见源码keys.go文件。`ebiten.KeyLeft`表示左方向键，`ebiten.KeyRight`表示右方向键，`ebiten.KeySpace`表示空格。

为了代码清晰，我们定义一个`Input`结构来处理输入：

```go
type Input struct {
  msg string
}

func (i *Input) Update() {
  if ebiten.IsKeyPressed(ebiten.KeyLeft) {
    fmt.Println("←←←←←←←←←←←←←←←←←←←←←←←")
    i.msg = "left pressed"
  } else if ebiten.IsKeyPressed(ebiten.KeyRight) {
    fmt.Println("→→→→→→→→→→→→→→→→→→→→→→→")
    i.msg = "right pressed"
  } else if ebiten.IsKeyPressed(ebiten.KeySpace) {
    fmt.Println("-----------------------")
    i.msg = "space pressed"
  }
}
```

Game结构中添加一个`Input`类型的字段，并且为了方便新增`NewGame`方法用于创建Game对象：

```go
type Game struct {
  input *Input
}

func NewGame() *Game {
  return &Game{
    input: &Input{msg: "Hello, World!"},
  }
}
```

Game结构的`Update`方法中，我们需要调用`Input`的`Update`方法触发按键的判断：

```go
func (g *Game) Update() error {
  g.input.Update()
  return nil
}
```

Game的`Draw`方法中将显示`Input`的`msg`字段：

```go
func (g *Game) Draw(screen *ebiten.Image) {
  ebitenutil.DebugPrint(screen, g.input.msg)
}
```

将main函数中创建Game对象的方式修改如下：

```go
game := NewGame()

if err := ebiten.RunGame(game); err != nil {
  log.Fatal(err)
}
```

修改的完整main.go代码如下；

```go
package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Input struct {
	msg string
}

type Game struct {
	input *Input
}

func NewGame() *Game {
	return &Game{
		input: &Input{msg: "Hello, World!"},
	}
}

// NewGame()函数用于初始化游戏：为Game结构体的信息赋值

func (i *Input) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		fmt.Println("←←You Pressed Left ←←")
		i.msg = "left pressed"
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		fmt.Println("→→You Pressed Right→→")
		i.msg = "right pressed"
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) {
		fmt.Println("--You Pressed Space--")
		i.msg = "space pressed"
	}
}

// 更新游戏内容显示

func (g *Game) Update() error {
	g.input.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, g.input.msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

// 关于Game结构体写三个接口 实现游戏的初始化、更新、渲染

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("外星人入侵")

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

```

> 核心还是 Updata Draw 和 Layout三个函数，不过嵌套了一个内部结构体实现显示内容的更改

结果如下；

![image-20240903174434458](images/image-20240903174434458.png)

如果你一直按着那么下面的终端就会重复打印信息

## 设置背景

代码如下：（修改Draw函数实现更改背景）

```go
func (g *Game) Draw(screen *ebiten.Image) {
  screen.Fill(color.RGBA{R: 200, G: 200, B: 200, A: 255})
  ebitenutil.DebugPrint(screen, g.input.msg)
}
```

结果如下：

![image-20240903174838967](images/image-20240903174838967.png)

>核心就是 screen.Fill(color.RGBA{R: 200, G: 200, B: 200, A: 255})
>
>通过这行代码进行颜色选择就好，RGBA的具体含义建议google

## 第一次重构

目前为止，我们的实现了显示窗口和处理输入的功能。我们先分析一下目前的程序有哪些问题：

- 所有逻辑都堆在一个文件中，修改不便
- 逻辑中直接出现字面值，例如640/480，字符串"外星人入侵"等，每次修改都需要重新编译程序

在继续之前，我们先对代码组织结构做一次重构，这能让我们走得更远。

为了清晰，方便管理，我们逻辑拆分到4个文件中：

- game.go：编写Game对象，并实现相关方法，同时负责协调其他各个模块
- input.go：输入相关的逻辑
- config.go：专门负责配置相关的逻辑
- main.go：main函数，负责创建Game对象，运行游戏

为了程序的灵活修改，我们将程序中的可变项都作为配置存放在文件中，程序启动时自动读取这个配置文件。我选择json作为配置文件的格式：(json格式用来写配置信息是非常推荐的呢)

config.json

```json
{
  "screenWidth": 640,
  "screenHeight": 480,
  "title": "外星人入侵",
  "bgColor": {
    "r": 230,
    "g": 230,
    "b": 230,
    "a": 255
  }
}
```

然后定义配置的结构和加载配置的函数：

config.go

```golang
package main

import (
	"encoding/json"
	"image/color"
	"log"
	"os"
)

type Config struct {
	ScreenWidth  int        `json:"screenWidth"`
	ScreenHeight int        `json:"screenHeight"`
	Title        string     `json:"title"`
	BgColor      color.RGBA `json:"bgColor"`
}

func loadConfig() *Config {
	f, err := os.Open("./config.json")
	if err != nil {
		log.Fatalf("os.Open failed: %v\n", err)
	}

	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg) // 从文件中读取json数据并解码到cfg中
	if err != nil {
		log.Fatalf("json.Decode failed: %v\n", err)
	}

	return &cfg
}

```

将游戏核心逻辑移到game.go文件中，定义游戏对象结构和创建游戏对象的方法：

game.go

```golang
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	input *Input
	cfg   *Config
}

func NewGame() *Game {
	cfg := loadConfig()
	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetWindowTitle(cfg.Title)

	return &Game{
		input: &Input{
			msg: "hello world",
		},
		cfg: cfg,
	}
}

```

先从配置文件中加载配置，然后根据配置设置游戏窗口大小和标题。拆分之后，`Draw`和`Layout`方法实现如下：

main.go

```golang
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// 更新游戏内容显示

func (g *Game) Update() error {
	g.input.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.cfg.BgColor)
	ebitenutil.DebugPrint(screen, g.input.msg)
  // 这里的绘制顺序不能搞错
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.cfg.ScreenWidth / 2, g.cfg.ScreenHeight / 2
}

// 关于Game结构体写三个接口 实现游戏的初始化、更新、渲染

func main() {

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

```

input.go

```go
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
)

type Input struct {
	msg string
}

// NewGame()函数用于初始化游戏：为Game结构体的信息赋值

func (i *Input) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		fmt.Println("←←You Pressed Left ←←")
		i.msg = "left pressed"
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		fmt.Println("→→You Pressed Right→→")
		i.msg = "right pressed"
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) {
		fmt.Println("--You Pressed Space--")
		i.msg = "space pressed"
	}
}

```



第一次重构到此完成，现在来看一下文件结构，是否更清晰了呢？

```fallback
├── config.go
├── config.json
├── game.go
├── input.go
└── main.go
```

注意，因为拆分成了多个文件，所以运行程序不能再使用`go run main.go`命令了，需要改为`go run .`

![image-20240903211617123](images/image-20240903211617123.png)

> 效果与之前的整体性代码保持一致

## 显示图片

接下来我们尝试在屏幕底部中心位置显示一张飞船的图片：

![img](images/ship.png)

ebitengine引擎提供了`ebitenutil.NewImageFromFile`函数，传入图片路径即可加载该图片，so easy。为了很好的管理游戏中的各个实体，我们给每个实体都定义一个结构。先定义飞船结构：

ship.go

```go
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
}

func NewShip() *Ship {
	// 用ebiten自带的就能解析png图片
    img, _, err := ebitenutil.NewImageFromFile("../resource/ship.png")
    if err != nil {
        log.Fatal(err)
    }

  	// 通过对象获取width和height然后直接传递数值就行
    width, height := img.Size()
    ship := &Ship{
        image:  img,
        width:  width,
        height: height,
    }

    return ship
}
```

Go标准库提供了三种格式的解码包，`image/png`，`image/jpeg`，`image/gif`。也就是说标准库中没有bmp格式的解码包，所幸golang.org/x仓库没有让我们失望，golang.org/x/image/bmp提供了解析bmp格式图片的功能。我们这里不需要显式的使用对应的图片库，故使用`import _`这种方式，让`init`函数产生副作用。

然后在游戏对象中添加飞船类型的字段：

```golang
func NewGame() *Game {
  // 相同的代码省略...
  return &Game {
    input:   &Input{},
    ship:  NewShip(),
    cfg:  cfg,
  }
}
```

为了将飞船显示在屏幕底部中央位置，我们需要计算坐标。ebitengine采用如下所示的二维坐标系：

![img](images/ebiten8.png)

x轴向右，y轴向下，左上角为原点。我们需要计算飞船左上角的位置。由上图很容易计算出：

```fallback
x=(W1-W2)/2
y=H1-H2
```

为了在屏幕上显示飞船图片，我们需要调用`*ebiten.Image`的`DrawImage`方法，该方法的第二个参数可以用于指定坐标相对于原点的偏移：

```golang
func (g *Game) Draw(screen *ebiten.Image) {
  screen.Fill(g.cfg.BgColor)
  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(float64(g.cfg.ScreenWidth-g.ship.width)/2, float64(g.cfg.ScreenHeight-g.ship.height))
  screen.DrawImage(g.ship.image, op)
}
```

我们给`Ship`类型增加一个绘制自身的方法，传入屏幕对象screen和配置，让代码更好维护：

```golang
func (ship *Ship) Draw(screen *ebiten.Image, cfg *Config) {
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(float64(cfg.ScreenWidth-ship.width)/2, float64(cfg.ScreenHeight-ship.height))
    screen.DrawImage(ship.image, op)
}
```

main.go函数中的Draw就可以简化为：

```go
func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(g.cfg.BgColor)
    g.ship.Draw(screen, g.cfg)
}
```

文件结构如下；

> 相比于上一步多了`ship.go`然后修改下game结构体和main函数即可

```bash
mobai@mobaideAir Project5 % tree .
.
├── backup
│   └── main.go
├── config.go
├── config.json
├── game.go
├── go.mod
├── go.sum
├── input.go
├── main.go
├── resource
│   └── ship.png
└── ship.go
```

为了防止歧义或者哪步搞错，下面给出完整的`main.go`和`ship.go`

main.go

```go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// 更新游戏内容显示

func (g *Game) Update() error {
	g.input.Update()
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

```

ship.go

```go
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
}

func NewShip() *Ship {
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

```

运行结果如下：
![image-20240904123515860](images/image-20240904123515860.png)

> 可见 飞船 是正好在绘制窗口的正中间显示的：达到了我们想要的结构

## 移动飞船

现在我们来实现使用左右方向键来控制飞船的移动。首先给飞船的类型增加x/y坐标字段：

> 修改ship.go里面的Ship结构体

```go
type Ship struct {
    // 与前面的代码一样
    x float64 // x坐标
    y float64 // y坐标
}
```

我们前面已经计算出飞船位于屏幕底部中心时的坐标，在创建飞船时将该坐标赋给xy：

> 修改ship.go里面的NewShip函数

```go
func NewShip(screenWidth, screenHeight int) *Ship {
  ship := &Ship{
    // ...
    x: float64(screenWidth-width) / 2,
    y: float64(screenHeight - height),
  }

  return ship
}
```

由于`NewShip`计算初始坐标需要屏幕尺寸，故增加屏幕宽、高两个参数，由`NewGame`方法传入：

> 也就是上面NewShip中所需要使用的两个参数： screenWidth和screeenHeight  修改game.go里面的NewGame函数

```golang
func NewGame() *Game {
  // 与上面的代码一样

  return &Game{
    input: &Input{},
    ship:  NewShip(cfg.ScreenWidth, cfg.ScreenHeight),  // 添加这里传入两个参数
    cfg:   cfg,
  }
}
```

然后我们在`Input`的`Update`方法中根据按下的是左方向键还是右方向键来更新飞船的坐标：

> 修改input.go里面的Update函数

```golang
type Input struct{}

func (i *Input) Update(ship *Ship) {
  if ebiten.IsKeyPressed(ebiten.KeyLeft) {
    ship.x -= 1
  } else if ebiten.IsKeyPressed(ebiten.KeyRight) {
    ship.x += 1
  }
}
// 传入的是一个ship对象，然后对其x进行修改就行了
```

由于需要修改飞船坐标，`Game.Update`方法调用`Input.Update`时需要传入飞船对象：

> 修改main.go里面的Update函数

```golang
func (g *Game) Update() error {
  g.input.Update(g.ship)
  return nil
}
```



> 修改好上面提到的几个点，结果如下：  go run .

> 测试发现无法 正常改变x轴的坐标位置
>
> 找到ship的绘制有问题修改如下

```go
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

```

移动如下：

![image-20240904172616773](images/image-20240904172616773.png)

注意到，目前有两个问题：

- 移动太慢
- 可以移出屏幕

修改`config.json`如下

```json
"shipSpeedFactor": 3
```

修改`config.go`里面的config结构体：

```golang
type Config struct {
    // 一样的代码
    ShipSpeedFactor float64    `json:"shipSpeedFactor"`
}
```

修改`Input.Update`方法，每次更新`ShipSpeedFactor`个像素：

```golang
func (i *Input) Update(ship *Ship, cfg *Config) {
  if ebiten.IsKeyPressed(ebiten.KeyLeft) {
    ship.x -= cfg.ShipSpeedFactor
  } else if ebiten.IsKeyPressed(ebiten.KeyRight) {
    ship.x += cfg.ShipSpeedFactor
  }
}
```

因为在`Input.Update`方法中需要访问配置，因此增加`Config`类型的参数，由`Game.Update`方法传入：

```golang
func (g *Game) Update() error {
  g.input.Update(g.ship, g.cfg)
  return nil
}
```

> 修改成功，小船的移动速度快了很多，也很方便的修改一些参数信息：关于小船的出界问题见Part2



# Part2

## 限制飞船的活动范围

上一篇文章还留了个尾巴，细心的同学应该发现了：飞船可以移动出屏幕！！！现在我们就来限制一下飞船的移动范围。我们规定飞船可以左右超过半个身位，如下图所示：

![img](images/ebiten12.png)

很容易计算得出，左边位置的x坐标为：

```
x = -W2/2
```

右边位置的坐标为：

```
x = W1 - W2/2
```

修改input.go的代码如下：

```golang
func (i *Input) Update(ship *Ship, cfg *Config) {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		ship.x -= cfg.ShipSpeedFactor
		if ship.x < -float64(ship.width)/2 {
			ship.x = -float64(ship.width) / 2
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		ship.x += cfg.ShipSpeedFactor
		if ship.x > float64(cfg.ScreenWidth)-float64(ship.width)/2 {
			ship.x = float64(cfg.ScreenWidth) - float64(ship.width)/2
		}
	}
}

// 增加两个if语句实现越界检测即可
```

效果自行实践。可以实现飞机主体尽量不越出界了



## 发射子弹

我们不用另外准备子弹的图片，直接画一个矩形就ok。为了可以灵活控制，我们将子弹的宽、高、颜色以及速率都用配置文件来控制：

修改 `config.json` 如下

```json
{
  "_comment1": "下面是屏幕信息",
  "screenWidth": 640,
  "screenHeight": 480,
  "title": "外星人入侵",
  "bgColor": {
    "r": 230,
    "g": 230,
    "b": 230,
    "a": 255
  },
  "_comment2": "下面是飞船速度",
  "shipSpeedFactor": 3,
  "_comment3": "下面是子弹信息",
  "bulletWidth": 3,
  "bulletHeight": 15,
  "bulletSpeedFactor": 2,
  "bulletColor": {
    "r": 80,
    "g": 80,
    "b": 80,
    "a": 255
  }
}
```

> 注意 json 中不能直接放入注释信息，json5中可以放入注释信息，json中可以通过上面的方式子放入注释信息

新增一个文件bullet.go，定义子弹的结构类型和New方法：

`bullet.go`

```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

// 创建子弹结构体
type Bullet struct {
	image       *ebiten.Image
	width       int
	height      int
	x           float64
	y           float64
	speedFactor float64
}

// 创建子弹创建的方法
func NewBullet(cfg *Config, ship *Ship) *Bullet {
	rect := image.Rect(0, 0, cfg.BulletWidth, cfg.BulletHeight)
	img := ebiten.NewImage(rect.Dx(), rect.Dy())
	img.Fill(cfg.BulletColor)

	return &Bullet{
		image:       img,
		width:       cfg.BulletWidth,
		height:      cfg.BulletHeight,
		x:           ship.x + float64(ship.width-cfg.BulletWidth)/2,
		y:           ship.y - float64(cfg.BulletHeight),
		speedFactor: cfg.BulletSpeedFactor,
	}
}

```

> 子弹的初始化信息 是设置好了

首先根据配置的宽高创建一个rect对象，然后调用`ebiten.NewImageWithOptions`创建一个`*ebiten.Image`对象。

子弹都是从飞船头部发出的，所以它的横坐标等于飞船中心的横坐标，左上角的纵坐标=屏幕高度-飞船高-子弹高。

随便增加子弹的绘制方法：

```golang
func (bullet *Bullet) Draw(screen *ebiten.Image) {
  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(bullet.x, bullet.y)
  screen.DrawImage(bullet.image, op)
}
```

> 把这段代码直接放入 `bullet.go` 最下面即可

我们在Game对象中增加一个map来管理子弹：

```golang
type Game struct {
  // -------省略-------
  bullets map[*Bullet]struct{}
}

func NewGame() *Game {
  return &Game{
    // -------省略-------
    bullets: make(map[*Bullet]struct{}),
  }
}
```

然后在`Draw`方法中，我们需要将子弹也绘制出来：

```golang
func (g *Game) Draw(screen *ebiten.Image) {
  screen.Fill(g.cfg.BgColor)
  g.ship.Draw(screen)
  for bullet := range g.bullets {
    bullet.Draw(screen)
  }
}
```

子弹位置如何更新呢？在`Game.Update`中更新，与飞船类似，只是飞船只能水平移动，而子弹只能垂直移动。

```golang
func (g *Game) Update() error {
  for bullet := range g.bullets {
    bullet.y -= bullet.speedFactor
  }
  // -------省略-------
}
```

子弹的更新、绘制逻辑都完成了，可是我们还没有子弹！现在我们就来实现按空格发射子弹的功能。我们需要在`Input.Update`方法中判断空格键是否按下，由于该方法需要访问Game对象的多个字段，干脆传入Game对象： 修改 input.go 文件

> 这里修改了 `input.go` 中的Update函数及其传入的参数，所以需要修改 `main.go` 中的参数信息：直接传入整个game对象

```golang
func (i *Input) Update(g *Game) {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.ship.x -= g.cfg.ShipSpeedFactor
		if g.ship.x < -float64(g.ship.width)/2 {
			g.ship.x = -float64(g.ship.width) / 2
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.ship.x += g.cfg.ShipSpeedFactor
		if g.ship.x > float64(g.cfg.ScreenWidth)-float64(g.ship.width)/2 {
			g.ship.x = float64(g.cfg.ScreenWidth) - float64(g.ship.width)/2
		}
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) {
		bullet := NewBullet(g.cfg, g.ship)
		g.addBullet(bullet)
	}
}
```

对应修改的main.go中的 Update如下：

```go
func (g *Game) Update() error {
	g.input.Update(g)
	for bullet := range g.bullets {
		bullet.y -= bullet.speedFactor
	}
	// 子弹移动
	return nil
}
```

给Game对象增加一个`addBullet`方法：

```golang
func (g *Game) addBullet(bullet *Bullet) {
  g.bullets[bullet] = struct{}{}
}
```

直接运行结果如下：

![image-20240906123306061](images/image-20240906123306061.png)

可见子弹绘制逻辑和绘制实现了，目前有两个问题：

- 无法一边移动一边发射，仔细看看`Input.Update`方法中的代码，你能发现什么问题吗？
- 子弹太多了，我们想要限制子弹的数量。

下面来逐一解决这些问题。

第一个问题很好解决，因为在KeyLeft/KeyRight/KeySpace这三个判断中我们用了if-else。这样会优先处理移动的操作。将KeySpace移到一个单独的if中即可：

```go
func (i *Input) Update(g *Game) {
  // -------省略-------
  if ebiten.IsKeyPressed(ebiten.KeySpace) {
    bullet := NewBullet(g.cfg, g.ship)
    g.addBullet(bullet)
  }
}
```

第二个问题，在配置中，增加同时最多存在的子弹数：

> 实测 23 颗子弹可以确保从飞机头部到界面的最高部分以及能发射出子弹，发射完毕之后也能无缝衔接上去

```json
{
  "maxBulletNum": 23
}
```

```golang
type Config struct {
  MaxBulletNum      int        `json:"maxBulletNum"`
}
```

然后我们在`Input.Update`方法中判断，如果目前存在的子弹数小于`MaxBulletNum`才能创建新的子弹：

```golang
if ebiten.IsKeyPressed(ebiten.KeySpace) {
  if len(g.bullets) < g.cfg.MaxBulletNum {
    bullet := NewBullet(g.cfg, g.ship)
    g.addBullet(bullet)
  }
}
```

那么修改好的`input.go` 里面的 Update 函数如下：
```go
func (i *Input) Update(g *Game) {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.ship.x -= g.cfg.ShipSpeedFactor
		if g.ship.x < -float64(g.ship.width)/2 {
			g.ship.x = -float64(g.ship.width) / 2
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.ship.x += g.cfg.ShipSpeedFactor
		if g.ship.x > float64(g.cfg.ScreenWidth)-float64(g.ship.width)/2 {
			g.ship.x = float64(g.cfg.ScreenWidth) - float64(g.ship.width)/2
		}
	}
	// 修改这里为else 也就是飞船发射子弹的逻辑和左右移动是无关的
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if len(g.bullets) < g.cfg.MaxBulletNum {
			bullet := NewBullet(g.cfg, g.ship)
			g.addBullet(bullet)
		}
	}
}
```

<img src="images/image-20240906123956292.png" alt="image-20240906123956292" style="zoom:50%;" />

数量好像被限制了，但是不是我们配置的23,原来`Input.Update()`的调用间隔太短了，导致我们一次space按键会发射多个子弹。我们可以控制两个子弹之间的时间间隔。同样用配置文件来控制（单位毫秒）： 实际测试 150ms 比较合适

```json
{
  "bulletInterval": 150
}
```

```golang
type Config struct {
  BulletInterval    int64      `json:"bulletInterval"`
}
```

在`Input`结构中增加一个`lastBulletTime`字段记录上次发射子弹的时间：

```golang
type Input struct {
  lastBulletTime time.Time
}
```

距离上次发射子弹的时间大于`BulletInterval`毫秒，才能再次发射，发射成功之后更新时间

```golang
func (i *Input) Update(g *Game) {
  // -------省略-------
  	if ebiten.IsKeyPressed(ebiten.KeySpace) {
      if len(g.bullets) < g.cfg.MaxBulletNum &&
        time.Since(i.lastBulletTime).Milliseconds() > g.cfg.BulletInterval {
        bullet := NewBullet(g.cfg, g.ship)
        g.addBullet(bullet)
        i.lastBulletTime = time.Now()
      }
	}
}
```

又出现了一个问题，23个子弹飞出屏幕外之后还是不能发射子弹。我们需要把离开屏幕的子弹删除。这适合在`Game.Update`函数中做：

```golang
func (g *Game) Update() error {
	g.input.Update(g)
	for bullet := range g.bullets {
		bullet.y -= bullet.speedFactor
		if bullet.outOfScreen() {
			delete(g.bullets, bullet)
		}
	}
	// 子弹移动
	return nil
}
```

为了`Bullet`添加判断是否处于屏幕外的方法：

```golang
func (bullet *Bullet) outOfScreen() bool {
  return bullet.y < -float64(bullet.height)
}
```

> 结果是 全屏幕可以无缝衔接子弹