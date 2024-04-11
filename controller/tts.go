package controller

import (
	"fmt"
	"os/exec"
)

func GenerateSpokenWord(word string, language string)  {
	espeakSpokenWord(word, language)
}

func espeakSpokenWord(word string, language string) {
	filepath := "temporary/" + word + ".wav"
	cmd := exec.Command("espeak", "-v", "de", word, "-w", filepath) 
	output, er := cmd.Output()
	if er != nil {
		fmt.Println(er)
		return
	}
	fmt.Println(output)
}
