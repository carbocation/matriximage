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
	BrighterAmplitude
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
//
// Note that amplitude values are always zero or positive,
// which are amenable to representation in an image. Phase
// values are not, making them more challenging to represent
// in an image. Generally, phase is not manipulated.
func toImage(matrix *dsputils.Matrix, whichPart int) *Image {
	vals := matrix.To2D()
	dims := matrix.Dimensions()

	m := image.NewGray16(image.Rect(0, 0, dims[1], dims[0]))

	// Determine a scale factor
	maxAmp := 0.0
	minPhase, maxPhase := 0.0, 0.0
	for y := range vals {
		for x := range vals[y] {
			if v := cmplx.Abs(vals[y][x]); v > maxAmp {
				maxAmp = v
			}

			v := cmplx.Phase(vals[y][x])
			if v < minPhase {
				minPhase = v
			}
			if v > maxPhase {
				maxPhase = v
			}
		}
	}

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
				scale := 1.0 / float64(dims[0]*dims[1])

				// Shift to get centered coordinates
				part = cmplx.Abs(vals[shiftedY][shiftedX]) * scale

			} else if whichPart == BrighterAmplitude {
				scale := 1.0 / math.Hypot(float64(dims[0]), float64(dims[1]))

				// Shift to get centered coordinates
				part = cmplx.Abs(vals[shiftedY][shiftedX]) * scale
				//part = math.Log(cmplx.Abs(vals[shiftedY][shiftedX])/math.Hypot(float64(dims[0]), float64(dims[1]))) / math.Log(maxAmp) * float64(math.MaxUint8)

			} else if whichPart == Phase {
				offset := 0.0 - minPhase
				scale := float64(MaxUint) / (maxPhase - minPhase)

				// Shift to get centered coordinates
				part = (offset + cmplx.Phase(vals[shiftedY][shiftedX])) * scale //* 1.0 / math.Pi //* float64(dims[1]*dims[0])

			}

			if part > float64(MaxUint) {
				part = MaxUint
			} else if part < 0 {
				part = 0
			}

			v := uint16(part)

			m.SetGray16(x, y, color.Gray16{v})
		}
	}

	return &Image{Image: m}
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

func toBrighterAmplitude(matrix *dsputils.Matrix) *Image {
	return toImage(matrix, BrighterAmplitude)
}

func toPhaseImage(matrix *dsputils.Matrix) *Image {
	return toImage(matrix, Phase)
}
