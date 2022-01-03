package shadow2d

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	edgeNorth = 0
	edgeSouth = 1
	edgeEast  = 2
	edgeWest  = 3
)

var (
	polygons   = newPolygonMap()
	lastEdgeId edgeId
)

type edgeId int

func getNextEdgeId() edgeId {
	lastEdgeId++

	return lastEdgeId
}

type edge struct {
	id edgeId
	line
}

func newEdge(startX, startY, endX, endY float64) *edge {
	return &edge{id: getNextEdgeId(), line: newLine(newPoint(startX, startY), newPoint(endX, endY))}
}

type cell struct {
	address
	// Карта граней где ключ это сторона, а значение это id грани.
	edges map[int]edgeId
}

func newCell(address address) *cell {
	return &cell{
		address: address,
		edges:   map[int]edgeId{},
	}
}

type polygonMap struct {
	edges map[edgeId]*edge
	cells map[address]*cell
}

func newPolygonMap() *polygonMap {
	p := &polygonMap{edges: map[edgeId]*edge{}, cells: map[address]*cell{}}

	columns := ScreenWidth / TileSize
	rows := ScreenHeight / TileSize

	for column := 0; column < columns; column++ {
		for row := 0; row < rows; row++ {
			address := newAddress(column, row)
			p.cells[address] = newCell(address)
		}
	}

	return p
}

func (p *polygonMap) update() {
	columns := ScreenWidth / TileSize
	rows := ScreenHeight / TileSize

	p.edges = map[edgeId]*edge{}
	p.cells = map[address]*cell{}
	lastEdgeId = 0

	tileSize := float64(TileSize)

	for column := 0; column < columns; column++ {
		for row := 0; row < rows; row++ {
			address := newAddress(column, row)

			if !tiles.has(address.get()) {
				continue
			}

			cell := newCell(address)
			p.cells[cell.address] = cell

			// Проверяем есть ли тайл сверху. Если тайла нет, то необходимо добавить ребро.
			if !tiles.has(column, row-1) {
				leftCell, foundCell := p.getCell(column-1, row)
				if foundCell {
					edge, foundEdge := p.edges[leftCell.edges[edgeNorth]]
					if foundEdge {
						cell.edges[edgeNorth] = edge.id
						edge.end.x = edge.end.x + tileSize
					} else {
						x, y := convertAddressToCoordinates(column, row)
						edge := newEdge(x, y, x+tileSize, y)
						cell.edges[edgeNorth] = edge.id
						p.edges[edge.id] = edge
					}
				} else {
					x, y := convertAddressToCoordinates(column, row)
					edge := newEdge(x, y, x+tileSize, y)
					cell.edges[edgeNorth] = edge.id
					p.edges[edge.id] = edge
				}
			}

			// Проверяем есть ли тайл снизу. Если тайла нет, то необходимо добавить ребро.
			if !tiles.has(column, row+1) {
				leftCell, foundCell := p.getCell(column-1, row)
				if foundCell {
					edge, foundEdge := p.edges[leftCell.edges[edgeSouth]]
					if foundEdge {
						cell.edges[edgeSouth] = edge.id
						edge.end.x += tileSize
					} else {
						x, y := convertAddressToCoordinates(column, row)
						edge := newEdge(x, y+tileSize, x+tileSize, y+tileSize)
						cell.edges[edgeSouth] = edge.id
						p.edges[edge.id] = edge
					}
				} else {
					x, y := convertAddressToCoordinates(column, row)
					edge := newEdge(x, y+tileSize, x+tileSize, y+tileSize)
					cell.edges[edgeSouth] = edge.id
					p.edges[edge.id] = edge
				}
			}

			// Проверяем есть ли тайл справа. Если тайла нет, то необходимо добавить ребро.
			if !tiles.has(column+1, row) {
				upperCell, foundCell := p.getCell(column, row-1)
				if foundCell {
					edge, foundEdge := p.edges[upperCell.edges[edgeEast]]
					if foundEdge {
						cell.edges[edgeEast] = edge.id
						edge.end.y += tileSize
					} else {
						x, y := convertAddressToCoordinates(column, row)
						edge := newEdge(x+tileSize, y, x+tileSize, y+tileSize)
						cell.edges[edgeEast] = edge.id
						p.edges[edge.id] = edge
					}
				} else {
					x, y := convertAddressToCoordinates(column, row)
					edge := newEdge(x+tileSize, y, x+tileSize, y+tileSize)
					cell.edges[edgeEast] = edge.id
					p.edges[edge.id] = edge
				}
			}

			// Проверяем есть ли тайл слева. Если тайла нет, то необходимо добавить ребро.
			if !tiles.has(column-1, row) {
				upperCell, foundCell := p.getCell(column, row-1)
				if foundCell {
					edge, foundEdge := p.edges[upperCell.edges[edgeWest]]
					if foundEdge {
						cell.edges[edgeWest] = edge.id
						edge.end.y += float64(TileSize)
					} else {
						x, y := convertAddressToCoordinates(column, row)
						edge := newEdge(x, y, x, y+tileSize)
						cell.edges[edgeWest] = edge.id
						p.edges[edge.id] = edge
					}
				} else {
					x, y := convertAddressToCoordinates(column, row)
					edge := newEdge(x, y, x, y+tileSize)
					cell.edges[edgeWest] = edge.id
					p.edges[edge.id] = edge
				}
			}
		}
	}
}

func (p *polygonMap) draw(screen *ebiten.Image) {
	for _, edge := range p.edges {
		ebitenutil.DrawLine(screen, edge.start.x, edge.start.y, edge.end.x, edge.end.y, hexToRGB(edgeColor))
		ebitenutil.DrawRect(screen, edge.start.x-4, edge.start.y-4, 8, 8, hexToRGB("ff0000"))
		ebitenutil.DrawRect(screen, edge.end.x-4, edge.end.y-4, 8, 8, hexToRGB("ff0000"))
	}
}

func (p *polygonMap) hasCell(column, row int) bool {
	_, ok := p.cells[newAddress(column, row)]

	return ok
}

func (p *polygonMap) getCell(column, row int) (*cell, bool) {
	cell, ok := p.cells[newAddress(column, row)]

	return cell, ok
}

func (p *polygonMap) getVertexes() []point {
	unique := map[point]bool{}
	var vertexes []point

	for _, edge := range p.edges {
		if _, ok := unique[edge.start]; !ok {
			vertexes = append(vertexes, edge.start)
			unique[edge.start] = true
		}

		if _, ok := unique[edge.end]; !ok {
			vertexes = append(vertexes, edge.end)
			unique[edge.end] = true
		}
	}

	return vertexes
}
