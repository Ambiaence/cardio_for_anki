package controller

import (
    "os"
    "io/ioutil"
    "encoding/json"
    "fmt"
)

type Settings struct {
    Source string
    Target string
    Tts string
    Translate string
    Equivalents string
    DeckName string
}

func read_global_settings() Settings {
    var settings Settings
    jsonFile, err := os.Open("settings.json")

    if err != nil { 
        fmt.Println(err)
        panic("Settings file not found")
    }

    bytes, _ := ioutil.ReadAll(jsonFile)

    defer jsonFile.Close()

    var data map[string]string

    json.Unmarshal(bytes, &data)

    settings.Source = data["source"]
    settings.Target = data["target"]
    settings.Tts = data["tts"]
    settings.Translate = data["translate"]
    settings.Equivalents = data["equivalents"]
    settings.DeckName = data["deck_name"]

    return settings
}

var GlobalSettings = read_global_settings()
