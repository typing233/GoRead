package ui

import (
	"github.com/gdamore/tcell/v2"
)

type KeyAction int

const (
	ActionNone KeyAction = iota
	ActionScrollDown
	ActionScrollUp
	ActionHalfPageDown
	ActionHalfPageUp
	ActionPageDown
	ActionPageUp
	ActionTop
	ActionBottom
	ActionNextChapter
	ActionPrevChapter
	ActionQuit
)

type inputState int

const (
	stateIdle inputState = iota
	stateGPending
)

type InputHandler struct {
	state inputState
}

func NewInputHandler() *InputHandler {
	return &InputHandler{state: stateIdle}
}

func (ih *InputHandler) Handle(ev *tcell.EventKey) KeyAction {
	if ih.state == stateGPending {
		ih.state = stateIdle
		if ev.Key() == tcell.KeyRune && ev.Rune() == 'g' {
			return ActionTop
		}
		// Not 'g' — fall through and process as normal key
		return ih.processKey(ev)
	}

	return ih.processKey(ev)
}

func (ih *InputHandler) processKey(ev *tcell.EventKey) KeyAction {
	if ev.Key() == tcell.KeyCtrlD {
		return ActionHalfPageDown
	}
	if ev.Key() == tcell.KeyCtrlU {
		return ActionHalfPageUp
	}
	if ev.Key() == tcell.KeyCtrlF || ev.Key() == tcell.KeyPgDn {
		return ActionPageDown
	}
	if ev.Key() == tcell.KeyCtrlB || ev.Key() == tcell.KeyPgUp {
		return ActionPageUp
	}

	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'j':
			return ActionScrollDown
		case 'k':
			return ActionScrollUp
		case 'G':
			return ActionBottom
		case 'g':
			ih.state = stateGPending
			return ActionNone
		case 'q':
			return ActionQuit
		case 'l', 'n':
			return ActionNextChapter
		case 'h', 'p':
			return ActionPrevChapter
		case ' ':
			return ActionPageDown
		}
	}

	if ev.Key() == tcell.KeyDown {
		return ActionScrollDown
	}
	if ev.Key() == tcell.KeyUp {
		return ActionScrollUp
	}

	return ActionNone
}
