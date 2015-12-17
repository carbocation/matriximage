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
	RealFromDFT
	Amplitude
	LogAmplitude
	Phase
)

// Several things to note:
//
// 1. Data is lost when using this on the FFT output because
// there are values outside the range of UINT8
//
// 2. Amplitude, LogAmplitude, and Phase all imply that the
// signal being plotted is a result of a DFT operation. In
// that case, we shift over by X/2 and down by Y/2; i.e.,
// we aim to recenter the DFT so the star graph is centered.
// This means that these signals have to be un-shifted.
func toImage(matrix *dsputils.Matrix, whichPart int) *Image {
	vals := matrix.To2D()
	dims := matrix.Dimensions()

	m := image.NewGray(image.Rect(0, 0, dims[1], dims[0]))

	for y := range vals {
		for x := range vals[y] {
			shiftedY, shiftedX := (y+dims[0]/2)%dims[0], (x+dims[1]/2)%dims[1]

			var part float64
			if whichPart == Real {
				part = real(vals[y][x])
			} else if whichPart == RealFromDFT {
				// Re-shift so we get back the original coordinates
				part = real(vals[shiftedY][shiftedX])
			} else if whichPart == Amplitude {
				// Shift to get centered coordinates
				part = cmplx.Abs(vals[shiftedY][shiftedX]) / float64(dims[0]*dims[1])

				//part += math.MaxUint8 / 2
			} else if whichPart == LogAmplitude {
				// Shift to get centered coordinates
				part = math.Log(cmplx.Abs(vals[shiftedY][shiftedX]) / float64(dims[0]*dims[1]))
			} else if whichPart == Phase {
				// Shift to get centered coordinates
				part = cmplx.Phase(vals[shiftedY][shiftedX]) //* 1.0 / math.Pi //* float64(dims[1]*dims[0])

				//part += math.MaxUint8 / 2
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

func toRealFromDFTImage(matrix *dsputils.Matrix) *Image {
	return toImage(matrix, RealFromDFT)
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
