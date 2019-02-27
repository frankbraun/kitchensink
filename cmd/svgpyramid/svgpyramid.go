// svgpyramid draw a pyramid SVG infographic.
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	svg "github.com/ajstarks/svgo/float"
)

func choseLineStyle(color, highlightColor string, highlightText int, nums ...int) string {
	for _, num := range nums {
		if num == highlightText {
			return fmt.Sprintf("stroke:%s;stroke-width:3;stroke-linecap:round", highlightColor)
		}
	}
	return fmt.Sprintf("stroke:%s;stroke-width:3;stroke-linecap:round", color)
}

func choseTextStyle(color, highlightColor string, highlightText, num int) string {
	if num == highlightText {
		return fmt.Sprintf("text-anchor:middle;font-size:25px;fill:%s", highlightColor)
	}
	return fmt.Sprintf("text-anchor:middle;font-size:25px;fill:%s", color)
}

func svgPyramid(color, highlightColor string, highlightText int) {
	width := 500.0
	height := 500.0
	xOffset := 50.0
	yOffset := 50.0

	sideLength := width
	halfSideLength := sideLength / 2
	triangleHeight := math.Sqrt(sideLength*sideLength - halfSideLength*halfSideLength)

	canvas := svg.New(os.Stdout)
	canvas.Start(width+xOffset+xOffset, height+yOffset+yOffset)

	h := math.Sqrt(sideLength*sideLength-halfSideLength*halfSideLength) / 6.5
	w := halfSideLength / 6.5

	// left triangle line
	s := choseLineStyle(color, highlightColor, highlightText, 7)
	canvas.Line(6.5*w+xOffset, 0+yOffset, 5.0*w+xOffset, 1.5*h+yOffset, s)

	s = choseLineStyle(color, highlightColor, highlightText, 5)
	canvas.Line(5.0*w+xOffset, 1.5*h+yOffset, 3.5*w+xOffset, 3.0*h+yOffset, s)

	s = choseLineStyle(color, highlightColor, highlightText, 3)
	canvas.Line(3.5*w+xOffset, 3.0*h+yOffset, 2.0*w+xOffset, 4.5*h+yOffset, s)

	s = choseLineStyle(color, highlightColor, highlightText, 2)
	canvas.Line(2.0*w+xOffset, 4.5*h+yOffset, 1.0*w+xOffset, 5.5*h+yOffset, s)

	s = choseLineStyle(color, highlightColor, highlightText, 1)
	canvas.Line(1.0*w+xOffset, 5.5*h+yOffset, 0+xOffset, 6.5*h+yOffset, s)

	// right triangle line
	s = choseLineStyle(color, highlightColor, highlightText, 7)
	canvas.Line(6.5*w+xOffset, 0+yOffset, 8.0*w+xOffset, 1.5*h+yOffset, s)

	s = choseLineStyle(color, highlightColor, highlightText, 6)
	canvas.Line(8.0*w+xOffset, 1.5*h+yOffset, 9.5*w+xOffset, 3.0*h+yOffset, s)

	s = choseLineStyle(color, highlightColor, highlightText, 4)
	canvas.Line(9.5*w+xOffset, 3.0*h+yOffset, 11.0*w+xOffset, 4.5*h+yOffset, s)

	s = choseLineStyle(color, highlightColor, highlightText, 2)
	canvas.Line(11.0*w+xOffset, 4.5*h+yOffset, 12.0*w+xOffset, 5.5*h+yOffset, s)

	s = choseLineStyle(color, highlightColor, highlightText, 1)
	canvas.Line(12.0*w+xOffset, 5.5*h+yOffset, 13.0*w+xOffset, 6.5*h+yOffset, s)

	// bottom triangle line
	s = choseLineStyle(color, highlightColor, highlightText, 1)
	canvas.Line(0+xOffset, triangleHeight+yOffset, sideLength+xOffset, triangleHeight+yOffset, s)

	// first horizontal line
	s = choseLineStyle(color, highlightColor, highlightText, 5, 7)
	canvas.Line(5.0*w+xOffset, 1.5*h+yOffset, 6.5*w+xOffset, 1.5*h+yOffset, s)

	s = choseLineStyle(color, highlightColor, highlightText, 6, 7)
	canvas.Line(6.5*w+xOffset, 1.5*h+yOffset, 8.0*w+xOffset, 1.5*h+yOffset, s)

	// second horizontal line
	s = choseLineStyle(color, highlightColor, highlightText, 3, 5)
	canvas.Line(3.5*w+xOffset, 3.0*h+yOffset, 6.5*w+xOffset, 3.0*h+yOffset, s)

	s = choseLineStyle(color, highlightColor, highlightText, 4, 6)
	canvas.Line(6.5*w+xOffset, 3.0*h+yOffset, 9.5*w+xOffset, 3.0*h+yOffset, s)

	// third horizontal line
	s = choseLineStyle(color, highlightColor, highlightText, 2, 3)
	canvas.Line(2.0*w+xOffset, 4.5*h+yOffset, 6.5*w+xOffset, 4.5*h+yOffset, s)

	s = choseLineStyle(color, highlightColor, highlightText, 2, 4)
	canvas.Line(6.5*w+xOffset, 4.5*h+yOffset, 11.0*w+xOffset, 4.5*h+yOffset, s)

	// fourth horizontal line
	s = choseLineStyle(color, highlightColor, highlightText, 1, 2)
	canvas.Line(1.0*w+xOffset, 5.5*h+yOffset, 12.0*w+xOffset, 5.5*h+yOffset, s)

	// vertical line
	s = choseLineStyle(color, highlightColor, highlightText, 5, 6)
	canvas.Line(halfSideLength+xOffset, 1.5*h+yOffset, halfSideLength+xOffset, 3.0*h+yOffset, s)

	s = choseLineStyle(color, highlightColor, highlightText, 3, 4)
	canvas.Line(halfSideLength+xOffset, 3.0*h+yOffset, halfSideLength+xOffset, 4.5*h+yOffset, s)

	// 1.
	style := choseTextStyle(color, highlightColor, highlightText, 1)
	canvas.Text(halfSideLength+xOffset, 6.125*h+yOffset, "secure devices", style)

	// 2.
	style = choseTextStyle(color, highlightColor, highlightText, 2)
	canvas.Text(halfSideLength+xOffset, 5.125*h+yOffset, "secure software", style)

	// 3.
	style = choseTextStyle(color, highlightColor, highlightText, 3)
	canvas.Text(4.75*w+xOffset, 3.75*h+yOffset, "anon.", style)
	canvas.Text(4.5*w+xOffset, 4.25*h+yOffset, "messaging", style)

	// 4.
	style = choseTextStyle(color, highlightColor, highlightText, 4)
	canvas.Text((6.5+1.75)*w+xOffset, 3.75*h+yOffset, "digital", style)
	canvas.Text((6.5+1.75)*w+xOffset, 4.25*h+yOffset, "cash", style)

	// 5.
	style = choseTextStyle(color, highlightColor, highlightText, 5)
	canvas.Text(5.25*w+xOffset, 2.75*h+yOffset, "nyms", style)

	// 6.
	style = choseTextStyle(color, highlightColor, highlightText, 6)
	canvas.Text(7.75*w+xOffset, 2.75*h+yOffset, "DNMs", style)

	// 7.
	style = choseTextStyle(color, highlightColor, highlightText, 7)
	canvas.Text(halfSideLength+xOffset, 1.25*h+yOffset, "I/F", style)

	canvas.End()
}

func main() {
	color := flag.String("color", "black", "set pyramid color")
	highlightColor := flag.String("highlight-color", "white", "set pyramid highlighting color")
	highlightText := flag.Int("highlight", 0, "highlight text element #")
	flag.Parse()
	svgPyramid(*color, *highlightColor, *highlightText)
}
