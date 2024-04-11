package controller

import (
    "fmt"
    "regexp"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "log"
)

type Word string

type WordList []Word 

type Category struct {
    PartOfSpeech string
    Language string
    Definitions map[string]string
}

func WordEquivalents(word Word) WordList {
	wikiEquivalents(word)
	return WordList{"nice", "sick"}
}

func wikiEquivalents(word Word) WordList {
    url := fmt.Sprintf("https://en.wiktionary.org/api/rest_v1/page/definition/%s", word)
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }

    req.Header.Set("accept", "application/json; charset=utf-8; profile=\"https://www.mediawiki.org/wiki/Specs/definition/0.8.0\"")

    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }

    defer resp.Body.Close() 

    body, err := ioutil.ReadAll(resp.Body)
    

    if err != nil {
        log.Fatal(err)
    }


    var data map[string]([]interface{})

    err = json.Unmarshal(body, &data)

    if err != nil {
	fmt.Println("Err != nil")
    }

    raw, ok := data["de"]

    if !ok {
	fmt.Println("Error reading", word)
    }

    err = json.Unmarshal(body, &data)


    for _, value := range raw {  
	r := value.(map[string]interface{})
	defs := r["definitions"].([]interface{})
	str := defs[0].(map[string]interface{})
	fmt.Println(extractTextFromATag(str["definition"].(string)))
    }

    return WordList{"Nice"}
}

//This functions was ripped from Google Gemini.
func extractTextFromATag(html string) string { 
  // Regular expression with improved handling of attributes and whitespace:
  pattern := `<a[^>]*>(.*?)</a>`

  re := regexp.MustCompile(pattern)

  match := re.FindStringSubmatch(html) 

  if match != nil {
    return match[1] // Return the captured text inside the <a> tag
  } else {
    return "No text found" 
  }
}
