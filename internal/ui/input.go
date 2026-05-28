package ui

import (
	"time"

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
	state   inputState
	timer   *time.Timer
	pending chan KeyAction
}

func NewInputHandler() *InputHandler {
	return &InputHandler{
		state:   stateIdle,
		pending: make(chan KeyAction, 1),
	}
}

func (ih *InputHandler) Handle(ev *tcell.EventKey) KeyAction {
	select {
	case action := <-ih.pending:
		ih.cancelTimer()
		if ev.Rune() == 'g' && action == ActionNone {
			return ActionTop
		}
		result := ih.processKey(ev)
		if result != ActionNone {
			return result
		}
		return action
	default:
	}

	switch ih.state {
	case stateIdle:
		return ih.processKey(ev)
	case stateGPending:
		ih.cancelTimer()
		ih.state = stateIdle
		if ev.Rune() == 'g' {
			return ActionTop
		}
		return ih.processKey(ev)
	}
	return ActionNone
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
			ih.timer = time.AfterFunc(500*time.Millisecond, func() {
				ih.state = stateIdle
			})
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

func (ih *InputHandler) cancelTimer() {
	if ih.timer != nil {
		ih.timer.Stop()
		ih.timer = nil
	}
}
