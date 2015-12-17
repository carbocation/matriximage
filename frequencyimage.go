package matriximage

import (
	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
)

type FourierImage struct {
	*dsputils.Matrix // Is expected to be the result of a Fourier transform
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
