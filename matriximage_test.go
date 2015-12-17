package matriximage

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/mjibson/go-dsp/dsputils"
)

const Case = "rand"

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
			//if y < 62*dims[0]/128 || x < 62*dims[1]/128 || y > 65*dims[0]/128 || x > 65*dims[1]/128 {
			if (y >= dims[0]/2-2 && y <= dims[0]/2+2) && (x >= dims[1]/2-2 && x <= 1+dims[1]/2+2) {
				fmt.Println(y, x)
				maskMatrix.SetValue(complex(0, 0), []int{y, x})
			} else {
				maskMatrix.SetValue(complex(float64(math.MaxUint8), 0), []int{y, x})
			}
		}
	}

	maskedFourier, err := fourier.ApplyMatrixMask(maskMatrix)
	if err != nil {
		t.Fatal(err)
	}

	testSave(*maskedFourier.AmplitudeImage(), "masked-amplitude")
	testSave(*maskedFourier.BrighterAmplitudeImage(), "masked-brighter-amplitude")
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
