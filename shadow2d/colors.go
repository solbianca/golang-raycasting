package shadow2d

import (
	"image/color"
	"strconv"
)

var (
	tileColor = "2476FF"
	edgeColor = "ff0000"
)

func hexToRGB(hexStr string) color.Color {
	u, err := strconv.ParseUint(hexStr, 16, 0)
	if err != nil {
		panic(err)
	}

	return color.RGBA{
		R: uint8(u & 0xff0000 >> 16),
		G: uint8(u & 0xff00 >> 8),
		B: uint8(u & 0xff),
		A: 255,
	}
}
