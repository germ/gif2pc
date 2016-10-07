package main

import (
	"fmt"
	"image/gif"
	"image/color"
	"io/ioutil"
	"os"
	"text/template"
)

type Point struct {
	X	float64
	Y	float64
	Z	float64
	R	float64
	G	float64
	B	float64
}

type TmplFill struct {
	R	int
	G	int
	B	int
}
type Cloud map[color.Color][]Point

func main() {
	if len(os.Args) == 1 {
		panic("No input")
	}
	f, _ := os.Open(os.Args[1])
	defer f.Close()

	src, err := gif.DecodeAll(f)
	if err != nil {
		panic(err)
	}

	maxDepth := len(src.Delay)

	cloud := make(Cloud, 10)
	var depth int
	for gifIndex, gifVal := range(src.Image) {
		depth += src.Delay[gifIndex]
		// Remove inner solids
		for x := 0; x < src.Image[gifIndex].Rect.Max.X; x++ {
			// Extract color data
			for y := 0; y < src.Image[gifIndex].Rect.Max.Y; y++ {
				cen := gifVal.At(x,y)
				
				// Be sure not edge, remove inner surface
				if !(x == 0 || x == src.Image[gifIndex].Rect.Max.X || y == 0 || y == src.Image[gifIndex].Rect.Max.Y) {
					// Calc neighbor delta
					left := gifVal.At(x-1,y)
					right := gifVal.At(x+1,y)
					up := gifVal.At(x,y-1)
					down := gifVal.At(x-1,y+1)
					if (cen == left && cen == right && cen == up && cen == down && gifIndex != 0 && gifIndex != len(src.Image)-1) {
						continue
					}
				}

				//Holy crap an edge
				r,g,b, _ := cen.RGBA()
				pt := Point {
					X:	float64(x)/float64(maxDepth),
					Y:	float64(y)/float64(maxDepth),
					Z:	float64(gifIndex)/float64(len(src.Delay)),
					//Z:	(float64(depth) + ((float64(lum)/255.0)*float64(src.Delay[gifIndex])))/float64(maxDepth),
					R:	float64(r)/65535,
					G:	float64(g)/65535,
					B:	float64(b)/65535,
				}
				cloud[cen] = append(cloud[cen], pt)
			}
		}
	}

	cloud.Write(os.Args[1])
}

func (c Cloud)Normalize(depth int) {
}

func (c Cloud)Write(fname string) {
	// Set up filler
	filler, _ := template.New("script").Parse(ScriptTmpl)


	// Loop over every color channel
	chanNum := 0
	for _, colorMap := range c {
		// Generate point map
		var data []byte
		var colorPt Point
		for _, v := range colorMap {
			line := fmt.Sprintf("%f;%f;%f\n", v.X, v.Y, v.Z)
			data = append(data, []byte(line)...)
			colorPt = v
		}
		
		name := fmt.Sprintf("%v.txt", chanNum)
		ioutil.WriteFile(name, data, os.ModePerm)

		// Generate script
		colorData := TmplFill{
			R:	int(colorPt.R*255),
			G:	int(colorPt.G*255),
			B:	int(colorPt.B*255),
		}
		f, _ := os.Create(fmt.Sprintf("%v.mlx", chanNum))
		filler.Execute(f, colorData)
		f.Close()

		chanNum += 1
	}

	ioutil.WriteFile("shape.mlx", []byte(ShapeTmpl), os.ModePerm)
}

var ScriptTmpl = `
<!DOCTYPE FilterScript>
<FilterScript>
 <filter name="Per Face Color Function">
  <Param tooltip="function to generate Red component. Expected Range 0-255" description="func r = " type="RichString" value="{{.R}}" name="r"/>
  <Param tooltip="function to generate Green component. Expected Range 0-255" description="func g = " type="RichString" value="{{.G}}" name="g"/>
  <Param tooltip="function to generate Blue component. Expected Range 0-255" description="func b = " type="RichString" value="{{.B}}" name="b"/>
  <Param tooltip="function to generate Alpha component. Expected Range 0-255" description="func alpha = " type="RichString" value="255" name="a"/>
 </filter>
</FilterScript>
`
var ShapeTmpl = `
<!DOCTYPE FilterScript>
<FilterScript>
 <filter name="Alpha Complex/Shape">
  <Param tooltip="Compute the alpha value as percentage of the diagonal of the bbox" description="Alpha value" type="RichAbsPerc" value="0.058043" min="0" name="alpha" max="5.80434"/>
  <Param tooltip="Select the output. The Alpha Shape is the boundary of the Alpha Complex" description="Get:" enum_val0="Alpha Complex" enum_val1="Alpha Shape" enum_cardinality="2" type="RichEnum" value="0" name="Filtering"/>
 </filter>
</FilterScript>`
 
