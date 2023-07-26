package database

type LoginCreds struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Site     string `json:"site"`
}

type Note struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type Binary struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type Card struct {
	Cardname  string `json:"cardname"`
	Number    string `json:"number"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	ValidTill string `json:"valid till"`
	Code      string `json:"code"`
}

type EncryptedCard struct {
	Cardname string `json:"cardname"`
	Data     string `json:"data"`
}

type EncryptedData struct {
	Name string `json:"name"`
	Data string `json:"data"`
}