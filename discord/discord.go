package discord

import (
	"fmt"

    "github.com/bwmarrin/discordgo"
    scryfall "github.com/heroku/scrivener/scryfall"
)

func EmbedCard(card scryfall.Card) discordgo.MessageEmbed {
	image := discordgo.MessageEmbedImage {
		URL: card.Images.Normal,
		ProxyURL: card.Images.Normal,
	}

	ret := discordgo.MessageEmbed {
		Title: card.Name,
		Color: getColor(card.ColorIdentity),
		Image: &image,
		URL: card.Link,
	}

	return ret
}

func EmbedChoice(cardList []scryfall.Card) discordgo.MessageEmbed {
	prompt := "Please choose from the list below. (Type the number to choose, or 0 to cancel.)"
	for index, card := range cardList {
		prompt = prompt + fmt.Sprintf("\n**[%d]** - %s", index + 1, card.Name)
	}

	// Discord has its limits.
	if(len(prompt) >= 2048) {
		prompt = "Unfortunately there were so many matches that discord won't even let us ask about them all here. Please narrow it down."
	}

	ret := discordgo.MessageEmbed {
		Title: "Multiple Matches Found",
		Description: prompt,
	}

	return ret
}

func getColor(identity []string) int {
	if(len(identity) > 1){
		return 0x997300
	} else if (len(identity) == 0){
		return 0xaaaaaa
	}

	switch(identity[0]){
	case "R":
		return 0xD3202A
	case "G":
		return 0x00733E
	case "U":
		return 0x0E68AB
	case "W":
		return 0xF9FAF4
	case "B":
		return 0x000000
	default:
		return 0x150B00
	}
}

func RespondWithCard(card scryfall.Card, session *discordgo.Session, channelID string) {
	if(len(card.Faces) > 1){
		for _, face := range card.Faces {
			face.ColorIdentity = card.ColorIdentity
			RespondWithCard(face, session, channelID)
		}
	} else {
	    embed := EmbedCard(card)
	    session.ChannelMessageSendEmbed(channelID, &embed)
	}
}