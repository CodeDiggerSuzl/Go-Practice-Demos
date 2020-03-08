package mlib

import (
	"testing"
)

func TestOpts(t *testing.T) {

	mm := NewMusicManager()
	if mm == nil {
		t.Error("NewMusicManager failed")
	}
	if mm.Len() != 0 {
		t.Error("NewMusicManager failed, not empty!")
	}

	m0 := &MusicEntry{"1", "song1", "Unknown", "http://netease/12", "MP3"}

	m1 := &MusicEntry{"2", "song2", "Unknown", "http://netease/12", "MP3"}

	mm.Add(m0)
	mm.Add(m1)

	m := mm.Find(m0.Name)
	if m == nil {
		t.Error("NewMusicManager.Find failed")
	}
	if m.Id != m0.Id || m.Artist != m0.Artist || m.Name != m0.Name || m.Source != m0.Source || m.Type != m0.Type {
		t.Error("NewMusicManager.Find failed,not match")
	}

	a, _ := mm.BlurFind("song")
	if len(a) != 2 {
		t.Error("Blur Find failed")
	}

	g, err := mm.Get(0)
	if g == nil {
		t.Error("Get failed", err)
	}

	r := mm.Remove(1)
	if r == nil {
		t.Error("Remove failed")
	}
}
