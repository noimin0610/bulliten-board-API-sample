package messages

import (
	"encoding/json"
	"net/http"
	"time"
)

type Message struct {
	Name      string `json:"name"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

func AllMessages() []Message {
	return []Message{
		{"ヤンマガ読者", "漫画は面白いです。", "2021-03-24 21:00:00"},
		{"Glossom社員", "そうですね。", "2021-03-24 21:00:01"},
	}
}

func AppendMessage(name string, text string) Message {
	n := time.Now()
	return Message{
		name, text, n.Format("2006-01-02 15:04:05"),
	}
}

func Messages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		messages, _ := json.Marshal(AllMessages())
		w.Write(messages)

	case http.MethodPost:
		AppendMessage(r.FormValue("name"), r.FormValue("text"))
		w.WriteHeader(http.StatusCreated)

	default:
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
