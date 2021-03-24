package messages

import (
	"reflect"
	"testing"
)

func TestAllMessages(t *testing.T) {
	messages := AllMessages()
	ref := []Message{
		{"ヤンマガ読者", "漫画は面白いです。", "2021-03-24 21:00:00"},
		{"Glossom社員", "そうですね。", "2021-03-24 21:00:01"},
	}
	if !reflect.DeepEqual(messages, ref) {
		t.Fatalf("failed:\nAllMessages() = %+v, want %+v", messages, ref)
	}
}

func TestAppendMessage(t *testing.T) {
	message := AppendMessage("ヤンマガチーム", "こんにちは。")
	ref := Message{
		"ヤンマガチーム", "こんにちは。", "",
	}
	if !(message.Name == ref.Name && message.Text == ref.Text && message.Timestamp != "") {
		t.Fatalf("failed:\nAppendMessage() = %+v, want %+v", message, ref)
	}
}