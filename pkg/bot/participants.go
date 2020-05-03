package bot

import (
	"fmt"
	"strings"
	"time"
)

type participant struct {
	name    string
	addedAt time.Time
}

func (b *bot) addParticipant(p string) {
	if len(strings.TrimSpace(p)) > 0 {
		b.participants[strings.TrimSpace(p)] = participant{
			name:    p,
			addedAt: time.Now(),
		}
	}
}

func (b *bot) removeParticipant(p string) {
	delete(b.participants, p)
}

func (b *bot) printParticipants() string {
	var participants []string
	for key, _ := range b.participants {
		participants = append(participants, key)
	}
	return fmt.Sprintf("%v", strings.Join(participants, ", "))
}
