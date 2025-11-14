package menu

import (
	"fmt"
	"strconv"

	"strings"
	"time"
)

func (menu *Menu) GetWidgetList(limit int, offset int, filter Filter, tenantid string, websiteid int) ([]TblWidgets, int, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return []TblWidgets{}, 0, AuthError
	}

	menumodel.DataAccess = menu.DataAccess
	menumodel.Userid = menu.UserId

	_, totalcount, _ := menumodel.WidgetList(0, 0, filter, menu.DB, tenantid, websiteid)

	widgetlist, _, cerr := menumodel.WidgetList(limit, offset, filter, menu.DB, tenantid, websiteid)

	if cerr != nil {

		return []TblWidgets{}, 0, cerr
	}

	return widgetlist, int(totalcount), nil
}

// CreateWidget
func (menu *Menu) CreateWidget(widget *TblWidgets) (TblWidgets, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblWidgets{}, AuthError
	}

	widget.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	widgetdetail, err := menumodel.CreateWidget(menu.DB, widget)

	if err != nil {

		return TblWidgets{}, err

	}

	if widget.ProductIds != "" {

		newids := strings.Split(widget.ProductIds, ",")

		var intids []int
		for _, idStr := range newids {

			id, _ := strconv.Atoi(idStr)

			intids = append(intids, id)
		}
		createdon, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
		for _, val := range intids {

			widgerproduct := TblWidgetProducts{

				WidgetId:  widget.Id,
				ProductId: val,
				CreatedBy: widget.CreatedBy,
				CreatedOn: createdon,
				TenantId:  widget.TenantId,
			}
			err := menumodel.InsertWidgetProductIds(menu.DB, &widgerproduct)

			if err != nil {

				fmt.Println(err)
			}
		}
	}
	return widgetdetail, nil

}

func (menu *Menu) GetWidgetById(widgetid int, tenantid string) (TblWidgets, []TblWidgetProducts, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblWidgets{}, []TblWidgetProducts{}, AuthError
	}
	widgetdetail, product, err := menumodel.GetWidgetById(menu.DB, widgetid, tenantid)

	if err != nil {

		return TblWidgets{}, []TblWidgetProducts{}, err

	}
	return widgetdetail, product, nil
}

// UpdateWidget
func (menu *Menu) UpdateWidget(widget *TblWidgets, widgetid int) (TblWidgets, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblWidgets{}, AuthError
	}

	widget.Id = widgetid
	widget.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	widgetdetail, err := menumodel.UpdateWidget(menu.DB, widget)

	derr := menumodel.DeleteProductIds(menu.DB, widgetid,widget.TenantId)

	if derr != nil {
		fmt.Println(derr)
	}

	if widget.ProductIds != "" {

		newids := strings.Split(widget.ProductIds, ",")

		var intids []int
		for _, idStr := range newids {

			id, _ := strconv.Atoi(idStr)

			intids = append(intids, id)
		}
		createdon, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
		for _, val := range intids {

			widgerproduct := TblWidgetProducts{

				WidgetId:  widget.Id,
				ProductId: val,
				CreatedBy: widget.CreatedBy,
				CreatedOn: createdon,
				TenantId:  widget.TenantId,
			}
			err := menumodel.InsertWidgetProductIds(menu.DB, &widgerproduct)

			if err != nil {

				fmt.Println(err)
			}
		}
	}

	if err != nil {

		return TblWidgets{}, err

	}
	return widgetdetail, nil

}

//Delete widget

func (menu *Menu) DeleteWidgetById(widgetid int, modifiedby int, tenantid string) error {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return AuthError
	}

	var widgetdetail TblWidgets

	widgetdetail.DeletedBy = modifiedby

	widgetdetail.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	widgetdetail.IsDeleted = 1

	widgetdetail.Id = widgetid

	err := menumodel.DeleteWidgetById(&widgetdetail, tenantid, menu.DB)

	if err != nil {

		return err
	}

	return nil

}

// Status change
func (menu *Menu) WidgetStatusChange(widgetid int, status int, userid int, tenantid string) (bool, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return false, AuthError
	}

	var widgetdetail TblWidgets

	widgetdetail.ModifiedBy = userid

	widgetdetail.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	widgetdetail.Status = status

	widgetdetail.TenantId = tenantid

	widgetdetail.Id = widgetid

	err := menumodel.WidgetStatusChange(widgetdetail, menu.DB)

	if err != nil {

		return false, err

	}
	return true, nil
}

// GetwidgetbySlug
func (menu *Menu) GetWidgetBySlug(slug string, tenantid string) (TblWidgets, error) {

	if AuthError := AuthandPermission(menu); AuthError != nil {

		return TblWidgets{}, AuthError
	}
	widgetdetail, err := menumodel.GetWidgetBySlug(menu.DB, slug, tenantid)

	if err != nil {

		return TblWidgets{}, err

	}
	return widgetdetail, nil

}
