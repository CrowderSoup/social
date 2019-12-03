package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"

	"github.com/CrowderSoup/socialboat/models"
	"github.com/CrowderSoup/socialboat/services"
)

// MenuController handlers for menu routes
type MenuController struct {
	MenuService *services.MenuService
}

// NewMenuController returns a new menu controller
func NewMenuController(db *gorm.DB) *MenuController {
	return &MenuController{
		MenuService: services.NewMenuService(db),
	}
}

// InitRoutes initialize routes for this controller
func (c *MenuController) InitRoutes(g *echo.Group) {
	g.GET("", c.list)
	g.GET("/:menu_id", c.get)
	g.POST("/create", c.create)
	g.POST("/:menu_id/item/create", c.createItem)
	g.POST("/:menu_id", c.update)
	g.POST("/:menu_id/item/update", c.updateItem)
	g.DELETE("/:menu_id/delete", c.delete)
	g.DELETE("/:menu_id/item/:item_id", c.deleteItem)
}

func (c *MenuController) list(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	// Ensure the user is logged in
	err := bc.EnsureLoggedIn()
	if err != nil {
		return err
	}

	menus, err := c.MenuService.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get menus")
	}

	return bc.ReturnView(http.StatusOK, "menus/index", echo.Map{
		"menus": menus,
	})
}

func (c *MenuController) get(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	// Ensure the user is logged in
	err := bc.EnsureLoggedIn()
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(bc.Param("menu_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid menu id")
	}

	menu, err := c.MenuService.Find(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't find menu")
	}

	return bc.ReturnView(http.StatusOK, "menus/menu", echo.Map{
		"menu": menu,
	})
}

func (c *MenuController) create(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	// Ensure the user is logged in
	err := bc.EnsureLoggedIn()
	if err != nil {
		return err
	}

	name := strings.TrimSpace(bc.FormValue("name"))

	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Menu Name is required!")
	}

	menu := models.Menu{
		Name: name,
	}

	err = c.MenuService.Create(&menu)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create menu")
	}

	menus, err := c.MenuService.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get menus")
	}

	return bc.ReturnView(http.StatusOK, "menus/index", echo.Map{
		"menus": menus,
	})
}

func (c *MenuController) createItem(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	// Ensure the user is logged in
	err := bc.EnsureLoggedIn()
	if err != nil {
		return err
	}

	name := strings.TrimSpace(bc.FormValue("item_name"))
	URL := strings.TrimSpace(bc.FormValue("item_url"))

	if name == "" || URL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name and url are required")
	}

	weight, err := strconv.ParseUint(strings.TrimSpace(bc.FormValue("item_weight")), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "weight must be a number")
	}

	menuID, err := strconv.ParseUint(bc.Param("menu_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "couldn't parse menu_id")
	}

	menuItem := models.MenuItem{
		Name:   name,
		URL:    URL,
		Weight: uint(weight),
		MenuID: uint(menuID),
	}

	err = c.MenuService.CreateItem(&menuItem)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't create menu item")
	}

	return bc.Redirect(http.StatusSeeOther, fmt.Sprintf("/menus/%d", menuID))
}

func (c *MenuController) update(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	// Ensure the user is logged in
	err := bc.EnsureLoggedIn()
	if err != nil {
		return err
	}

	menuID, err := strconv.ParseUint(bc.Param("menu_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "couldn't parse menu_id")
	}

	menu, err := c.MenuService.Find(uint(menuID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get menu")
	}

	menuName := strings.TrimSpace(bc.FormValue("menu_name"))
	if menuName != "" {
		menu.Name = menuName
	}

	err = c.MenuService.Update(&menu)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't update menu")
	}

	menus, err := c.MenuService.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get menus")
	}

	return bc.ReturnView(http.StatusOK, "menus/index", echo.Map{
		"menus": menus,
	})
}

func (c *MenuController) updateItem(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	// Ensure the user is logged in
	err := bc.EnsureLoggedIn()
	if err != nil {
		return err
	}

	itemID, err := strconv.ParseUint(bc.FormValue("item_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "couldn't parse menu item id")
	}

	menuItem, err := c.MenuService.FindItem(uint(itemID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't find menu item")
	}

	name := strings.TrimSpace(bc.FormValue("item_name"))
	URL := strings.TrimSpace(bc.FormValue("item_url"))

	if name == "" || URL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name and url are required")
	}

	weight, err := strconv.ParseUint(strings.TrimSpace(bc.FormValue("item_weight")), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "weight must be a number")
	}

	menuItem.Name = name
	menuItem.URL = URL
	menuItem.Weight = uint(weight)

	err = c.MenuService.UpdateItem(menuItem)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update menu item")
	}

	menuID, err := strconv.ParseUint(bc.Param("menu_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "couldn't parse menu_id")
	}

	return bc.Redirect(http.StatusSeeOther, fmt.Sprintf("/menus/%d", menuID))
}

func (c *MenuController) delete(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	// Ensure the user is logged in
	err := bc.EnsureLoggedIn()
	if err != nil {
		return err
	}

	menuID, err := strconv.ParseUint(bc.Param("menu_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "couldn't parse menu_id")
	}

	err = c.MenuService.Delete(uint(menuID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't delete menu")
	}

	return nil
}

func (c *MenuController) deleteItem(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	// Ensure the user is logged in
	err := bc.EnsureLoggedIn()
	if err != nil {
		return err
	}

	itemID, err := strconv.ParseUint(bc.Param("item_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "couldn't parse menu item id")
	}

	err = c.MenuService.DeleteItem(uint(itemID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't delete menu item")
	}

	return nil
}
