package slack

import (
    "fmt"
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

func NewCard(name string, image string) Card {
    imageAttach := Attachment {
        Title: name,
        URL: image,
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

func NewCardChoice(searchString string) CardChoice {
    return cardChoice {
        Display: "ephemeral",
        Text: fmt.Sprintf("Searching for '%s' returned multiple results. And sadly the UI for the next step isn't ready.", searchString),
    }
}