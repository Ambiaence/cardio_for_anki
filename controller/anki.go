package controller

import (
    "bytes"
    "encoding/json"
    "net/http"
)

var deck_name = GlobalSettings.DeckName

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
        panic("Error marshaling anki card.")
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
}
