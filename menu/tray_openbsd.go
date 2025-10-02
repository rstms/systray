//go:build openbsd

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
