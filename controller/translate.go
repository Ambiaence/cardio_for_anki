package controller

import (
    "context"
    "fmt"

    translate "cloud.google.com/go/translate/apiv3"
    translatepb "cloud.google.com/go/translate/apiv3/translatepb"
)

func TranslateSentence(input string) string {
    return google_translate("Dummy")
}

func google_translate(input string) string {
    ctx := context.Background()

    c, err := translate.NewTranslationClient(ctx)

    if err != nil {
	fmt.Println(err)
    }

    defer c.Close()

    req := &translatepb.TranslateTextRequest{
	Contents: []string{"ich bin anna"},
	SourceLanguageCode: "en",
	TargetLanguageCode: "de",
    }

    resp, err := c.TranslateText(ctx, req)
    if err != nil {
	fmt.Println(err)
    }

    fmt.Println(resp.Translations)
    return "return"
}
