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

```go
import (
    _ "golang.org/x/image/bmp"
)

type Ship struct {
  image  *ebiten.Image
  width  int
  height int
}

func NewShip() *Ship {
  img, _, err := ebitenutil.NewImageFromFile("../images/ship.bmp")
  if err != nil {
    log.Fatal(err)
  }

  width, height := img.Size()
  ship := &Ship{
    image:  img,
    width:  width,
    height: height,
  }

  return ship
}
```