package slack

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
        Display: "ephemeral",
        Name: name,
    }

    ret.Attachments = append(ret.Attachments, imageAttach)

    return ret
}