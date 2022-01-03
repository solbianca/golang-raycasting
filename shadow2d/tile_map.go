package shadow2d

import "github.com/hajimehoshi/ebiten/v2"

var (
	tiles *tileMap
)

func init() {
	tiles = newTileMap()
}

type tileMap struct {
	tiles map[address]*tile
}

func newTileMap() *tileMap {
	t := &tileMap{tiles: map[address]*tile{}}

	columns := ScreenWidth / TileSize
	rows := ScreenHeight / TileSize

	for column := 0; column < columns; column++ {
		address := newAddress(column, 0)
		t.tiles[address] = newTile(address)

		address = newAddress(column, rows-1)
		t.tiles[address] = newTile(address)
	}

	for row := 0; row < rows; row++ {
		address := newAddress(0, row)
		t.tiles[address] = newTile(address)

		address = newAddress(columns-1, row)
		t.tiles[address] = newTile(address)
	}

	return t
}

func (t *tileMap) update() {
	column, row := convertCoordinatesToAddress(input.mouse.getCoordinates())

	if input.mouse.isRightButtonReleased {
		if t.has(column, row) {
			t.delete(column, row)
		} else {
			t.set(column, row)
		}
	}
}

func (t *tileMap) draw(screen *ebiten.Image) {
	for _, tile := range t.tiles {
		tile.draw(screen)
	}
}

func (t *tileMap) set(column, row int) {
	address := newAddress(column, row)
	t.tiles[address] = newTile(address)
}

func (t *tileMap) has(column, row int) bool {
	_, ok := t.tiles[newAddress(column, row)]

	return ok
}

func (t *tileMap) delete(column, row int) {
	delete(t.tiles, newAddress(column, row))
}

func (t *tileMap) get(column, row int) (*tile, bool) {
	tile, ok := t.tiles[newAddress(column, row)]

	return tile, ok
}

type tile struct {
	address
	point

	img *ebiten.Image
}

func newTile(address address) *tile {
	x, y := convertAddressToCoordinates(address.get())
	t := &tile{address: address, point: newPoint(x, y)}

	t.img = ebiten.NewImage(TileSize, TileSize)
	t.img.Fill(hexToRGB(tileColor))

	return t
}

func (t *tile) update() {
}

func (t *tile) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(t.point.x, t.point.y)

	screen.DrawImage(t.img, op)
}
