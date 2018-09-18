package slack

import (
    "log"
    scryfall "github.com/heroku/scrivener/scryfall"
)

type actionOption struct {
    Text string `json:"text"`
    Value string `json:"value"`
}

type Action struct {
    Name       string         `json:"name"`
    Text       string         `json:"text"`
    ActionType string         `json:"type"`
    Value      string         `json:"value"`
    Options    []actionOption `json:"options"`
}

type Attachment struct {
    Text           string   `json:"text"`
    Fallback       string   `json:"fallback"`
    CallbackID     string   `json:"callback_id"`
    Color          string   `json:"color"`
    AttachmentType string   `json:"attachment_type"`
    Actions        []Action `json:"actions"`
    Title          string   `json:"title"`
    Link           string   `json:"title_link"`
    URL            string   `json:"image_url"`
}

type Card interface {}

type card struct {
    Attachments []Attachment `json:"attachments"`
    Display string `json:"response_type"`
    Name string `json:"text"`
}

func NewCard(scry scryfall.Card) Card {
    ret := card {
        Attachments: []Attachment {},
        Display: "in_channel",
    }

    if(scry.Faces == nil){
        imageAttach := Attachment {
            Title: scry.Name,
            URL: scry.Images.Large,
            Link: scry.Link,
        }
        ret.Attachments = append(ret.Attachments, imageAttach)
    }else{
        for _, face := range scry.Faces {
            imageAttach := Attachment {
                Title: face.Name,
                URL: face.Images.Large,
            }
            ret.Attachments = append(ret.Attachments, imageAttach)
        }
    }


    return ret
}

type CardChoice interface {}

type cardChoice struct {
    Attachments []Attachment `json:"attachments"`
    Display string `json:"response_type"`
    Text string `json:"text"`
}

func NewCardChoice(searchString string, cardList []scryfall.Card) CardChoice {
    log.Printf("New Card Choice for %s cards.", len(cardList))

    buttons := []Action{}
    // Let's try adding the menu before we trash the buttons.
    options := []actionOption {}
    for _, card := range cardList {
        options = append(options, actionOption {
            Text: card.Name,
            Value: card.Name,
        })
    }
    menuAction := Action {
        Name: "card",
        Text: "Please choose a card",
        ActionType: "select",
        Options: options,
    }
    buttons = append(buttons, menuAction)

    // Now the actual buttons.
    for _, card := range cardList {
        buttons = append(buttons, Action{
            Name: "card",
            Text: card.Name,
            ActionType: "button",
            Value: card.Name,
        })
    }

    attachments := []Attachment{}
    attachments = append(attachments, Attachment{
        Text: "Please choose one.",
        Fallback: "Please choose one.",
        CallbackID: "choose_card",
        Color: "acacac",
        AttachmentType: "default",
        Actions: buttons,
    })

    return cardChoice {
        Attachments: attachments,
        Display: "in_channel",
    }
}


type Callback struct {
    CallbackID string `json:"callback_id"`
    Actions []Action `json:"actions"`
}
