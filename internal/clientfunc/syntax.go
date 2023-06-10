package clientfunc

import "fmt"

func printRegisterSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: register <login> <password>")
}

func printLoginSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: login <login> <password>")
}

func printSetCardSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: setcard <cardname> <cardnumber> <cardholder name> <cardholder surname> <card valid till date in format 'dd:mm:yyyy'> <cvv code>")
}

func printGetCardSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: getcard <cardname>")
}

func printListCardsSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: listcards")
}

func printSetLoginCredsSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax:setlogincreds <metaname> <sitename> <login> <password>")
}

func printGetLoginCredsSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: getlogincreds <logincreds name>")
}

func printListLoginCredsSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: listlogincreds")
}

func printSetNoteSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax:setnote <name of the note> <text of the note>")
}

func printGetNoteSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: getnote <name of the note>")
}

func printListNotesSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: listnotes")
}

func printSetBinarySyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax:setbinary <name of the binary/binary file>")
	fmt.Println("Binary file to be send has to be placed in 'client/filessnd/' folder, and has to be named as stated above")
}

func printGetBinarySyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: getbinary <name of the binary>")
}

func printListBinariesSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: listbinaries")
}
