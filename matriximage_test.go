package matriximage

import (
	"fmt"
	"os"
	"testing"

	//"github.com/mjibson/go-dsp/dsputils"
)

const Case = "4x2"

func TestFrequencyImage(t *testing.T) {
	m, err := FromFile(fmt.Sprintf("./images/%s.gif", Case))
	if err != nil {
		t.Fatal(err)
	}

	fi := m.FFT()

	testSave(*fi.Real, "real")
	testSave(*fi.Imag, "imag")

	recovered, err := fi.IFFTN()
	if err != nil {
		t.Fatal(err)
	}

	testSave(*recovered, "recovered")

	/*
		array := [][]complex128{
			dsputils.ToComplex([]float64{1, 2, 3, 9}),
			dsputils.ToComplex([]float64{8, 5, 1, 2}),
			dsputils.ToComplex([]float64{9, 8, 7, 2}),
		}

		fmt.Println(dsputils.MakeMatrix2(array))
		fmt.Println(dsputils.MakeEmptyMatrix([]int{4, 2}))
	*/

	fmt.Println(m.Image)
	fmt.Println(toImage(m.toGrayMatrix(), false).Image)

	fftmatrix := m.fftn()
	fmt.Println(fftmatrix)

	fmt.Println(fi.Real.Image)
	fmt.Println(fi.Imag.Image)

	fmt.Println(recovered.Image)
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
