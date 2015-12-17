package matriximage

import (
	//"fmt"
	//"math"

	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
)

// Consists of isolated real and imaginary components
type FrequencyImage struct {
	Amp   *Image
	Phase *Image

	matrix *dsputils.Matrix
}

func (fi *FrequencyImage) toGrayMatrix() (*dsputils.Matrix, error) {

	amp := fi.Amp.toGrayMatrix()
	phase := fi.Phase.toGrayMatrix()
	return mergeAmpPhase(amp, phase), nil

	/*

		if fi.Amp.Bounds() != fi.Phase.Bounds() {
			return nil, fmt.Errorf("Bounds for Real component differ from those of the Imag component")
		}

		amp := fi.Amp.toGrayMatrix().To2D()
		phas := fi.Phase.toGrayMatrix().To2D()

		dims := fi.Amp.toGrayMatrix().Dimensions()

		matrix := dsputils.MakeEmptyMatrix(dims)

		for i := 0; i < dims[1]; i++ {
			for j := 0; j < dims[0]; j++ {
				// The real part of the real matrix is the real part of the frequency matrix
				// The real part of the imag matrix is the imag part of the frequency matrix
				radius := real(amp[j][i])

				rl := radius * math.Cos(real(phas[j][i]))
				im := radius * math.Sin(real(phas[j][i]))

				v := complex(rl, im)

				matrix.SetValue(v, []int{j, i})
			}
		}

		return matrix, nil
	*/
}

func (fi *FrequencyImage) IFFTN() (*Image, error) {
	// Easy option that ignores the saved images
	return toRealImage(fft.IFFTN(fi.matrix)), nil

	/*
		matrix, err := fi.toGrayMatrix()
		if err != nil {
			return nil, err
		}

		inverse := fft.IFFTN(matrix)

		return toRealImage(inverse), nil
	*/
}
