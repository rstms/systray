//go:build openbsd

// FIXME: test if this works on openbsd with the new upstream

package menu

import (
	_ "embed"
)

//go:embed icon.png
var DefaultIconData []byte

type SystrayMenuItem struct {
	Title   string
	Tooltip string
}

func (m *Menu) startup() error {
	return nil
}

func (m *Menu) shutdown() error {
	return nil
}

/*
func (m *MenuItem) start() {
}

func (m *MenuItem) stop() {
}
*/
