package shadow2d

func convertCoordinatesToAddress(x, y int) (column, row int) {
	column = x / TileSize
	row = y / TileSize

	return column, row
}

func convertAddressToCoordinates(column, row int) (x, y float64) {
	x = float64(column * TileSize)
	y = float64(row * TileSize)

	return x, y
}

type address struct {
	column, row int
}

func newAddress(column int, row int) address {
	return address{column: column, row: row}
}

func (a address) get() (column int, row int) {
	return a.column, a.row
}

type point struct {
	x, y float64
}

func newPoint(x, y float64) point {
	return point{x: x, y: y}
}

func (p point) get() (x float64, y float64) {
	return p.x, p.y
}
