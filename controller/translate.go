package controller

import (
    "os"
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "golang.org/x/net/http2"
    "log"
)

var deepl_key = readDeeplApiKey()

func TranslateSentence(input string) string {
    return DeeplTranslate(input) 
}

func readDeeplApiKey() string {
    fmt.Println("readDeeplApi")
    jsonFile, err := os.Open("keys.json")

    if err != nil { 
        fmt.Println(err)
        panic("key file not found")
    }

    bytes, _ := ioutil.ReadAll(jsonFile)

    defer jsonFile.Close()
    var data map[string]string

    json.Unmarshal(bytes, &data)

    key := data["deepl"]
    fmt.Println("What?", string(key))
    return key
}

func DeeplTranslate(input string) string {
    url := "https://api-free.deepl.com/v2/translate"
    apiKey := deepl_key // Replace with your actual DeepL API key

    data := map[string]interface{}{
        "text":         []string{input},
        "target_lang":  "EN",
        "source_lang":  "DE",
    }

    jsonData, _ := json.Marshal(data)

    headers := http.Header{
        "Authorization":  []string{"DeepL-Auth-Key " + apiKey},
        "User-Agent":     []string{"YourApp/1.2.3"},
        "Content-Length": []string{fmt.Sprintf("%d", len(jsonData))},
        "Content-Type":   []string{"application/json"},
    }

    transport := &http2.Transport{}
    client := &http.Client{Transport: transport}
    request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    request.Header = headers

    response, err := client.Do(request)

    if err != nil {
        fmt.Println(err)
        panic("err != nil")
    }

    defer response.Body.Close()

    fmt.Println("Status Code:", response.StatusCode)

    bodyBytes, err := io.ReadAll(response.Body)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(bodyBytes))

    var json_data map[string]Translations 

    err = json.Unmarshal(bodyBytes, &json_data)

    if err != nil {
	fmt.Println("Err != nil")
    }

    translations := json_data["translations"]
    translation := translations[0]
    text := translation["text"]

    if translations == nil {
        panic("NIL")
    }
    return text 
}
