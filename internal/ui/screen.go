package ui

import (
	"github.com/gdamore/tcell/v2"
)

func InitScreen() (tcell.Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := s.Init(); err != nil {
		return nil, err
	}
	s.EnableMouse()
	s.Clear()
	return s, nil
}
