package controller

import (
	"os/exec"
)

func GenerateSpokenWord(word string, language string)  {
	espeakSpokenWord(word, language)
}

func espeakSpokenWord(word string, language string) {
	filepath := "temporary/" + word + ".wav"
	cmd := exec.Command("espeak", "-v", language, word, "-w", filepath) 
	_, er := cmd.Output()
	if er != nil {
		return
	}
}
