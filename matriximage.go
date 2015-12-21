package matriximage

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"

	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
)

const MaxUint = math.MaxUint16

type Image struct {
	image.Image
}

func (m Image) DFT() FourierImage {
	return FourierImage{Matrix: m.fftn()}
}

// Work with gray for now
// Returns a matrix without rescaling values
func (m Image) ToGrayMatrix() *dsputils.Matrix {
	var img *image.Gray16
	switch t := m.Image.(type) {
	case *image.Gray16:
		img = t
	default:
		img = ImageToGray(m)
	}

	// Generate 0-based dimensions
	min, max := img.Bounds().Min, img.Bounds().Max
	lenY, lenX := max.Y-min.Y, max.X-min.X

	matrix := dsputils.MakeEmptyMatrix([]int{lenY, lenX})

	scale := 1.0

	for i := 0; i < lenX; i++ {
		for j := 0; j < lenY; j++ {

			v := scale * float64(img.Gray16At(i+min.X, j+min.Y).Y)

			matrix.SetValue(complex(v, 0), []int{j, i})
		}
	}

	return matrix
}

func (m Image) fftn() *dsputils.Matrix {
	matrix := m.ToGrayMatrix()
	return fft.FFTN(matrix)
}

func FromFile(filename string) (*Image, error) {
	infile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, err := image.Decode(infile)
	if err != nil {
		return nil, err
	}

	grayImage := ImageToGray(src)

	return &Image{Image: grayImage}, nil
}

func ImageToGray(m image.Image) *image.Gray16 {
	b := m.Bounds()
	gray := image.NewGray16(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			gray.SetGray16(x, y, color.Gray16Model.Convert(m.At(x, y)).(color.Gray16))
		}
	}
	return gray
}

func (m Image) ToFile(named string) error {
	outfile, err := os.Create(named)
	if err != nil {
		return err
	}
	defer outfile.Close()

	return png.Encode(outfile, m.Image)
}
