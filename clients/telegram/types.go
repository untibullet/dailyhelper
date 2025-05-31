package telegram

type UpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int      `json:"update_id"`
	Message *Message `json:"message,omitempty"`
}

type Message struct {
	Text string `json:"text,omitempty"`
	From *User  `json:"from,omitempty"`
	Chat Chat   `json:"chat"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username,omitempty"`
}

type Chat struct {
	ID int `json:"id"`
}
