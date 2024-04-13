package main
// Return the captured text inside the <a> tag
import (
	"strconv"
	_ "fmt"
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

var red_button_color = color.NRGBA{R: 229, G: 163, B: 163 ,A: 255}
var blue_button_color = color.NRGBA{R: 132, G: 205,  B: 238, A: 255}

const (
	start = iota
	choice = iota
	curate = iota
	end = iota
)

type mode struct {
	value int
}

func (m *mode) next_mode (){
	m.value = m.value + 1
	if m.value > end {
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
	chosen bool
	area *widget.Clickable
	number int	
}


var (
	theme *material.Theme
	red_button_theme *material.Theme

	update_words = true

	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}

	word_buttons = []WordButton{}
	definition_buttons = []WordButton{}

	definition_list = &widget.List{
		List: layout.List{
			Axis: layout.Horizontal,
		},

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
var sentence string 
var control_state = mode{value: 0}

var current_curated_position = 0 
var current_curated_word *WordButton
var reset_curate = true
var curated_word_definitions controller.WordList 
var update_word_definitions = true

func next_chosen_button() *WordButton {
	for {
		current_curated_position = current_curated_position + 1
		if (current_curated_position == len(word_buttons)) {
			current_curated_position = 0
		}

		if word_buttons[current_curated_position].chosen { 
			return &word_buttons[current_curated_position]
		}
	}
}

func main() {
	controller.GenerateSpokenWord("The word", "English")
	controller.TranslateSentence("Dummy")
	//controller.WordEquivalents("bestimmt")

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
	theme = material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	theme.ContrastBg = blue_button_color

	red_button_theme = material.NewTheme()
	red_button_theme.ContrastBg = red_button_color

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

func sentence_input(gtx C) D {
	editor_style := material.Editor(theme, lineEditor, "Hint")
	editor_style.Font.Style = font.Italic
	border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
	spaced := func(gtx C) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, editor_style.Layout)
	}
	dimensions := border.Layout(gtx, spaced)
	return dimensions 
}

func target_word_selection(gtx C) D {
	flex := layout.Flex{}
	list_style := material.List(theme, button_list)

	buttons_new := func(gtx C, i int) D {
		var button WordButton
		text := strconv.Itoa(i) 
		text += " - "
		text += words[i] 

		button.word = words[i] 
		button.area = new(widget.Clickable)
		button.text = text
		button.chosen = false
		button.number = i

		word_buttons = append(word_buttons, button)
		update_words = false;
		return material.Button(theme, button.area, button.text).Layout(gtx)
	}

	buttons_old := func(gtx C, i int) D {
		button := word_buttons[i]
		if button.chosen == true {
			return material.Button(red_button_theme, button.area, button.text).Layout(gtx)
		}
		return material.Button(theme, button.area, button.text).Layout(gtx)
	}

	button_generator := buttons_old

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

func word_being_curated(gtx C) D { 
	text := current_curated_word.word
	label := material.H3(theme, text)
	return label.Layout(gtx)
}

func definitions(gtx C) D {
	flex := layout.Flex{}
	list_style := material.List(theme, definition_list)

	buttons_new := func(gtx C, i int) D {
		var button WordButton
		text := strconv.Itoa(i) 
		text += " - "
		text += string(curated_word_definitions[i])

		button.word = string(curated_word_definitions[i]) 
		button.area = new(widget.Clickable)
		button.text = text
		button.chosen = false
		button.number = i

		definition_buttons = append(definition_buttons, button)
		update_word_definitions = false;
		return material.Button(theme, button.area, button.text).Layout(gtx)
	}

	buttons_old := func(gtx C, i int) D {
		button := definition_buttons[i]
		if button.chosen == true {
			return material.Button(red_button_theme, button.area, button.text).Layout(gtx)
		}
		return material.Button(theme, button.area, button.text).Layout(gtx)
	}

	button_generator := buttons_old

	if (update_word_definitions) {
		definition_buttons = definition_buttons[:0]
		button_generator = buttons_new
	}

	update_word_definitions = false

	anon_list := func(gtx C) D {
		return list_style.Layout(gtx, len(curated_word_definitions), button_generator)
	}
	return flex.Layout(gtx, layout.Flexed(1, anon_list))
}

func dashboard(gtx layout.Context, th *material.Theme) layout.Dimensions {
	submited := false

	for {
		e, ok := lineEditor.Update(gtx)

		if !ok {
			break
		}

		if e, ok := e.(widget.SubmitEvent); ok {
			sentence = string(e.Text)
			words = strings.Fields(e.Text)
			lineEditor.SetText("")
			update_words = true;
			control_state.next_mode()
			submited = true
		}
	}

	for {
		key_event, ok := gtx.Source.Event(key.Filter{})
		if !ok {
			break
		}

		if submited  == true {
			break
		}


		k, ok := key_event.(key.Event)

		if !ok {
			continue
		}

		
		if k.Name == "⏎" && k.State == key.Press {
			control_state.next_mode()
		}

		if (control_state.value == choice) && (k.State == key.Press) {

			if len(k.Name) != 1 {
				break
			}


			value, err := strconv.ParseInt(string(k.Name), 10, 64)
			if err != nil {
				break
			}

			if value >= int64(len(word_buttons)) {
				break
			}
			
			word_buttons[value].chosen = !word_buttons[value].chosen
		}

		if (control_state.value == curate){
			if reset_curate == true {
				current_curated_word = next_chosen_button()
				curated_word_definitions = controller.WordEquivalents(controller.Word(current_curated_word.word))
				reset_curate = false
			}

			if (k.Name == "Tab") {
				update_word_definitions = true
				current_curated_word = next_chosen_button()
				curated_word_definitions = controller.WordEquivalents(controller.Word(current_curated_word.word))
			}
		}
	}

	stage_one := []layout.Widget{
		sentence_input,
	}

	stage_two := []layout.Widget{
		target_word_selection,
	}

	stage_three := []layout.Widget{
		word_being_curated,
		definitions,
	}

	widgets := &stage_one

	if (control_state.value == 0) {
		widgets = &stage_one
	} else if (control_state.value == 1) {
		widgets = &stage_two
		reset_curate = true
	} else if (control_state.value == 2) { 
		widgets = &stage_three
	}

	spaced := func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, (*widgets)[i])
	}
	list_style := material.List(th, list)
	dimensions := list_style.Layout(gtx, len((*widgets)), spaced)
	return dimensions
}
