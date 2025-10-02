package menu

import (
	"fmt"
	"log"
	"sync"
)

const Version = "0.0.15"

type MenuItemType int

const (
	MenuItemClickable MenuItemType = iota
	MenuItemCheckbox
	MenuItemSeparator
	MenuItemQuit
)

type MenuItem struct {
	Id          int
	Type        MenuItemType
	Title       string
	Tooltip     string
	checked     bool
	trayItem    *SystrayMenuItem
	exitHandler chan struct{}
	menu        *Menu
	subItems    []*MenuItem
}

type Menu struct {
	Title       string
	Tooltip     string
	iconData    []byte
	Clicked     chan *MenuItem // caller <- menu when items clicked
	Exited      chan struct{}  // caller <- menu when menu has exited
	clickMux    chan *MenuItem
	exitHandler chan struct{}
	items       []*MenuItem
	nextId      int
	started     bool
	stopped     bool
	qid         int
	wg          sync.WaitGroup
	debug       bool
}

func NewMenu(title, tooltip string, iconData []byte, clicked chan *MenuItem, exited chan struct{}) *Menu {
	m := Menu{
		Title:       title,
		Tooltip:     tooltip,
		Clicked:     clicked,
		Exited:      exited,
		iconData:    iconData,
		clickMux:    make(chan *MenuItem, 1),
		exitHandler: make(chan struct{}, 1),
		items:       []*MenuItem{},
		qid:         -1,
		debug:       false,
	}
	if len(iconData) == 0 {
		m.iconData = DefaultIconData
	}
	m.AddItem(title, tooltip)
	m.AddSeparator()
	return &m
}

func (m *Menu) Start() error {
	if m.debug {
		log.Printf("Menu.Start: started=%v stopped=%v\n", m.started, m.stopped)
	}
	if m.started {
		return fmt.Errorf("already started")
	}
	m.wg.Add(1)
	go m.handler()
	err := m.startup()
	if err != nil {
		return err
	}
	if m.debug {
		log.Printf("Menu.Start.Returning: started=%v stopped=%v\n", m.started, m.stopped)
	}
	return nil
}

func (m *Menu) Run() error {
	if m.debug {
		log.Printf("Menu.Run: started=%v stopped=%v\n", m.started, m.stopped)
	}
	err := m.Start()
	if err != nil {
		return err
	}
	err = m.Wait()
	if err != nil {
		return err
	}
	if m.debug {
		log.Printf("Menu.Run.Returning: started=%v stopped=%v\n", m.started, m.stopped)
	}
	return nil
}

func (m *Menu) handler() {
	if m.debug {
		log.Println("Menu.handler start")
		defer log.Println("Menu.handler exit")
	}
	defer m.wg.Done()
	for {
		select {
		case item := <-m.clickMux:
			if m.debug {
				log.Printf("Menu.handler read from mux: %d %s\n", (*item).Id, (*item).Title)
			}
			if m.Clicked != nil {
				m.Clicked <- item
			}
		case <-m.exitHandler:
			if m.debug {
				log.Println("Menu.handler read from shutdown")
			}
			if m.Exited != nil {
				m.Exited <- struct{}{}
			}
			return
		}
	}
}

func (m *Menu) Stop() error {
	if m.debug {
		log.Printf("Menu.Stop: started=%v stopped=%v\n", m.started, m.stopped)
	}
	err := m.shutdown()
	if err != nil {
		return err
	}
	if m.debug {
		log.Printf("Menu.Stop.Returning: started=%v stopped=%v\n", m.started, m.stopped)
	}
	return nil
}

func (m *Menu) Wait() error {
	if m.debug {
		log.Printf("Menu.Wait: started=%v stopped=%v\n", m.started, m.stopped)
	}
	m.wg.Wait()
	if m.debug {
		log.Printf("Menu.Wait.Returning: started=%v stopped=%v\n", m.started, m.stopped)
	}
	return nil
}

func (m *Menu) nextItem(itemType MenuItemType, title, tooltip string, checked bool) *MenuItem {
	item := MenuItem{
		Id:          m.nextId,
		Title:       title,
		Tooltip:     tooltip,
		Type:        itemType,
		checked:     checked,
		subItems:    []*MenuItem{},
		menu:        m,
		exitHandler: make(chan struct{}, 1),
	}
	m.nextId++
	return &item
}

func (m *Menu) AddItem(title, tooltip string) *MenuItem {
	item := m.nextItem(MenuItemClickable, title, tooltip, false)
	m.items = append(m.items, item)
	return item
}

func (m *Menu) AddQuitItem(title, tooltip string) *MenuItem {
	item := m.nextItem(MenuItemQuit, title, tooltip, false)
	m.qid = (*item).Id
	m.items = append(m.items, item)
	return item
}

func (m *Menu) AddCheckboxItem(title, tooltip string, checked bool) *MenuItem {
	item := m.nextItem(MenuItemCheckbox, title, tooltip, checked)
	m.items = append(m.items, item)
	return item
}

func (m *Menu) AddSeparator() {
	item := m.nextItem(MenuItemSeparator, "", "", false)
	m.items = append(m.items, item)
}

func (i *MenuItem) AddItem(title, tooltip string) *MenuItem {
	item := i.menu.nextItem(MenuItemClickable, title, tooltip, false)
	i.subItems = append(i.subItems, item)
	return item
}

func (i *MenuItem) AddQuitItem(title, tooltip string) *MenuItem {
	item := i.menu.nextItem(MenuItemQuit, title, tooltip, false)
	i.menu.qid = (*item).Id
	i.subItems = append(i.subItems, item)
	return item
}

func (i *MenuItem) AddCheckboxItem(title, tooltip string, checked bool) *MenuItem {
	item := i.menu.nextItem(MenuItemCheckbox, title, tooltip, checked)
	i.subItems = append(i.subItems, item)
	return item
}

func (i *MenuItem) AddSeparator() {
	item := i.menu.nextItem(MenuItemSeparator, "", "", false)
	i.subItems = append(i.subItems, item)
}
