package mp

import (
	"fmt"
	"time"
)

type MP3Player struct {
	stat     int
	progress int
}

func (p *MP3Player) Play(source string) {
	fmt.Println("🗣 mp3 playing")
	p.progress = 0
	for p.progress < 100 {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("🚥🚥🚥🚥🚥🚥")
		p.progress += 5
	}
	fmt.Println("✅ PlayMusic done")
}
