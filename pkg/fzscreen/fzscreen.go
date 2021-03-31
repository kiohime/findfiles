package fzscreen

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/junegunn/fzf/src/algo"
	"github.com/junegunn/fzf/src/util"
)

func findPos(input string, pattern string) *[]int {

	caseSensitive := false
	forward := true
	normalize := forward
	chars := util.ToChars([]byte(input))

	_, pos := algo.FuzzyMatchV2(caseSensitive, normalize, forward, &chars, []rune(pattern), true, nil)
	return pos
}

func SetColor(s string) color.RGBA {
	var c color.RGBA
	switch s {
	case "red":
		c = color.RGBA{255, 0, 0, 255}
	case "green":
		c = color.RGBA{0, 255, 0, 255}
	case "blue":
		c = color.RGBA{0, 0, 255, 255}
	case "black":
		c = color.RGBA{0, 0, 0, 255}
	case "white":
		c = color.RGBA{255, 255, 255, 255}
	}

	return c
}

func nextPos(pi int, pos *[]int) (int, int) {
	pi--
	if pi < 0 {
		return pi, int(^uint(0) >> 1)
	}
	return pi, (*pos)[pi]
}

func sliceStr(s string, pos *[]int) []string {
	if pos == nil {
		return nil
	}
	pi := len(*pos) - 1
	if pi < 0 {
		return []string{s}
	}
	p := (*pos)[pi]
	hl := false
	segments := []string{}
	seg := ""
	for i, r := range s {
		if !hl {
			if i != p {
				seg += string(r)
			} else {
				segments = append(segments, seg)
				hl = true
				seg = string(r)
				pi, p = nextPos(pi, pos)
			}
		} else {
			if i == p {
				seg += string(r)
				pi, p = nextPos(pi, pos)
			} else {
				segments = append(segments, seg)
				hl = false
				seg = string(r)
			}
		}
	}
	if seg != "" {
		segments = append(segments, seg)
	}
	return segments
}

func getHighlightColor(hl bool) color.RGBA {
	if hl {
		return SetColor("red")
	}
	return SetColor("black")
}

func NewFyneString(s string, pos *[]int) *FyneString {
	ct := []*canvas.Text{}

	slice := sliceStr(s, pos)
	fmt.Println(*pos)
	fmt.Println(s)
	colored := false
	for _, s := range slice {
		if s == "" {
			colored = !colored
			continue
		}
		for _, r := range s {
			r := string(r)
			ct = append(ct, canvas.NewText(r, getHighlightColor(colored)))

		}
		if colored {
			s = strings.ToUpper(s)
		}
		fmt.Print(s)
		colored = !colored
	}
	fmt.Println("")

	newString := FyneString{segments: ct}
	return &newString
}

type FyneString struct {
	segments []*canvas.Text
	// pos      *[]int
}

func (o *FyneString) renderString(rPos fyne.Position) fyne.CanvasObject {
	const kerning = -1

	res := fyne.NewContainer()
	posX := rPos.X
	posY := rPos.Y
	var nextOff fyne.Size
	for _, segment := range o.segments {
		off := fyne.MeasureText(segment.Text, segment.TextSize, segment.TextStyle)
		posX += nextOff.Width + kerning
		segment.Move(fyne.NewPos(posX, posY))
		res.Objects = append(res.Objects, segment)
		nextOff = off
	}
	return res
}

func Render(inData []string, pat string) fyne.CanvasObject {
	const leading = -3

	result := fyne.NewContainer()
	resX := result.Position().X
	resY := result.Position().Y
	for _, in := range inData {
		strPos := findPos(in, pat)
		if strPos == nil {
			continue
		}
		// fyne.MeasureText(s)
		strData := NewFyneString(in, strPos)
		sample := strData.segments[0]
		off := fyne.MeasureText(in, sample.TextSize, sample.TextStyle)
		result.Objects = append(result.Objects, strData.renderString(result.Position()))
		resY = off.Height + resY + leading

		result.Move(fyne.NewPos(resX, resY))
	}
	return result
}
