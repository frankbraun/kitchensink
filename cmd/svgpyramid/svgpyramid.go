// svgpyramid draw a pyramid SVG infographic.
package main

import (
	"math"
	"os"

	"github.com/ajstarks/svgo"
)

func main() {
	width := 500
	height := 500
	xOffset := 50
	yOffset := 50

	sideLength := width
	halfSideLength := sideLength / 2
	triangleHeight := int(math.Sqrt(float64(sideLength*sideLength - halfSideLength*halfSideLength)))

	canvas := svg.New(os.Stdout)
	canvas.Start(width+xOffset+xOffset, height+yOffset+yOffset)
	s := "stroke:white;stroke-width:3"
	canvas.Polygon([]int{
		0 + xOffset,
		sideLength + xOffset,
		halfSideLength + xOffset,
	},
		[]int{
			triangleHeight + yOffset,
			triangleHeight + yOffset,
			0 + yOffset,
		},
		"fill:none;"+s)

	h := math.Sqrt(float64(sideLength*sideLength-halfSideLength*halfSideLength)) / 6.5
	w := float64(halfSideLength) / 6.5
	canvas.Line(int(5.0*w)+xOffset, int(1.5*h)+yOffset, int(8.0*w)+xOffset, int(1.5*h)+yOffset, s)
	canvas.Line(int(3.5*w)+xOffset, int(3.0*h)+yOffset, int(9.5*w)+xOffset, int(3.0*h)+yOffset, s)
	canvas.Line(int(2.0*w)+xOffset, int(4.5*h)+yOffset, int(11.0*w)+xOffset, int(4.5*h)+yOffset, s)
	canvas.Line(int(1.0*w)+xOffset, int(5.5*h)+yOffset, int(12.0*w)+xOffset, int(5.5*h)+yOffset, s)
	canvas.Line(halfSideLength+xOffset, int(1.5*h)+yOffset, halfSideLength+xOffset, int(4.5*h)+yOffset, s)
	canvas.Text(halfSideLength+xOffset, int(1.25*h)+yOffset, "I/F", "text-anchor:middle;font-size:25px;fill:white")
	canvas.Text(int(5.25*w)+xOffset, int(2.75*h)+yOffset, "nyms", "text-anchor:middle;font-size:25px;fill:white")
	canvas.Text(int(7.75*w)+xOffset, int(2.75*h)+yOffset, "DNMs", "text-anchor:middle;font-size:25px;fill:white")
	canvas.Text(int(4.75*w)+xOffset, int(3.75*h)+yOffset, "anon.", "text-anchor:middle;font-size:25px;fill:white")
	canvas.Text(int(4.5*w)+xOffset, int(4.25*h)+yOffset, "messaging", "text-anchor:middle;font-size:25px;fill:white")
	canvas.Text(int((6.5+1.75)*w)+xOffset, int(3.75*h)+yOffset, "digital", "text-anchor:middle;font-size:25px;fill:white")
	canvas.Text(int((6.5+1.75)*w)+xOffset, int(4.25*h)+yOffset, "cash", "text-anchor:middle;font-size:25px;fill:white")
	canvas.Text(halfSideLength+xOffset, int(5.125*h)+yOffset, "secure software", "text-anchor:middle;font-size:25px;fill:white")
	canvas.Text(halfSideLength+xOffset, int(6.125*h)+yOffset, "secure devices", "text-anchor:middle;font-size:25px;fill:white")
	canvas.End()
}
