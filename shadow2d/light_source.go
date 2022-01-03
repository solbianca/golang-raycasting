package shadow2d

import (
	"image/color"
	"math"
	"sort"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	light          *lightSource
	radius         = 1000.0
	triangleImage  = ebiten.NewImage(ScreenWidth, ScreenHeight)
	shadowImage    = ebiten.NewImage(ScreenWidth, ScreenHeight)

	visibleVertexes []vertex
)

func init() {
	light = newLightSource()
	visibleVertexes = []vertex{}

	triangleImage.Fill(color.White)
}

type lightSource struct {
	sourceImg *ebiten.Image
}

func newLightSource() *lightSource {
	ls := &lightSource{}

	dc := gg.NewContext(10, 10)
	dc.DrawCircle(5, 5, 4)
	dc.SetHexColor("ffff00")
	dc.Fill()
	ls.sourceImg = ebiten.NewImageFromImage(dc.Image())

	return ls
}

func (l *lightSource) update() {
	if !input.mouse.isLeftButtonPressed {
		return
	}

	visibleVertexes = []vertex{}

	sourceX, sourceY := float64(input.mouse.x), float64(input.mouse.y)
	sourcePoint := newPoint(sourceX, sourceY)

	for _, vertex := range polygons.getVertexes() {
		var rdx, rdy float64

		rdx = vertex.x - sourcePoint.x
		rdy = vertex.y - sourcePoint.y

		baseAngle := math.Atan2(rdy, rdx)

		angle := 0.0

		for j := 0; j < 3; j++ {
			switch j {
			case 0:
				angle = baseAngle - 0.0001
			case 1:
				angle = baseAngle
			case 2:
				angle = baseAngle + 0.0001
			}

			rdx = radius * math.Cos(angle)
			rdy = radius * math.Sin(angle)

			minT1 := math.MaxFloat64
			minPX, minPY, minAngle := 0.0, 0.0, 0.0
			isIntersect := false

			for _, edge := range polygons.edges {
				sdx := edge.end.x - edge.start.x
				sdy := edge.end.y - edge.start.y

				if math.Abs(sdx-rdx) > 0.0 && math.Abs(sdy-rdy) > 0.0 {
					// t2 is normalised distance from line segment start to line segment end of intersect point
					t2 := (rdx*(edge.start.y-sourcePoint.y) + (rdy * (sourcePoint.x - edge.start.x))) / (sdx*rdy - sdy*rdx)
					// t1 is normalised distance from source along ray to ray length of intersect point
					t1 := (edge.start.x + sdx*t2 - sourcePoint.x) / rdx

					if t1 > 0 && t2 >= 0 && t2 <= 1.0 {
						if t1 < minT1 {
							minT1 = t1
							minPX = sourcePoint.x + rdx*t1
							minPY = sourcePoint.y + rdy*t1
							minAngle = math.Atan2(minPY-sourcePoint.y, minPX-sourcePoint.x)
							isIntersect = true
						}
					}
				}
			}

			if isIntersect {
				visibleVertexes = append(visibleVertexes, newVertex(minAngle, newPoint(minPX, minPY)))
			}
		}
	}

	sort.Slice(visibleVertexes, func(i, j int) bool {
		return visibleVertexes[i].angle > visibleVertexes[j].angle
	})
}

func cross(aPoint1, aPoint2, bPoint1, bPoint2 point) {

}

func (l *lightSource) draw(screen *ebiten.Image) {
	x, y := float64(input.mouse.x), float64(input.mouse.y)

	if !input.mouse.isLeftButtonPressed {
		return
	}

	// draw light source
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x-5, y-5)
	screen.DrawImage(l.sourceImg, op)

	shadowImage.Fill(color.Black)

	// Subtract ray triangles from shadow
	opt := &ebiten.DrawTrianglesOptions{}
	opt.Address = ebiten.AddressRepeat
	opt.CompositeMode = ebiten.CompositeModeSourceOut
	for i, vertex := range visibleVertexes {
		nextVertex := visibleVertexes[(i+1)%len(visibleVertexes)]

		// Draw triangle of area between rays
		v := rayVertices(x, y, nextVertex.x, nextVertex.y, vertex.x, vertex.y)
		shadowImage.DrawTriangles(v, []uint16{0, 1, 2}, triangleImage, opt)
	}

	// Draw shadow
	op = &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.7)
	screen.DrawImage(shadowImage, op)
}

type line struct {
	start point
	end   point
}

func newLine(start, end point) line {
	return line{start: start, end: end}
}

type vertex struct {
	angle float64
	point
}

func newVertex(angle float64, point point) vertex {
	return vertex{angle: angle, point: point}
}

func rayVertices(x1, y1, x2, y2, x3, y3 float64) []ebiten.Vertex {
	return []ebiten.Vertex{
		{DstX: float32(x1), DstY: float32(y1), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x2), DstY: float32(y2), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x3), DstY: float32(y3), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
}
