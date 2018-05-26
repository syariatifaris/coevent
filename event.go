package covent

type Event struct {
	Name    string      `json:"name"`
	Content interface{} `json:"content"`
}

type RegEvtData struct {
	Topic  string   `json:"topic"`
	Events []string `json:"events"`
}
