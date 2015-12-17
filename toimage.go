package matriximage

import (
	"image"
	"image/color"
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/dsputils"
)

const (
	Real = iota
	Imaginary
	Amplitude
	LogAmplitude
	Phase
)

// Data is lost when using this on the FFT output because
// there are values outside the range of UINT8
func toImage(matrix *dsputils.Matrix, whichPart int) *Image {
	vals := matrix.To2D()
	dims := matrix.Dimensions()

	m := image.NewGray(image.Rect(0, 0, dims[1], dims[0]))

	for y := range vals {
		for x := range vals[y] {
			var part float64
			if whichPart == Real {
				part = real(vals[y][x])
			} else if whichPart == Imaginary {
				part = imag(vals[y][x])
			} else if whichPart == Amplitude {
				part = cmplx.Abs(vals[(y+dims[0]/2)%dims[0]][(x+dims[1]/2)%dims[1]]) / float64(dims[0]*dims[1])

				//part += math.MaxUint8 / 2
			} else if whichPart == LogAmplitude {
				part = math.Log(cmplx.Abs(vals[(y+dims[0]/2)%dims[0]][(x+dims[1]/2)%dims[1]]) / float64(dims[0]*dims[1]))
			} else if whichPart == Phase {
				part = cmplx.Phase(vals[(y+dims[0]/2)%dims[0]][(x+dims[1]/2)%dims[1]]) //* 1.0 / math.Pi //* float64(dims[1]*dims[0])

				part += math.MaxUint8 / 2
			}

			if part > float64(math.MaxUint8) {
				part = math.MaxUint8
			} else if part < 0 {
				part = 0
			}

			v := uint8(part)

			m.SetGray(x, y, color.Gray{v})
		}
	}

	return &Image{m}
}

func positiveMod(n, k int) int {
	if n < 0 {
		return k - ((-1 * n) % k)
	}

	return n % k
}

func toRealImage(matrix *dsputils.Matrix) *Image {
	return toImage(matrix, Real)
}

func toImaginaryImage(matrix *dsputils.Matrix) *Image {
	return toImage(matrix, Imaginary)
}

func toAmplitudeImage(matrix *dsputils.Matrix) *Image {
	return toImage(matrix, Amplitude)
}

func toLogAmplitudeImage(matrix *dsputils.Matrix) *Image {
	return toImage(matrix, LogAmplitude)
}

func toPhaseImage(matrix *dsputils.Matrix) *Image {
	return toImage(matrix, Phase)
}
