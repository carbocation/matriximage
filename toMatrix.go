package matriximage

import (
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/dsputils"
)

func splitAmpPhase(fftn *dsputils.Matrix) (*dsputils.Matrix, *dsputils.Matrix) {
	dims := fftn.Dimensions()

	abs := dsputils.MakeEmptyMatrix(dims)
	phase := dsputils.MakeEmptyMatrix(dims)

	// Y
	for i := 0; i < dims[1]; i++ {
		// X
		for j := 0; j < dims[0]; j++ {
			here := []int{j, i}

			valueHere := fftn.Value(here)

			abs.SetValue(complex(cmplx.Abs(valueHere), 0), here)
			phase.SetValue(complex(cmplx.Phase(valueHere), 0), here)
		}
	}

	return abs, phase
}

// Merges amplitude and phase matrices to recover the FFT matrix
// Reminder: layout is matrix[y][x]
func mergeAmpPhase(amp, phase *dsputils.Matrix) *dsputils.Matrix {
	ad := amp.Dimensions()
	pd := phase.Dimensions()

	if len(ad) != 2 || len(pd) != 2 {
		return nil
	}

	for k := range ad {
		if ad[k] != pd[k] {
			return nil
		}
	}

	n := dsputils.MakeEmptyMatrix(ad)

	// Y
	for i := 0; i < ad[1]; i++ {
		// X
		for j := 0; j < ad[0]; j++ {
			here := []int{j, i}

			radius := real(amp.Value(here))

			rl := radius * math.Cos(real(phase.Value(here)))
			im := radius * math.Sin(real(phase.Value(here)))

			v := complex(rl, im)

			n.SetValue(v, here)
		}
	}

	return n
}
