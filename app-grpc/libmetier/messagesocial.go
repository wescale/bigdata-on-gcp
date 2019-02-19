package libmetier

import (
	"encoding/json"
	"log"
	"time"
)

// MessageSocial a common social msg
type MessageSocial struct {
	Data      string    `json:"data"`
	User      string    `json:"user"`
	Source    string    `json:"source"`
	Tag       string    `json:"tag"`
	Date      time.Time `json:"timestamp"`
	Sentiment float32   `json:"sentiment"`
	ID        string    `json:"id"`
}

// ToMessageSocial public function
func ToMessageSocial(mstpl []byte) MessageSocial {
	var ms MessageSocial
	err := json.Unmarshal(mstpl, &ms)
	if err != nil {
		log.Println(err)
	}
	return ms
}

// ToByteArray another public function
func (ms MessageSocial) ToByteArray() []byte {
	b, err := json.Marshal(ms)
	if err != nil {
		log.Println(err)
	}
	return []byte(b)
}
