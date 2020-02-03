package mlib

import (
	"errors"
	"strings"
)

type MusicEntry struct {
	Id     string
	Name   string
	Artist string
	Source string
	Type   string
}

type MusicManager struct {
	musics []MusicEntry
}

func NewMusicManager() *MusicManager {
	return &MusicManager{make([]MusicEntry, 0)}
}

func (m *MusicManager) Len() int {
	return len(m.musics)
}

func (m *MusicManager) Get(idx int) (music *MusicEntry, err error) {
	if idx < 0 || idx >= len(m.musics) {
		return nil, errors.New("index out of range!")
	}
	return &m.musics[idx], nil
}

// find by name
func (m *MusicManager) Find(name string) *MusicEntry {
	if len(m.musics) == 0 {
		return nil
	}
	for _, m := range m.musics {
		if m.Name == name {
			return &m
		}
	}
	return nil
}

func (m *MusicManager) BlurFind(name string) (result []MusicEntry, err error) {
	rst := make([]MusicEntry, 0)

	if len(m.musics) == 0 {
		return nil, nil
	}

	for _, m := range m.musics {
		if strings.Contains(m.Name, name) {
			rst = append(rst, m)
		}

	}

	return rst, nil
}

func (m *MusicManager) Add(music *MusicEntry) {
	m.musics = append(m.musics, *music)
}

func (m *MusicManager) Remove(index int) *MusicEntry {
	if index < 0 || index > len(m.musics) {
		return nil
	}
	removeMusic := &m.musics[index]
	m.musics = append(m.musics[:index], m.musics[index+1:]...)
	return removeMusic
}
