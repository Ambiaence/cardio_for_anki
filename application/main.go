package main

import (
	"strconv"
	"fmt"
	"log"
	"os"
	"strings"
	"image/color"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"georgeallen.net/audio_cards/controller"
)

type (
	D = layout.Dimensions
	C = layout.Context
)

var (
	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}

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
)

var words = []string{"One","two","three"}  

func main() {
	controller.GenerateSpokenWord("The word", "English")
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(window *app.Window) error {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))


	var ops op.Ops

	for {
		switch event := window.NextEvent().(type)  {
		case app.DestroyEvent:
			return event.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, event)
			kitchen(gtx, theme)
			event.Frame(gtx.Ops)
		}
	}
}

func kitchen(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for {
		e, ok := lineEditor.Update(gtx)
		if !ok {
			break
		}
		if e, ok := e.(widget.SubmitEvent); ok {
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
