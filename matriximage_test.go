package matriximage

import (
	"fmt"
	"os"
	"testing"

	"github.com/mjibson/go-dsp/dsputils"
)

const Case = "sin"

func TestMask(t *testing.T) {
	m, err := FromFile(fmt.Sprintf("./images/%s.gif", Case))
	if err != nil {
		t.Fatal(err)
	}

	fourier := m.DFT()

	maskMatrix := dsputils.MakeEmptyMatrix(fourier.Dimensions()) //fourier.AmplitudeImage().toGrayMatrix()
	dims := maskMatrix.Dimensions()

	// Create a square mask, where only the middle part of the amplitude is recovered.
	for y := 0; y < dims[0]; y++ {
		for x := 0; x < dims[1]; x++ {
			rad := 1
			//if (y >= dims[0]/2-rad && y <= dims[0]/2+rad) && (x >= dims[1]/2-rad && x <= 1+dims[1]/2+rad) {
			if (y < dims[0]/2-rad || y > dims[0]/2+rad) || (x < dims[1]/2-rad || x > 1+dims[1]/2+rad) {
				maskMatrix.SetValue(complex(0, 0), []int{y, x})
			} else {
				maskMatrix.SetValue(complex(float64(MaxUint), 0), []int{y, x})
			}
		}
	}

	maskedFourier, err := fourier.ApplyMatrixMask(maskMatrix)
	if err != nil {
		t.Fatal(err)
	}

	testSave(*maskedFourier.AmplitudeImage(), "masked-amplitude")
	testSave(*maskedFourier.BrighterAmplitudeImage(), "masked-brighter-amplitude")
	testSave(*maskedFourier.PhaseImage(), "masked-phase")
	testSave(*maskedFourier.IDFTImage(), "masked-recovered")
}

func TestFourierImage(t *testing.T) {
	m, err := FromFile(fmt.Sprintf("./images/%s.gif", Case))
	if err != nil {
		t.Fatal(err)
	}

	fourier := m.DFT()

	testSave(*fourier.AmplitudeImage(), "amp")
	testSave(*fourier.BrighterAmplitudeImage(), "logamp")
	testSave(*fourier.PhaseImage(), "phase")
	testSave(*fourier.IDFTImage(), "recovered")
}

func TestToMatrix(t *testing.T) {
	m, err := FromFile("./images/square.gif")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(m.Image)
}

func TestLoad(t *testing.T) {
	_, err := FromFile("./images/square.gif")
	if err != nil {
		t.Fatal(err)
	}
}

func testSave(m Image, name string) {
	os.MkdirAll("./images.out", 0755)

	filename := fmt.Sprintf("./images.out/%s-%s.png", Case, name)

	m.ToFile(filename)
}
