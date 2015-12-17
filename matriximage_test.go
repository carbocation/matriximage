package matriximage

import (
	"fmt"
	"os"
	"testing"

	//"github.com/mjibson/go-dsp/dsputils"
	//"github.com/mjibson/go-dsp/fft"
)

const Case = "sin"

func TestRecombine(t *testing.T) {
	m, err := FromFile(fmt.Sprintf("./images/%s.gif", Case))
	if err != nil {
		t.Fatal(err)
	}

	matrix := m.toGrayMatrix()
	transformedImg := m.FFT()

	recoveredImg, err := transformedImg.IFFTN()
	if err != nil {
		t.Fatal(err)
	}

	if close := recoveredImg.toGrayMatrix().PrettyClose(matrix); !close {
		t.Errorf("Recovered matrix is not the same as the original matrix. First values: %v %v", recoveredImg.toGrayMatrix().Value([]int{0, 0}), matrix.Value([]int{0, 0}))
	}
}

func TestSplitMatrix(t *testing.T) {
	m, err := FromFile(fmt.Sprintf("./images/%s.gif", Case))
	if err != nil {
		t.Fatal(err)
	}

	transformedImg := m.FFT()
	transformed := transformedImg.matrix
	amp, phase := splitAmpPhase(transformed)
	recombined := mergeAmpPhase(amp, phase)

	if close := transformed.PrettyClose(recombined); !close {
		t.Fatalf("Splitting and recombining the FFT data failed")
	}
}

func TestFrequencyImage(t *testing.T) {
	m, err := FromFile(fmt.Sprintf("./images/%s.gif", Case))
	if err != nil {
		t.Fatal(err)
	}

	fi := m.FFT()

	testSave(*fi.Amp, "amp")
	testSave(*fi.Phase, "phase")

	recovered, err := fi.IFFTN()
	if err != nil {
		t.Fatal(err)
	}

	testSave(*recovered, "recovered")

	/*
		fmt.Println("m", m.Image)
		fmt.Println("m.real", toImage(m.toGrayMatrix(), Real).Image)

		fmt.Println("fft.Amp", toImage(fi.matrix, Amplitude).Image)
		fmt.Println("fft.LogAmp", toImage(fi.matrix, LogAmplitude).Image)
		fmt.Println("fft.Amp", fi.Amp.Image)
		fmt.Println("fft.Phase", fi.Phase.Image)

		fmt.Println("recovered", recovered.Image)
	*/

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
