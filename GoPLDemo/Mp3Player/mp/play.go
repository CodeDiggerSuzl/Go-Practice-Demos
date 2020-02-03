package mp

import "fmt"

type Player interface {
	Play(source string)
}

func Play(source, mtype string) {
	var p Player
	switch mtype {
	// 接口赋值
	case "MP3":
		p = &MP3Player{}
	case "WNV":
		p = &MWVPlayer{}
	default:
		fmt.Println(" 🤭 Not support media type")
		return
	}
	p.Play(source)
}
