package main

import (
	"strconv"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"time"
	"strings"

	"gioui.org/app"
	//"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/gpu/headless"
	"gioui.org/io/event"
	//"gioui.org/io/input"
	"gioui.org/layout"
	"gioui.org/op"
	//"gioui.org/op/clip"
	//"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	screenshot = flag.String("screenshot", "", "save a screenshot to a file and exit")
	disable    = flag.Bool("disable", false, "disable all widgets")
)

var words = []string{"Test 1", "Test 2", "Test 3"}

type iconAndTextButton struct {
	theme  *material.Theme
	button *widget.Clickable
	icon   *widget.Icon
	word   string
}

func main() {
	flag.Parse()
	ic, err := widget.NewIcon(icons.ContentAdd)
	if err != nil {
		log.Fatal(err)
	}
	icon = ic
	progressIncrementer = make(chan float32)
	if *screenshot != "" {
		if err := saveScreenshot(*screenshot); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save screenshot: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	go func() {
		for {
			time.Sleep(time.Second)
			progressIncrementer <- 0.1
		}
	}()

	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func saveScreenshot(f string) error {
	const scale = 1.5
	sz := image.Point{X: 800 * scale, Y: 600 * scale}
	w, err := headless.NewWindow(sz.X, sz.Y)
	if err != nil {
		return err
	}
	gtx := layout.Context{
		Ops: new(op.Ops),
		Metric: unit.Metric{
			PxPerDp: scale,
			PxPerSp: scale,
		},
		Constraints: layout.Exact(sz),
	}
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	kitchen(gtx, th)
	w.Frame(gtx.Ops)
	img := image.NewRGBA(image.Rectangle{Max: sz})
	err = w.Screenshot(img)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}
	return ioutil.WriteFile(f, buf.Bytes(), 0o666)
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	events := make(chan event.Event)
	acks := make(chan struct{})

	go func() {
		for {
			ev := w.NextEvent()
			events <- ev
			<-acks
			if _, ok := ev.(app.DestroyEvent); ok {
				return
			}
		}
	}()

	var ops op.Ops
	for {
		select {
		case e := <-events:
			switch e := e.(type) {
			case app.DestroyEvent:
				acks <- struct{}{}
				return e.Err
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				if *disable {
					gtx = gtx.Disabled()
				}
				if checkbox.Update(gtx) {
					if checkbox.Value {
						transformTime = e.Now
					} else {
						transformTime = time.Time{}
					}
				}
				kitchen(gtx, th)
				e.Frame(gtx.Ops)
			}
			acks <- struct{}{}
		case p := <-progressIncrementer:
			progress += p
			if progress > 1 {
				progress = 0
			}
			w.Invalidate()
		}
	}
}

var (
	editor     = new(widget.Editor)
	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	button            = new(widget.Clickable)
	greenButton       = new(widget.Clickable)
	iconTextButton    = new(widget.Clickable)
	iconButton        = new(widget.Clickable)
	flatBtn           = new(widget.Clickable)
	disableBtn        = new(widget.Clickable)
	radioButtonsGroup = new(widget.Enum)

	button_list = &widget.List{
		List: layout.List{
			Axis: layout.Horizontal,
		},

	}

	list              = &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
	progress            = float32(0)
	progressIncrementer chan float32
	green               = true
	topLabel            = "Hello, Gio"
	topLabelState       = new(widget.Selectable)
	icon                *widget.Icon
	checkbox            = new(widget.Bool)
	swtch               = new(widget.Bool)
	transformTime       time.Time
	float               = new(widget.Float)
)

type (
	D = layout.Dimensions
	C = layout.Context
)

func (b iconAndTextButton) Layout(gtx layout.Context) layout.Dimensions {
	return material.ButtonLayout(b.theme, b.button).Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			textIconSpacer := unit.Dp(5)

			layIcon := layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: textIconSpacer}.Layout(gtx, func(gtx C) D {
					var d D
					if b.icon != nil {
						size := gtx.Dp(unit.Dp(56)) - 2*gtx.Dp(unit.Dp(16))
						gtx.Constraints = layout.Exact(image.Pt(size, size))
						d = b.icon.Layout(gtx, b.theme.ContrastFg)
					}
					return d
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					l := material.Body1(b.theme, b.word)
					l.Color = b.theme.Palette.ContrastFg
					return l.Layout(gtx)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel)
		})
	})
}

func kitchen(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for {
		e, ok := lineEditor.Update(gtx)
		if !ok {
			break
		}
		if e, ok := e.(widget.SubmitEvent); ok {
			topLabel = e.Text
			words = strings.Fields(e.Text)
			lineEditor.SetText("")
			for _, word := range words {
				fmt.Print(word, "\n")
			}
		}
	}

	editor := func(gtx C) D {
		editor_style := material.Editor(th, lineEditor, "Hint")
		editor_style.Font.Style = font.Italic
		border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
		spaced := func(gtx C) D {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, editor_style.Layout)
		}
		dimensions := border.Layout(gtx, spaced)
		return dimensions 
	}

	buttons := func(gtx C) D {
		flex := layout.Flex{}
		list_style := material.List(th, button_list)
		button_generator := func(gtx C, i int) D {
			fmt.Println("Test?")
			button := new(widget.Clickable)
			text := strconv.Itoa(i) 
			text += ": "
			text += words[i]
			return material.Button(th, button, text).Layout(gtx)
		}

		anon_list := func(gtx C) D {
			return list_style.Layout(gtx, len(words), button_generator)
		}
		return flex.Layout(gtx, layout.Flexed(1, anon_list))
	}


	widgets := []layout.Widget{
		editor,
		buttons,
	}

	spaced := func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
	}

	list_style := material.List(th, list)
	dimensions := list_style.Layout(gtx, len(widgets), spaced)
	return dimensions
}
