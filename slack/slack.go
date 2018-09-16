package slack

import (
    "fmt"
    "strings"
    scryfall "github.com/heroku/scrivener/scryfall"
)

type Action struct {
    Name       string `json:"name"`
    Text       string `json:"text"`
    ActionType string `json:"type"`
    Value      string `json:"value"`
}

type Attachment struct {
    Text           string   `json:"text"`
    Fallback       string   `json:"fallback"`
    CallbackID     string   `json:"callback_id"`
    Color          string   `json:"color"`
    AttachmentType string   `json:"attachment_type"`
    Actions        []Action `json:"actions"`
    Title          string   `json:"title"`
    URL            string   `json:"image_url"`
}

type Card interface {}

type card struct {
    Attachments []Attachment `json:"attachments"`
    Display string `json:"response_type"`
    Name string `json:"text"`
}

func NewCard(scry scryfall.Card) Card {
    imageAttach := Attachment {
        Title: scry.Name,
        URL: scry.Images.Large,
    }

    ret := card {
        Attachments: []Attachment {},
        Display: "in_channel",
    }

    ret.Attachments = append(ret.Attachments, imageAttach)

    return ret
}

type CardChoice interface {}

type cardChoice struct {
    Attachments []Attachment `json:"attachments"`
    Display string `json:"response_type"`
    Text string `json:"text"`
}

func NewCardChoice(searchString string, cardList []scryfall.Card) CardChoice {
    var text strings.Builder
    text.WriteString(fmt.Sprintf("Searching for '_%s_' returned multiple results:", searchString))
    for _, card := range cardList {
        text.WriteString(fmt.Sprintf("\n%s", card.Name))
    }
    text.WriteString("\n")

    return cardChoice {
        Display: "ephemeral",
        Text: text.String(),
    }
}