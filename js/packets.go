package js

type RequestSendEmail struct {
	To   string `json:"to"`
	From string `json:"from"`
	Body string `json:"body"`
}

type ResponseSendEmail struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}
