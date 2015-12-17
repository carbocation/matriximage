package matriximage

import (
	"fmt"

	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
)

type FourierImage struct {
	*dsputils.Matrix // Is expected to be the result of a Fourier transform
}

// Generate a new FourierImage by applying a mask to the current one
// The mask is assumed to have its frequency origin at the center, which
// is shifted by N/2 and M/2 from the top left origin
func (f FourierImage) ApplyImageMask(img *Image) (*FourierImage, error) {
	mask := img.ToGrayMatrix()

	return f.ApplyMatrixMask(mask)
}

func (f FourierImage) ApplyMatrixMask(mask *dsputils.Matrix) (*FourierImage, error) {
	dims := f.Dimensions()
	maskDims := mask.Dimensions()

	if len(dims) != len(maskDims) {
		return nil, fmt.Errorf("Mask had different number of dimensions from image")
	}
	for i := 0; i < len(dims); i++ {
		if dims[i] != maskDims[i] {
			return nil, fmt.Errorf("Size of mask dimension %d (%d) differs from size of image dimension %d (%d)", i, maskDims[i], i, dims[i])
		}
	}

	m := dsputils.MakeEmptyMatrix(dims)

	// While the mask has been shifted with respect to the orientation of the matrix,
	// the internal matrix representation has not. So, we pull the mask value from
	// a shifted value, and pull the matrix value from an unshifted value.
	for y := 0; y < dims[0]; y++ {
		for x := 0; x < dims[1]; x++ {
			shiftedY, shiftedX := (y+dims[0]/2)%dims[0], (x+dims[1]/2)%dims[1]

			v := f.Value([]int{y, x})
			maskV := mask.Value([]int{shiftedY, shiftedX})

			// Multiply the value at this pixel by a factor from 0-1, depending on
			// the mask's amplitude compared to the max possible amplitude
			product := v * complex((real(maskV)/MaxUint), (real(maskV)/MaxUint))

			m.SetValue(product, []int{y, x})
		}
	}

	return &FourierImage{Matrix: m}, nil
}

// For making masks
func (f FourierImage) AmplitudeImage() *Image {
	return toAmplitudeImage(f.Matrix)
}

// For making masks
func (f FourierImage) BrighterAmplitudeImage() *Image {
	return toBrighterAmplitude(f.Matrix)
}

// For... visualization?
func (f FourierImage) PhaseImage() *Image {
	return toPhaseImage(f.Matrix)
}

// Inverse FFT image output
func (f FourierImage) IDFTImage() *Image {
	return toRealImage(f.idft())
}

// Inverse FFT
func (f FourierImage) idft() *dsputils.Matrix {
	return fft.IFFTN(f.Matrix)
}
