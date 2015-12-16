package matriximage

import (
	"fmt"
	//"math"

	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
)

// Consists of isolated real and imaginary components
type FrequencyImage struct {
	Real *Image
	Imag *Image

	matrix *dsputils.Matrix
}

func (fi *FrequencyImage) toGrayMatrix() (*dsputils.Matrix, error) {
	if fi.Real.Bounds() != fi.Imag.Bounds() {
		return nil, fmt.Errorf("Bounds for Real component differ from those of the Imag component")
	}

	rl := fi.Real.toGrayMatrix().To2D()
	im := fi.Imag.toGrayMatrix().To2D()

	dims := fi.Real.toGrayMatrix().Dimensions()

	//scale := math.Sqrt(float64(dims[0] * dims[1]))
	scale := 1.0

	matrix := dsputils.MakeEmptyMatrix(dims)

	for i := 0; i < dims[1]; i++ {
		for j := 0; j < dims[0]; j++ {
			// The real part of the real matrix is the real part of the frequency matrix
			// The real part of the imag matrix is the imag part of the frequency matrix
			v := complex(scale*real(rl[j][i]), scale*real(im[j][i]))

			matrix.SetValue(v, []int{j, i})
		}
	}

	return matrix, nil
}

func (fi *FrequencyImage) IFFTN() (*Image, error) {
	return toRealImage(fft.IFFTN(fi.matrix)), nil

	matrix, err := fi.toGrayMatrix()
	if err != nil {
		return nil, err
	}

	inverse := fft.IFFTN(matrix)

	return toRealImage(inverse), nil
}
