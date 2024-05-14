package controller

import (
    "os"
    "io/ioutil"
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

var deck_name = read_deck_name()

func read_deck_name() string {
    jsonFile, err := os.Open("settings.json")

    if err != nil { 
        fmt.Println(err)
        panic("Settings file not found")
    }

    bytes, _ := ioutil.ReadAll(jsonFile)

    defer jsonFile.Close()
    var data map[string]string

    json.Unmarshal(bytes, &data)

    key := data["deck_name"]
    return key
}

func CreateCard(front string, back string) {
    data := map[string]interface{}{
        "action":  "addNote",
        "version": 6,
        "params": map[string]interface{}{
            "note": map[string]interface{}{
                "deckName":  deck_name,
                "modelName": "Basic",
                "fields": map[string]string{
                    "Front": front,
                    "Back":  back,
                },
                "options": map[string]bool{
                    "allowDuplicate": false,
                },
            },
        },
    }

    payload_bytes, err := json.Marshal(data)
    if err != nil {
        fmt.Println("Error marshaling JSON:", err)
        return
    }
    body := bytes.NewReader(payload_bytes)

    req, err := http.NewRequest("POST", "http://localhost:8765", body)

    if err != nil {
        panic("Error creating request:")
        return
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}

    resp, err := client.Do(req)

    if err != nil {
        panic("Error sending request:")
        return
    }

    defer resp.Body.Close()
    fmt.Println("Anki Response Status:", resp.Status)
}
