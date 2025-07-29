# Menu Package

The Menu module in SpurtCMS manages site navigation structures like headers, footers, and sidebars. It allows creation of hierarchical menus with nested items and custom ordering. 

## Features

- Dynamic Menu Creation – Create and manage multiple menus (e.g., header, footer, sidebar) with nested structures. 
- Nested Menu Items – Support for parent-child hierarchy to build multi-level navigation.
- Drag-and-Drop Ordering – Reorder menu items easily via UI or APIs for precise control.
- Multilingual Support – Localize menu labels for different languages and regions.
- Status Toggle – Enable or disable menu items without deleting them.

# Installation

``` bash
go get github.com/spurtcms/menu
```


# Usage Example


``` bash

import (
	"github.com/spurtcms/auth"
	"github.com/spurtcms/menu"
)

func main() {

	Auth := auth.AuthSetup(auth.Config{
		UserId:     1,
		ExpiryTime: 2,
		SecretKey:  "SecretKey@123",
		DB:         &gorm.DB{},
		RoleId:     1,
	})

	token, _ := Auth.CreateToken()

	Auth.VerifyToken(token, SecretKey)

	permission, _ := NewAuth.IsGranted("Menu", auth.CRUD, TenantId)

	MenuConfig = menu.MenuSetup(menu.Config{

		DB:               DB,
		AuthEnable:       true,
		PermissionEnable: true,
		Auth:             NewAuth,
	})

	//Menus

	if permisison {

		//Menu List

		Menulist, Total_count, err := MenuConfig.MenuList(10, 0, filter{}, 1)

		if err != nil {
			ErrorLog.Printf("menu list  error: %s", err)
		}

		//create Menu

		Create := menu.MenuCreate{
			MenuName:    menuname,
			Description: menudesc,
			Status:      menustatus,
			TenantId:    1,
			CreatedBy:   1,
			ParentId:    0,
		}

		_, err := MenuConfig.CreateMenus(Create)

		//update menu
		menudetails := menu.MenuCreate{
			MenuName:    c.PostForm("menu_name"),
			Description: c.PostForm("menu_desc"),
			Status:      menustatus,
			ModifiedBy:  1,
			TenantId:    1,
			Id:          1,
		}
		_, err := MenuConfig.UpdateMenu(menudetails)

		// delete Course
		err := MenuConfig.DeleteMenu(1, 1, 1)

	}
}


```

# Getting help
If you encounter a problem with the package,please refer [Please refer [(https://www.spurtcms.com/documentation/cms-admin)] or you can create a new Issue in this repo[https://github.com/spurtcms/menu/issues]. 
