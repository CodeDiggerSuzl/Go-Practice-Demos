package mp

import (
	"fmt"
	"time"
)

type MWVPlayer struct {
	stat     int
	progress int
}

func (p *MWVPlayer) Play(source string) {
	fmt.Println("🗣 MWV playing")
	p.progress = 0
	for p.progress < 100 {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("🚥🚥🚥🚥🚥🚥")
		p.progress += 5
	}
	fmt.Println("✅ PlayMusic done")
}
