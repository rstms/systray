//go:build windows

package menu

import (
	_ "embed"
	"fmt"
	"github.com/rstms/systray"
	"log"
	"runtime"
)

//go:embed icon.ico
var DefaultIconData []byte

type SystrayMenuItem struct {
	item *systray.MenuItem
}

func NewSystrayMenuItem(item *systray.MenuItem) *SystrayMenuItem {
	return &SystrayMenuItem{item: item}
}

func (s *SystrayMenuItem) Clicked() chan struct{} {
	return s.item.ClickedCh
}

func (s *SystrayMenuItem) AddSubMenuItem(title, tooltip string) *SystrayMenuItem {
	return &SystrayMenuItem{item: s.item.AddSubMenuItem(title, tooltip)}
}

func (s *SystrayMenuItem) AddSubMenuItemCheckbox(title, tooltip string, checked bool) *SystrayMenuItem {
	return &SystrayMenuItem{item: s.item.AddSubMenuItemCheckbox(title, tooltip, checked)}
}

func (m *Menu) startup() error {
	if m.debug {
		log.Println("Menu.startup")
	}
	if m.started {
		return fmt.Errorf("already started")
	}
	m.wg.Add(1)
	go func() {
		if m.debug {
			log.Println("Menu.EventLoop started")
			defer log.Println("Menu.EventLoop exiting")
		}
		defer m.wg.Done()
		runtime.LockOSThread()
		systray.Run(m.onReady, m.onExit)
	}()
	m.started = true
	return nil
}

func (m *Menu) shutdown() error {
	if m.debug {
		log.Println("Menu.shutdown")
	}
	if !m.started {
		return fmt.Errorf("never started")
	}
	if m.stopped {
		return fmt.Errorf("already stoped")
	}
	for _, item := range m.items {
		item.stop()
	}
	m.stopped = true
	m.exitHandler <- struct{}{}
	return nil
}

func (m *Menu) onReady() {
	if m.debug {
		log.Println("Menu.onReady")
	}
	// Set the icon and tooltip
	systray.SetTitle(m.Title)
	systray.SetTooltip(m.Title)
	systray.SetIcon(m.iconData)

	if m.qid < 0 {
		m.AddQuitItem("Quit", "Shutdown "+m.Title)
	}

	for _, item := range m.items {
		switch item.Type {
		case MenuItemClickable, MenuItemQuit:
			item.start(NewSystrayMenuItem(systray.AddMenuItem(item.Title, item.Tooltip)))
		case MenuItemCheckbox:
			item.start(NewSystrayMenuItem(systray.AddMenuItemCheckbox(item.Title, item.Tooltip, item.checked)))
		case MenuItemSeparator:
			systray.AddSeparator()
		default:
			panic(fmt.Sprintf("unexpected item: %v", item))
		}
	}
}

func (m *Menu) onExit() {
	if m.debug {
		log.Println("Menu.onExit")
	}
	m.Stop()
}

func (i *MenuItem) start(trayItem *SystrayMenuItem) {
	if i.menu.debug {
		log.Printf("MenuItem.start %d %s\n", i.Id, i.Title)
	}
	i.trayItem = trayItem
	i.menu.wg.Add(1)
	go i.handler()
	for _, subItem := range i.subItems {
		switch subItem.Type {
		case MenuItemClickable, MenuItemQuit:
			subItem.start(i.trayItem.AddSubMenuItem(subItem.Title, subItem.Tooltip))
		case MenuItemCheckbox:
			subItem.start(i.trayItem.AddSubMenuItemCheckbox(subItem.Title, subItem.Tooltip, subItem.checked))
		case MenuItemSeparator:
			systray.AddSeparator()
		default:
			panic(fmt.Sprintf("unexpected subItem: %v", subItem))
		}
	}
}

func (i *MenuItem) handler() {
	defer i.menu.wg.Done()
	defer systray.Quit()
	if i.menu.debug {
		defer log.Printf("MenuItem handler exit %d %s\n", i.Id, i.Title)
		log.Printf("MenuItem handler start %d %s\n", i.Id, i.Title)
	}
	for {
		select {
		case <-i.trayItem.Clicked():
			if i.menu.debug {
				log.Printf("received ClickedCh:  %d %s\n", i.Id, i.Title)
			}
			i.menu.clickMux <- i
			if i.Id == i.menu.qid {
				if i.menu.debug {
					log.Println("Quit Item clicked; calling systray.Quit()")
				}
				i.exitHandler <- struct{}{}
			}
		case <-i.exitHandler:
			return
		}
	}
}

func (i *MenuItem) stop() {
	if i.menu.debug {
		log.Printf("MenuItem.stop %d %s\n", i.Id, i.Title)
	}
	i.exitHandler <- struct{}{}
	for _, subItem := range i.subItems {
		subItem.stop()
	}
}
