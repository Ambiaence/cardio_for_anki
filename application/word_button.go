package main

import (
	"gioui.org/widget"
)

type WordButton struct {
	text string
	word string
	chosen bool
	area *widget.Clickable
	number int	
}

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
