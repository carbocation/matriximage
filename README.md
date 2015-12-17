# matriximage

Golang utility to convert between (grayscale) images and their Fourier transforms, and to manipulate them in frequency space

```sh
go get github.com/carbocation/matriximage
```

##Example:
```go
package main

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/mjibson/go-dsp/dsputils"
)

// Work with the sine wave image
const Case = "sin"
const MaxUint = math.MaxUint16

func main() {
	Basic()
	Masked()
}

// Reads an image, performs a Fourier transform, then outputs various components
// as visible images
func Basic() {
	// Read an image from file. This helper function also converts it to Gray16
	m, err := FromFile(fmt.Sprintf("./images/%s.gif", Case))
	if err != nil {
		t.Fatal(err)
	}

	// Take the 2D Fourier transform of the image
	fourier := m.DFT()

	// Ensure the output folder exists
	os.MkdirAll("./images.out", 0755)

	// Write the amplitude, (brightened) amplitude, phase, and then recombined gray image to disk
	save(*fourier.AmplitudeImage(), "amp")
	save(*fourier.BrighterAmplitudeImage(), "logamp")
	save(*fourier.PhaseImage(), "phase")
	save(*fourier.IDFTImage(), "recovered")
}

// Reads an image, applies a mask in the frequency space, performs the inverse 
// Fourier transform, and writes the resulting recovered image to disk
func Masked() {
	// Read an image from file. This helper function also converts it to Gray16
	m, err := FromFile(fmt.Sprintf("./images/%s.gif", Case))
	if err != nil {
		t.Fatal(err)
	}

	// Take the 2D Fourier transform of the image
	fourier := m.DFT()

	// Create a square mask, where only the middle part of the amplitude is recovered.
	maskMatrix := dsputils.MakeEmptyMatrix(fourier.Dimensions()) //fourier.AmplitudeImage().toGrayMatrix()
	dims := maskMatrix.Dimensions()
	for y := 0; y < dims[0]; y++ {
		for x := 0; x < dims[1]; x++ {
			rad := 1
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

	save(*maskedFourier.AmplitudeImage(), "masked-amplitude")
	save(*maskedFourier.BrighterAmplitudeImage(), "masked-brighter-amplitude")
	save(*maskedFourier.PhaseImage(), "masked-phase")
	save(*maskedFourier.IDFTImage(), "masked-recovered")
}

func save(m Image, name string) {
	os.MkdirAll("./images.out", 0755)

	filename := fmt.Sprintf("./images.out/%s-%s.png", Case, name)

	m.ToFile(filename)
}
```