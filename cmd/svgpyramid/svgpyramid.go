// svgpyramid draw a pyramid SVG infographic.
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/ajstarks/svgo"
)

func choseStyle(color, highlightColor string, highlightText, num int) string {
	if num == highlightText {
		return fmt.Sprintf("text-anchor:middle;font-size:25px;fill:%s", highlightColor)
	}
	return fmt.Sprintf("text-anchor:middle;font-size:25px;fill:%s", color)
}

func svgPyramid(color, highlightColor string, highlightText int) {
	width := 500
	height := 500
	xOffset := 50
	yOffset := 50

	sideLength := width
	halfSideLength := sideLength / 2
	triangleHeight := int(math.Sqrt(float64(sideLength*sideLength - halfSideLength*halfSideLength)))

	canvas := svg.New(os.Stdout)
	canvas.Start(width+xOffset+xOffset, height+yOffset+yOffset)
	s := fmt.Sprintf("stroke:%s;stroke-width:3", color)
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

	// 1.
	style := choseStyle(color, highlightColor, highlightText, 1)
	canvas.Text(halfSideLength+xOffset, int(6.125*h)+yOffset, "secure devices", style)

	// 2.
	style = choseStyle(color, highlightColor, highlightText, 2)
	canvas.Text(halfSideLength+xOffset, int(5.125*h)+yOffset, "secure software", style)

	// 3.
	style = choseStyle(color, highlightColor, highlightText, 3)
	canvas.Text(int(4.75*w)+xOffset, int(3.75*h)+yOffset, "anon.", style)
	canvas.Text(int(4.5*w)+xOffset, int(4.25*h)+yOffset, "messaging", style)

	// 4.
	style = choseStyle(color, highlightColor, highlightText, 4)
	canvas.Text(int((6.5+1.75)*w)+xOffset, int(3.75*h)+yOffset, "digital", style)
	canvas.Text(int((6.5+1.75)*w)+xOffset, int(4.25*h)+yOffset, "cash", style)

	// 5.
	style = choseStyle(color, highlightColor, highlightText, 5)
	canvas.Text(int(5.25*w)+xOffset, int(2.75*h)+yOffset, "nyms", style)

	// 6.
	style = choseStyle(color, highlightColor, highlightText, 6)
	canvas.Text(int(7.75*w)+xOffset, int(2.75*h)+yOffset, "DNMs", style)

	// 7.
	style = choseStyle(color, highlightColor, highlightText, 7)
	canvas.Text(halfSideLength+xOffset, int(1.25*h)+yOffset, "I/F", style)

	canvas.End()
}

func main() {
	color := flag.String("color", "black", "set pyramid color")
	highlightColor := flag.String("highlight-color", "white", "set pyramid highlighting color")
	highlightText := flag.Int("highlight", 0, "highlight text element #")
	flag.Parse()
	svgPyramid(*color, *highlightColor, *highlightText)
}
