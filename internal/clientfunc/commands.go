package clientfunc

import (
	"fmt"
	"strings"

	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/database"
	"github.com/gambruh/gophkeeper/internal/helpers"
)

// Cards
func (c *Client) SetCardCommand(input []string) {
	input = helpers.SplitFurther(input)

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

	err := c.saveCardInStorage(cardData)
	if err != nil {
		fmt.Println("error in client saving data to storage in SetCardCommand:", err)
	}
	fmt.Printf("Card %s saved to the storage!\n", cardData.Cardname)
}

func (c *Client) GetCardCommand(input []string) {
	input = helpers.SplitFurther(input)
	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 2 {
		printGetCardSyntax()
		return
	}

	cardname := input[1]

	card, err := c.getCardFromStorage(cardname)
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
	input = helpers.SplitFurther(input)

	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 1 {
		printListCardsSyntax()
		return
	}

	cards, err := c.listCardsFromStorage()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Cards:")
		for _, card := range cards {
			fmt.Println("  ", card)
		}
	}
}

// LoginCreds
func (c *Client) SetLoginCredsCommand(input []string) {
	input = helpers.SplitFurther(input)

	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 5 {
		printSetLoginCredsSyntax()
		return
	}
	var logincreds database.LoginCreds
	for i, data := range input {
		switch i {
		case 0:
		case 1:
			logincreds.Name = data
		case 2:
			logincreds.Site = data
		case 3:
			logincreds.Login = data
		case 4:
			logincreds.Password = data
		}
	}

	err := c.saveLoginCredsInStorage(logincreds)
	if err != nil {
		fmt.Println("error in client saving data to storage in SetLoginCredsCommand:", err)
	}
	fmt.Printf("Login credentials %s saved to the storage!\n", logincreds.Name)
}

func (c *Client) GetLoginCredsCommand(input []string) {
	input = helpers.SplitFurther(input)

	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 2 {
		printGetLoginCredsSyntax()
		return
	}

	logincredname := input[1]

	logincreds, err := c.getLoginCredsFromStorage(logincredname)
	if err != nil {
		if err == ErrDataNotFound {
			fmt.Println("No data in local storage")
			return
		} else {
			fmt.Println("error when trying to get data from local storage:", err)
			return
		}
	}
	//result of the command, if no errors
	fmt.Printf("%+v\n", logincreds)
}

func (c *Client) ListLoginCredsCommand(input []string) {
	input = helpers.SplitFurther(input)

	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 1 {
		printListLoginCredsSyntax()
		return
	}

	logincreds, err := c.listLoginCredsFromStorage()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Login credentials:")
		for _, logincred := range logincreds {
			fmt.Println("  ", logincred)
		}
	}

}

// Notes
func (c *Client) SetNoteCommand(input []string) {
	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	notename, notetext, ok := strings.Cut(input[1], " ")
	if !ok {
		printSetNoteSyntax()
		return
	}

	if !strings.HasPrefix(notetext, `"`) && !strings.HasSuffix(notetext, `"`) {
		printSetNoteSyntax()
		return
	}
	notetext = notetext[1 : len(notetext)-1]

	var note database.Note
	note.Name = notename
	note.Text = notetext

	err := c.saveNoteInStorage(note)
	if err != nil {
		fmt.Println("error in client saving data to storage in SetNoteCommand:", err)
	}
	fmt.Printf("Note %s saved to the storage!\n", note.Name)
}

func (c *Client) GetNoteCommand(input []string) {
	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 2 {
		printGetNoteSyntax()
		return
	}

	notename := input[1]

	note, err := c.getNoteFromStorage(notename)
	if err != nil {
		if err == ErrDataNotFound {
			fmt.Println("No data in local storage")
			return
		} else {
			fmt.Println("error when trying to get data from local storage:", err)
			return
		}
	}
	//result of the command, if no errors
	fmt.Printf("%+v\n", note)
}

func (c *Client) ListNotesCommand(input []string) {
	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 1 {
		printListNotesSyntax()
		return
	}

	notes, err := c.listNotesFromStorage()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Notes:")
		for _, note := range notes {
			fmt.Println("  ", note)
		}
	}
}

func (c *Client) SetBinaryCommand(input []string) {
	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 2 {
		printSetBinarySyntax()
		return
	}

	binaryname := input[1]

	binary, err := helpers.PrepareBinary(binaryname, config.ClientCfg.BinInputFolder)
	if err != nil {
		fmt.Println("smth wrong with binary filepath:", err)
		return
	}

	err = c.saveBinaryInStorage(binary)
	if err != nil {
		fmt.Println("error in client saving data to storage in SetBinaryCommand:", err)
	}
	fmt.Printf("Binary %s saved to the storage!\n", binary.Name)
}

func (c *Client) GetBinaryCommand(input []string) {
	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 2 {
		printGetBinarySyntax()
		return
	}

	binaryname := input[1]

	binary, err := c.getBinaryFromStorage(binaryname)
	if err != nil {
		if err == ErrDataNotFound {
			fmt.Println("No data in local storage")
			return
		} else {
			fmt.Println("error when trying to get data from local storage:", err)
			return
		}
	}
	//result of the command, if no errors
	fmt.Printf("binary file named %s has been created in ./filesrcv folder\n", binary.Name)
}

func (c *Client) ListBinariesCommand(input []string) {
	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 1 {
		printListBinariesSyntax()
		return
	}

	binaries, err := c.listBinariesFromStorage()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Binaries:")
		for _, binary := range binaries {
			fmt.Println("  ", binary)
		}
	}
}
