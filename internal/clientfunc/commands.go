package clientfunc

import (
	"fmt"

	"github.com/gambruh/gophkeeper/internal/database"
)

// Cards
func (c *Client) SetCardCommand(input []string) {
	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 7 {
		printSetCardSyntax()
		return
	}
	var cardData database.Card
	for i, data := range input {
		switch i {
		case 0:
		case 1:
			cardData.Cardname = data
		case 2:
			cardData.Number = data
		case 3:
			cardData.Name = data
		case 4:
			cardData.Surname = data
		case 5:
			cardData.ValidTill = data
		case 6:
			cardData.Code = data
		}
	}

	err := c.sendCardToDB(cardData)
	switch err {
	case ErrMetanameIsTaken:
		fmt.Println("There are already card with this name in database. Please provide new cardname or edit current")
	case ErrLoginRequired:
		fmt.Println("Please login to the server")
	case nil:
		fmt.Println("Saved card to the vault")
	default:
		fmt.Println("error in client sending data to database in SetCardCommand:", err)
	}

	err = c.saveCardInStorage(cardData)
	if err != nil {
		fmt.Println("error in client saving data to storage in SetCardCommand:", err)
	}
}

func (c *Client) GetCardCommand(input []string) {
	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 2 {
		printGetCardSyntax()
		return
	}

	cardname := input[1]

	card, err := c.getCardFromDB(cardname)
	switch err {
	case ErrDataNotFound:
		fmt.Println("Card with that name not found in the DB")
	case ErrLoginRequired:
		fmt.Println("Please login to the server")
	case ErrBadRequest:
		fmt.Println("Please contact devs to change API interaction, wrong request")
	case ErrServerIsDown:
		fmt.Println("Internal server error")
	case nil:
		fmt.Printf("%+v\n", card)
	}

	card, err = c.getCardFromLocalStorage(cardname)
	if err != nil {
		if err == database.ErrDataNotFound {
			fmt.Println("No data in local storage")
			return
		} else {
			fmt.Println("error when trying to get card data from local storage:", err)
			return
		}
	}
	//result of the command, if no errors
	fmt.Printf("%+v\n", card)
}

func (c *Client) ListCardsCommand(input []string) {

	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 1 {
		printListCardsSyntax()
		return
	}

	cards, err := c.listCardsFromDB()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Cards:")
	for _, card := range cards {
		fmt.Println("  ", card)
	}
	fmt.Println("Use getcard <cardname> to acquire card data")
}
