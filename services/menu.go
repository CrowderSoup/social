package services

import (
	"github.com/jinzhu/gorm"

	"github.com/CrowderSoup/socialboat/models"
)

// MenuService a service to get the menus
type MenuService struct {
	DB *gorm.DB
}

// NewMenuService returns a new menu service
func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{
		DB: db,
	}
}

// Create creates a menu
func (s *MenuService) Create(menu *models.Menu) error {
	if err := s.DB.Create(menu).Error; err != nil {
		return err
	}

	return nil
}

// CreateItem creates a menu item
func (s *MenuService) CreateItem(menuItem *models.MenuItem) error {
	if err := s.DB.Create(menuItem).Error; err != nil {
		return err
	}

	return nil
}

// Update updates a menu
func (s *MenuService) Update(menu *models.Menu) error {
	if err := s.DB.Save(menu).Error; err != nil {
		return err
	}

	return nil
}

// UpdateItem updates a menu item
func (s *MenuService) UpdateItem(menuItem *models.MenuItem) error {
	if err := s.DB.Save(menuItem).Error; err != nil {
		return err
	}

	return nil
}

// GetAll get's all the menus
func (s *MenuService) GetAll() (map[string]models.Menu, error) {
	var menus []models.Menu
	if err := s.DB.Preload("MenuItems", func(db *gorm.DB) *gorm.DB {
		return db.Order("menu_items.weight")
	}).Find(&menus).Error; err != nil {
		return nil, err
	}

	result := make(map[string]models.Menu)
	for _, m := range menus {
		result[m.Name] = m
	}

	return result, nil
}

// GetAllForView get's all menus in ViewMenus form
func (s *MenuService) GetAllForView() (*models.ViewMenus, error) {
	menus, err := s.GetAll()
	if err != nil {
		return nil, err
	}

	return models.NewViewMenus(menus), nil
}

// Find attempts to find a menu and it's items by the menu id
func (s *MenuService) Find(id uint) (models.Menu, error) {
	var menu models.Menu
	if err := s.DB.Preload("MenuItems", func(db *gorm.DB) *gorm.DB {
		return db.Order("menu_items.weight")
	}).First(&menu, id).Error; err != nil {
		return menu, err
	}

	return menu, nil
}

// FindItem attempts to find a menu item by it's ID
func (s *MenuService) FindItem(id uint) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := s.DB.First(&menuItem, id).Error; err != nil {
		return nil, err
	}

	return &menuItem, nil
}

// Delete delete's the menu
func (s *MenuService) Delete(id uint) error {
	menu, err := s.Find(id)
	if err != nil {
		return err
	}

	if err := s.DB.Delete(menu).Error; err != nil {
		return err
	}

	return nil
}

// DeleteItem delete's the menu item
func (s *MenuService) DeleteItem(id uint) error {
	menuItem, err := s.FindItem(id)
	if err != nil {
		return err
	}

	if err := s.DB.Delete(menuItem).Error; err != nil {
		return err
	}

	return nil
}
