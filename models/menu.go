package models

import "github.com/jinzhu/gorm"

// Menu a menu to be used on the front-end
type Menu struct {
	gorm.Model

	Name      string
	MenuItems []MenuItem
}

// MenuItem is an item for a menu
type MenuItem struct {
	gorm.Model

	Name   string
	URL    string
	Weight uint
	MenuID uint
}

// ViewMenus data structure for menus in views
type ViewMenus struct {
	menus map[string]Menu
}

// NewViewMenus returns a new ViewMenus
func NewViewMenus(menus map[string]Menu) *ViewMenus {
	return &ViewMenus{
		menus: menus,
	}
}

// GetMenu returns a menu for a view
func (vm *ViewMenus) GetMenu(name string) Menu {
	return vm.menus[name]
}
