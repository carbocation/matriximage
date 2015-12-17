package matriximage

import (
	"fmt"
	"os"
	"testing"
)

const Case = "churchill"

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
