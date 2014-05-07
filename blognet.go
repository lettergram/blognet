package main

import (
"image"
"image/color"
"image/draw"
"image/png"
"code.google.com/p/draw2d/draw2d"

"log"
"os"
"os/exec"
)

/** Sets Colors for easy use **/
var (
	white color.Color = color.RGBA{255, 255, 255, 255}
	black color.Color = color.RGBA{0, 0, 0, 255}
	teal  color.Color = color.RGBA{0, 200, 200, 255}
	red   color.Color = color.RGBA{200, 30, 30, 255}
	green color.Color = color.RGBA{0, 255, 0, 255}

	max_x = 1024
	max_y = 1024
)

/** Interface for displaying Shapes **/
type displayShape interface {
	drawShape() *image.RGBA
}

/** Struct for rectangle **/
type Rectangle struct {
	p image.Point
	length, width int
}

/** Struct for circle **/
type Circle struct {
	p image.Point
	r int
}

/**
 * The following is taken from:
 * http://blog.golang.org/go-imagedraw-package
 * Excellent use of masking!!!
 */
func (c *Circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *Circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *Circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}


/**
 * Rectangle Draw shape function
 */
func (r Rectangle) drawShape() *image.RGBA{
	return image.NewRGBA(image.Rect(0, 0, r.length, r.width))
}

/**
 * Draws lines between a source and destination,
 * this is a bit ugly, but should be fairly straight forward
 */
func drawlines(src_x float64, src_y float64, dest_x float64, dest_y float64, m *image.RGBA){
	gc := draw2d.NewGraphicContext(m)
	gc.MoveTo(src_x, src_y)
	gc.LineTo(dest_x, dest_y)
	gc.SetStrokeColor(red)
	gc.SetLineWidth(3)
	gc.Stroke()
}

func buildMap() *image.RGBA{

	surface := Rectangle{ length:max_y, width:max_x } // Surface to draw on
	painter := Rectangle{ length:max_y, width:max_x } // Colored Mask Layer

	m := surface.drawShape()
	cc  := painter.drawShape()

	draw.Draw(m, m.Bounds(), &image.Uniform{teal}, image.ZP, draw.Src)
	draw.Draw(cc, cc.Bounds(), &image.Uniform{red}, image.ZP, draw.Src)

	/** Generates two circles **/
	draw.DrawMask(m, m.Bounds(), cc, image.ZP, &Circle{image.Point{640, 640}, 20}, image.ZP, draw.Over)
	draw.DrawMask(m, m.Bounds(), cc, image.ZP, &Circle{image.Point{520, 520}, 35}, image.ZP, draw.Over)
	draw.DrawMask(m, m.Bounds(), cc, image.ZP, &Circle{image.Point{407, 180}, 55}, image.ZP, draw.Over)

	/** Draws Line Between Circles **/
	drawlines(640, 640, 520, 520, m)
	drawlines(520, 520, 407, 180, m)

	return m
}

// show  a specified file by Preview.app for OS X(darwin)
func Show(name string) {
	command := "open"
	arg1 := "-a"
	arg2 := "/Applications/Preview.app"
	cmd := exec.Command(command, arg1, arg2, name)
	err := cmd.Run()
	
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	m := buildMap();

	w, _ := os.Create("blogmap.png")
	defer w.Close()
	png.Encode(w, m) //Encode writes the Image m to w in PNG format.

	Show(w.Name())
}
