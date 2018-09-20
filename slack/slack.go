package slack

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    scryfall "github.com/heroku/scrivener/scryfall"
)

type Response interface {}

type response struct {
    ResponseType string `json:"response_type"`
    Text string `json:"text"`
}

func NewResponse(responseType string, text string) Response {
    return response {
        ResponseType: responseType,
        Text: text,
    }
}

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
    Selected   []actionOption `json:"selected_options"`
}

type Callback struct {
    ID string `json:"callback_id"`
    Actions []Action `json:"actions"`
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

type Card interface {
    Response
}

type card struct {
    Attachments []Attachment `json:"attachments"`
    Display string `json:"response_type"`
    Name string `json:"text"`
}

func NewCard(scry scryfall.Card, linkOnly bool) Card {
    ret := card {
        Attachments: []Attachment {},
        Display: "in_channel",
    }

    if(scry.Images.Normal != ""){
        imageAttach := Attachment {
            Title: scry.Name,
            Link: scry.Link,
        }
        if(!linkOnly){
            imageAttach.URL = scry.Images.Normal
        }
        ret.Attachments = append(ret.Attachments, imageAttach)
    }else if(scry.Faces != nil){
        for _, face := range scry.Faces {
            imageAttach := Attachment {
                Title: face.Name,
                Link: face.Link,
            }
            if(!linkOnly){
                imageAttach.URL = face.Images.Normal
            }
            ret.Attachments = append(ret.Attachments, imageAttach)
        }
    }else{
        imageAttach := Attachment {
            Title: fmt.Sprintf("%s - Image not available.", scry.Name),
            Link: scry.Link,
        }
        ret.Attachments = append(ret.Attachments, imageAttach)
    }


    return ret
}

type CardChoice interface {}

type cardChoice struct {
    Attachments []Attachment `json:"attachments"`
    Display string `json:"response_type"`
    Text string `json:"text"`
}

func NewCardChoice(searchString string, cardList []scryfall.Card, linkOnly bool) CardChoice {
    numCards := len(cardList)
    log.Printf("New Card Choice for %d cards.", numCards)

    actions := []Action{}
    var callbackID string

    if(numCards > 5){
        // Slack won't let us use more than five buttons, so we need a menu.
        callbackID = "cardMenu"
        // Let's try adding the menu before we trash the buttons.
        options := []actionOption {}
        for _, card := range cardList {
            opt := actionOption {
                Text: card.Name,
                Value: card.Name,
            }
            if(linkOnly){
                opt.Value = fmt.Sprintf("%s --link", opt.Value)
            }
            options = append(options, opt)
        }
        menuAction := Action {
            Name: "card",
            Text: "Please choose a card",
            ActionType: "select",
            Options: options,
        }
        actions = append(actions, menuAction)
    } else {
        callbackID = "cardButton"
        // Now the actual buttons.
        for _, card := range cardList {
            actions = append(actions, Action{
                Name: "card",
                Text: card.Name,
                ActionType: "button",
                Value: card.Name,
            })
        }
    }

    attachments := []Attachment{}
    attachments = append(attachments, Attachment{
        Text: "There were multiple matches.",
        Fallback: "There were multiple matches.",
        CallbackID: callbackID,
        Color: "acacac",
        AttachmentType: "default",
        Actions: actions,
    })

    return cardChoice {
        Attachments: attachments,
        Display: "in_channel",
    }
}

func Respond(message Response, url string) {
    jsonString, _ := json.Marshal(message)
    http.Post(url, "application/json", bytes.NewBuffer(jsonString))
}