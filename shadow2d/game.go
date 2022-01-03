package shadow2d

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	ScreenWidth  = 512
	ScreenHeight = 512
	TileSize     = 32

	bgImage *ebiten.Image
)

func init() {
	path := "assets/rock.png"
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("file not founded by path [%a], %v", path, err))
	}
	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		log.Fatal(err)
	}
	bgImage = ebiten.NewImageFromImage(img)
}

type gameObject interface {
	update()
	draw(screen *ebiten.Image)
}

type Game struct {
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Update() error {
	input.update()

	tiles.update()
	polygons.update()
	light.update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	screen.DrawImage(bgImage, nil)

	drawGrid(screen)

	tiles.draw(screen)
	polygons.draw(screen)
	light.draw(screen)

	column, row := convertCoordinatesToAddress(input.mouse.x, input.mouse.y)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Clumn:[%d] Row:[%d]", column, row))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nEdges:[%d]", len(polygons.edges)))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\nRays:[%d]", len(visibleVertexes)))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\nCursor position:[%d]:[%d]", input.mouse.x, input.mouse.y))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func drawGrid(screen *ebiten.Image) {
	for step := 0; step <= ScreenWidth; step += TileSize {
		ebitenutil.DrawLine(screen, float64(step), float64(0), float64(step), float64(ScreenHeight), hexToRGB("333333"))
	}

	for step := 0; step <= ScreenHeight; step += TileSize {
		ebitenutil.DrawLine(screen, float64(0), float64(step), float64(ScreenWidth), float64(step), hexToRGB("333333"))
	}
}
