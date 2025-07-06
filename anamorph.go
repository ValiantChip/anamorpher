package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"os"

	"github.com/ValiantChip/anamorpher/lib/anamorpher"
)

func ReturnWithCode() int {
	var radius int
	var degrees int
	var inputPath string
	var outputpath string
	var outheight int
	var outwidth int
	var scale float64
	var help bool

	flag.CommandLine.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "anamorph [input path] [options]\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}

	flag.IntVar(&radius, "r", 100, "radius")
	flag.IntVar(&degrees, "d", 45, "degrees")
	flag.StringVar(&outputpath, "o", "./out.jpg", "output path")
	flag.IntVar(&outheight, "height", 0, "output height")
	flag.IntVar(&outwidth, "w", 0, "output width")
	flag.Float64Var(&scale, "s", 1.0, "scale")
	flag.BoolVar(&help, "h", false, "\nshow this message")
	flag.Parse()

	if help {
		flag.CommandLine.Usage()
		return 0
	}

	inputPath = flag.Arg(0)
	if inputPath == "" {
		fmt.Printf("Error: input path must be set\n")
		return 2
	}

	if radius < 1 {
		fmt.Printf("Error: radius must be set to a value greater than 0\n")
		return 2
	}

	if scale < 1.0 {
		fmt.Printf("Error: scale must not be less than 1\n")
		return 2
	}

	fl, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Error: invalid input path\n")
		return 2
	}
	defer fl.Close()

	out, err := os.Create(outputpath)
	if err != nil {
		fmt.Printf("Error: error creating output file\n")
		return 1
	}

	defer out.Close()

	img, err := jpeg.Decode(fl)
	if err != nil {
		fmt.Printf("Error: error decoding image: %s\n", err.Error())
		return 1
	}

	radius = int(float64(radius) * scale)

	mod := image.NewNRGBA(image.Rect(0, 0, int(float64(outwidth)*scale), int(float64(outheight)*scale)))

	morph := anamorpher.New(img, mod, radians(float64(degrees)), float64(radius))

	maxBounds := morph.MaximumRequiredBounds()
	maxBounds.Max.X = int(float64(maxBounds.Max.X))
	maxBounds.Max.Y = int(float64(maxBounds.Max.Y))
	if outheight == 0 || outwidth == 0 {
		nMod := image.NewNRGBA(maxBounds)

		morph.Mod = nMod
	}

	err = morph.Anamorph()
	if err != nil {
		fmt.Printf("Error: error anamorphing image: %s\n", err.Error())
		return 1
	}

	err = jpeg.Encode(out, morph.Mod, nil)
	if err != nil {
		fmt.Printf("Error: error encoding image: %s\n", err.Error())
		return 1
	}

	fmt.Printf("image successfully anamorphed\n")
	return 0
}

func main() {
	code := ReturnWithCode()
	if code == 2 {
		flag.CommandLine.Usage()
	}

	os.Exit(code)
}

func radians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}
