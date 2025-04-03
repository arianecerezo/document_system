package main

import (
	"github.com/arianecerezo/document_system/models"
	"github.com/uadmin/uadmin"
)

func main() {
	uadmin.Register(
		models.Folder{},
		models.FolderGroup{},
		models.FolderUser{},
		models.Channel{},
		models.Document{},
		models.DocumentGroup{},
		models.DocumentUser{},
		models.DocumentVersion{},
	)
	uadmin.RegisterInlines(
		models.Folder{},
		map[string]string{
			"foldergroup": "FolderID",
			"folderuser":  "FolderID",
		},
	)

	uadmin.RegisterInlines(
		models.Document{},
		map[string]string{
			"documentgroup":   "DocumentID",
			"documentuser":    "DocumentID",
			"documentversion": "DocumentID",
		},
	)
	docS := uadmin.Schema["document"]

	docS.FormModifier = func(s *uadmin.ModelSchema, m interface{}, u *uadmin.User) {
		d, _ := m.(*models.Document)

		if !u.Admin && d.CreatedBy != "" {
			s.FieldByName("CreatedBy").ReadOnly = "true"
		}
	}
	//////
	docS.ListModifier = func(s *uadmin.ModelSchema, u *uadmin.User) (string, []interface{}) {
		if !u.Admin {
			return "user_id = ?", []interface{}{u.ID}
		}
		return "", []interface{}{}
	}

	uadmin.Schema["document"] = docS
	uadmin.SiteName = "Document System"

	uadmin.StartServer()
}
