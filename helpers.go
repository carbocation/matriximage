package matriximage

import (
	"image"
	"image/color"
	"math"

	"github.com/mjibson/go-dsp/dsputils"
)

func toImage(matrix *dsputils.Matrix, imaginary bool) *Image {
	vals := matrix.To2D()
	dims := matrix.Dimensions()

	m := image.NewGray16(image.Rect(0, 0, dims[1], dims[0]))

	for y := range vals {
		for x := range vals[y] {
			part := real(vals[y][x])
			if imaginary {
				part = imag(vals[y][x])
			}

			if part > float64(math.MaxUint16) {
				part = math.MaxUint16
			} else if part < 0 {
				part = 0
			}

			v := uint16(part)

			m.SetGray16(x, y, color.Gray16{v})
		}
	}

	return &Image{m}
}

func toRealImage(matrix *dsputils.Matrix) *Image {
	return toImage(matrix, false)
}

func toImaginaryImage(matrix *dsputils.Matrix) *Image {
	return toImage(matrix, true)
}
