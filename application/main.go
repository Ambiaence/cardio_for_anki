package main

import (
	"fmt"
	"strconv"
	"log"
	"os"
	"strings"
	"image/color"
	"errors"

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
var submited bool

const (
	start = iota
	choice = iota
	curate = iota
	create = iota
)

type (
	D = layout.Dimensions
	C = layout.Context
)

type EquivalentWords []string

type EquivalentMap map[int]EquivalentWords

var (
	equivalents = make(EquivalentMap)

	theme *material.Theme
	red_button_theme *material.Theme

	update_words = true
	update_equivalents = true

	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}

	word_buttons = []WordButton{}

	chosen_equivalents = []string{} 

	definition_buttons = []WordButton{}

	definition_list = &widget.List{
		List: layout.List{
			Axis: layout.Horizontal,
		},

	}

	chosen_equivalents_list = &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
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

	words = []string{"one","two","three"}
	control_state = Mode{value: 0}

	current_curated_position = 0 
	current_curated_word *WordButton
	curated_word_definitions controller.WordList 

	reset_curate = true
	update_word_definitions = true

	sentence string 
	source_sentence string
)


func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		w.Option(app.Title("Cardio For Anki"));
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

func chosen_equivalents_display(gtx C) D {
	flex := layout.Flex{}
	list_style := material.List(theme, chosen_equivalents_list)

	list_new := func(gtx C, i int) D {
		if i == 0 {
			for key, value := range equivalents {
				line := word_buttons[key].word + " -> "
				for _, v := range value {
					line = line + " / " + v 
				}
				chosen_equivalents = append(chosen_equivalents, line)
			}
				
		}
		string_line := chosen_equivalents[i]
		return material.H6(theme, string_line).Layout(gtx)
	}

	list_old := func(gtx C, i int) D {
		string_line := chosen_equivalents[i]
		return material.H6(theme, string_line).Layout(gtx)
	}

	list_generator := list_old

	if (update_equivalents) {
		chosen_equivalents = chosen_equivalents[:0]
		list_generator = list_new
	}

	update_equivalents = false

	anon_list := func(gtx C) D {
		return list_style.Layout(gtx, len(equivalents), list_generator)
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

func handle_sentence_editor(gtx layout.Context, th *material.Theme) {
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
}

func source_sentence_label(gtx layout.Context) layout.Dimensions{
	label := "Source Language Sentence: " + source_sentence 
	return material.H6(theme, label).Layout(gtx)
}

func target_sentence_label(gtx layout.Context) layout.Dimensions{
	label := "Target Language Sentence: " + sentence 
	return material.H6(theme, label).Layout(gtx)
}

func handle_curate_inputs(k key.Event, gtx layout.Context) {
	if reset_curate == true {
		current_curated_word = next_chosen_button()
		curated_word_definitions = controller.WordEquivalents(controller.Word(current_curated_word.word))
		reset_curate = false
	}

	if (k.Name == "Tab") {
		var equivalent_list EquivalentWords

		for index, button  := range definition_buttons {
			_ = index

			if button.chosen == false {
				continue;
			}

			equivalent_list = append(equivalent_list, button.word)
			fmt.Println(equivalent_list)
		}

		equivalents[current_curated_word.number] = equivalent_list

		update_word_definitions = true
		update_equivalents = true
		current_curated_word = next_chosen_button()
		curated_word_definitions = controller.WordEquivalents(controller.Word(current_curated_word.word))
	}
}

func handle_choice_inputs(k key.Event, gtx layout.Context) error {
	if len(k.Name) != 1 {
		return errors.New("len(k.Name) !=1.")
	}

	value, err := strconv.ParseInt(string(k.Name), 10, 64)

	if err != nil {
		return errors.New("Could not parse integer.")
	}

	if value >= int64(len(word_buttons)) {
		return errors.New("value <= int64*len(word_buttons))")
	}
	
	word_buttons[value].chosen = !word_buttons[value].chosen

	return nil
}

func handle_state_related_inputs(gtx layout.Context) {
	for {
		key_event, ok := gtx.Source.Event(key.Filter{})

		if !ok {
			break
		}

		if submited == true {
			break
		}


		k, ok := key_event.(key.Event)

		if !ok {
			continue
		}

		
		if k.Name == "⏎" && k.State == key.Press {
			if control_state.value == curate {
				source_sentence = controller.TranslateSentence(sentence)
			}
			control_state.next_mode()
		}

		if (control_state.value == choice) && (k.State == key.Press) {
			err := handle_choice_inputs(k ,gtx)

			if err != nil {
				break
			}
		}

		if (control_state.value == curate) {
			handle_curate_inputs(k, gtx)
		}

		if (control_state.value == curate) && (k.State == key.Press) {
			if len(k.Name) != 1 {
				break
			}

			value, err := strconv.ParseInt(string(k.Name), 10, 64)
			if err != nil {
				break
			}

			if value >= int64(len(definition_buttons)) {
				break
			}
			
			definition_buttons[value].chosen = !definition_buttons[value].chosen
		}

		if (control_state.value == create) && (k.Name == "Tab") {
			front := ""
			back := ""

			front = sentence
			back = source_sentence + "/n" 

			for _, equivalents := range chosen_equivalents {
				back = back + "[" + equivalents + "] " 
			}
			controller.CreateCard(front, back)
		}
	}
}

func dashboard(gtx layout.Context, th *material.Theme) layout.Dimensions {
	submited = false

	handle_sentence_editor(gtx, th)

	handle_state_related_inputs(gtx)

	stage_one := []layout.Widget{
		sentence_input,
	}

	stage_two := []layout.Widget{
		target_word_selection,
	}

	stage_three := []layout.Widget{
		word_being_curated,
		definitions,
		chosen_equivalents_display,
	}

	stage_four := []layout.Widget{
		chosen_equivalents_display,
		source_sentence_label,
		target_sentence_label,
	}

	widgets := &stage_one

	if (control_state.value == 0) {
		widgets = &stage_one
	} else if (control_state.value == 1) {
		widgets = &stage_two
		reset_curate = true
	} else if (control_state.value == 2) { 
		widgets = &stage_three
	} else if (control_state.value == 3) { 
		widgets = &stage_four
	}

	spaced := func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, (*widgets)[i])
	}

	list_style := material.List(th, list)
	dimensions := list_style.Layout(gtx, len((*widgets)), spaced)
	return dimensions
}
