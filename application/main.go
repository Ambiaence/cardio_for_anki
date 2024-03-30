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
	"gioui.org/io/key"
	_"gioui.org/io/system"
	"georgeallen.net/audio_cards/controller"
)

type mode struct {
	value int
}

func (m *mode) next_mode (){
	fmt.Println("-----")
	fmt.Println(m.value)
	m.value = m.value + 1
	fmt.Println(m.value)

	if m.value > 2 {
		m.value = 0
	}
}

type (
	D = layout.Dimensions
	C = layout.Context
)

type WordButton struct {
	text string
	word string
	area *widget.Clickable
	number int	
}

var (
	update_words = true
	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}

	word_buttons = []WordButton{}

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
var sentence string 
var control_state = mode{value: 0}

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
			//ev, _ := event.Source.Event(key.Filter{}) 
			//print(ev)
			gtx := app.NewContext(&ops, event)
			dashboard(gtx, theme)
			event.Frame(gtx.Ops)
		}
	}
}

func dashboard(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for {
		e, ok := lineEditor.Update(gtx)

		if !ok {
			break
		}

		if e, ok := e.(widget.SubmitEvent); ok {
			sentence = string(e.Text)
			words = strings.Fields(e.Text)
			lineEditor.SetText("")
			for _, word := range words {
				fmt.Print(word, "\n")
			}
			update_words = true;
		}
	}

	for {
		key_event, ok := gtx.Source.Event(key.Filter{})
		if !ok {
			break
		}

		key, ok := key_event.(key.Event)
		if !ok {
			continue
		}
		
		if key.Name == "Tab" {
			control_state.next_mode()
		}

		switch state := control_state.value; state {
			case 0: 
				fmt.Println("One")
			case 1:
				fmt.Println("Two")
			case 2:
				fmt.Println("Thr")
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

		buttons_new := func(gtx C, i int) D {
			var button WordButton
			text := strconv.Itoa(i) 
			text += " - "
			text += words[i]

			button.word = words[i] 
			button.area = new(widget.Clickable)
			button.text = text
			button.number = i

			word_buttons = append(word_buttons, button)
			update_words = false;
			return material.Button(th, button.area, button.text).Layout(gtx)
		}

		buttons_old := func(gtx C, i int) D {
			button := word_buttons[i]
			return material.Button(th, button.area, button.text).Layout(gtx)
		}

		button_generator :=  buttons_old

		if (update_words) {
			word_buttons = word_buttons[:0]
			button_generator = buttons_new
		}

		update_words = false
		anon_list := func(gtx C) D {
			return list_style.Layout(gtx, len(words), button_generator)
		}
		return flex.Layout(gtx, layout.Flexed(1, anon_list))
	}

	stage_one := []layout.Widget{
		buttons,
	}

	stage_two := []layout.Widget{
		editor,
	}

	widgets := &stage_one

	if (control_state.value == 1) {
		widgets = &stage_one
	} else {
		widgets = &stage_two
	}


	spaced := func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, (*widgets)[i])
	}

	list_style := material.List(th, list)
	dimensions := list_style.Layout(gtx, len((*widgets)), spaced)
	return dimensions
}
