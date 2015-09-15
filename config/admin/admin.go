package admin

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor"
	"github.com/qor/qor-example/app/models"
	"github.com/qor/qor-example/config"
	"github.com/qor/qor-example/db"
	"github.com/qor/qor/admin"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/qor/validations"
)

var Admin *admin.Admin
var Countries = []string{"China", "Japan", "USA"}

func init() {
	Admin = admin.New(&qor.Config{DB: db.Publish.DraftDB()})
	Admin.SetAuth(Auth{})

	assetManager := Admin.AddResource(&admin.AssetManager{}, &admin.Config{Invisible: true})

	product := Admin.AddResource(&models.Product{}, &admin.Config{Menu: []string{"Product Management"}})
	product.Meta(&admin.Meta{Name: "MadeCountry", Type: "select_one", Collection: Countries})
	product.Meta(&admin.Meta{Name: "Description", Type: "rich_editor", Resource: assetManager})

	for _, country := range Countries {
		var country = country
		product.Scope(&admin.Scope{Name: country, Group: "Made Country", Handle: func(db *gorm.DB, ctx *qor.Context) *gorm.DB {
			return db.Where("made_country = ?", country)
		}})
	}

	product.IndexAttrs("-ColorVariations")

	Admin.AddResource(&models.Color{}, &admin.Config{Menu: []string{"Product Management"}})
	Admin.AddResource(&models.Size{}, &admin.Config{Menu: []string{"Product Management"}})
	Admin.AddResource(&models.Category{}, &admin.Config{Menu: []string{"Product Management"}})

	Admin.AddResource(&models.Order{}, &admin.Config{Menu: []string{"Order Management"}})

	store := Admin.AddResource(&models.Store{}, &admin.Config{Menu: []string{"Store Management"}})
	store.IndexAttrs("-Latitude", "-Longitude")
	store.AddValidator(func(record interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
		if meta := metaValues.Get("Name"); meta != nil {
			if name := utils.ToString(meta.Value); strings.TrimSpace(name) == "" {
				return validations.NewError(record, "Name", "Name can't be blank")
			}
		}
		return nil
	})

	Admin.AddResource(config.Config.I18n, &admin.Config{Menu: []string{"Site Management"}})

	setting := Admin.AddResource(&models.Setting{}, &admin.Config{Singleton: true})
	setting.UseTheme("location")
	setting.Meta(&admin.Meta{Name: "StoreAddress", Type: "single_edit"})
	setting.Meta(&admin.Meta{Name: "CustomerSupportAddress", Type: "single_edit"})

	user := Admin.AddResource(&models.User{})
	user.IndexAttrs("ID", "Email", "Name", "Gender", "Role")

	Admin.AddResource(db.Publish)
}