package models

import (
	"fmt"
	"strings"

	"github.com/uadmin/uadmin"
)

// DocumentUser !
type DocumentUser struct {
	uadmin.Model
	User       uadmin.User
	UserID     uint
	Document   Document
	DocumentID uint
	Read       bool
	Add        bool
	Edit       bool
	Delete     bool
}

func (d *DocumentUser) String() string {

	uadmin.Preload(d)

	return d.User.String()
}

func (d Document) GetPermissions(user uadmin.User) (Read bool, Add bool, Edit bool, Delete bool) {
	if user.Admin {
		Read = true
		Add = true
		Edit = true
		Delete = true
	}

	if d.FolderID != 0 {
		folderGroup := FolderGroup{}

		uadmin.Get(&folderGroup, "group_id = ? AND folder_id = ?", user.UserGroupID, d.FolderID)

		if folderGroup.ID != 0 {
			Read = folderGroup.Read
			Add = folderGroup.Add
			Edit = folderGroup.Edit
			Delete = folderGroup.Delete
		}

		folderUser := FolderUser{}

		uadmin.Get(&folderUser, "user_id = ? AND folder_id = ?", user.ID, d.FolderID)

		if folderUser.ID != 0 {
			Read = folderUser.Read
			Add = folderUser.Add
			Edit = folderUser.Edit
			Delete = folderUser.Delete
		}
	}
	documentGroup := DocumentGroup{}

	uadmin.Get(&documentGroup, "group_id = ? AND document_id = ?", user.UserGroupID, d.ID)

	if documentGroup.ID != 0 {
		Read = documentGroup.Read
		Add = documentGroup.Add
		Edit = documentGroup.Edit
		Delete = documentGroup.Delete
	}

	documentUser := DocumentUser{}

	uadmin.Get(&documentUser, "user_id = ? AND document_id = ?", user.ID, d.ID)

	if documentUser.ID != 0 {
		Read = documentUser.Read
		Add = documentUser.Add
		Edit = documentUser.Edit
		Delete = documentUser.Delete
	}

	return
}

// Count !
func (d Document) Count(a interface{}, query interface{}, args ...interface{}) int {
	Q := fmt.Sprint(query)
	if strings.Contains(Q, "user_id = ?") {
		qParts := strings.Split(Q, " AND ")
		tempArgs := []interface{}{}
		tempQuery := []string{}
		for i := range qParts {
			if qParts[i] != "user_id = ?" {
				tempArgs = append(tempArgs, args[i])

				tempQuery = append(tempQuery, qParts[i])
			}
		}
		query = strings.Join(tempQuery, " AND ")
		args = tempArgs
	}
	return uadmin.Count(a, query, args...)
}

func (d Document) AdminPage(order string, asc bool, offset int, limit int, a interface{}, query interface{}, args ...interface{}) (err error) {
	if offset < 0 {
		offset = 0
	}

	userID := uint(0)

	Q := fmt.Sprint(query)

	if strings.Contains(Q, "user_id = ?") {
		uadmin.Trail(uadmin.DEBUG, "1")

		qParts := strings.Split(Q, " AND ")

		tempArgs := []interface{}{}
		tempQuery := []string{}

		for i := range qParts {
			if qParts[i] != "user_id = ?" {
				tempArgs = append(tempArgs, args[i])

				tempQuery = append(tempQuery, qParts[i])
			} else {
				uadmin.Trail(uadmin.DEBUG, "UserID: %d", args[i])

				userID, _ = (args[i]).(uint)
			}
		}
		query = strings.Join(tempQuery, " AND ")

		args = tempArgs
	}

	if userID == 0 {
		uadmin.Trail(uadmin.DEBUG, "2")

		err = uadmin.AdminPage(order, asc, offset, limit, a, query, args...)

		return err
	}

	user := uadmin.User{}

	uadmin.Get(&user, "id = ?", userID)

	docList := []Document{}
	tempList := []Document{}

	for {
		err = uadmin.AdminPage(order, asc, offset, limit, &tempList, query, args)
		uadmin.Trail(uadmin.DEBUG, "8: offset:%d, limit:%d", offset, limit)

		if err != nil {
			uadmin.Trail(uadmin.DEBUG, "3")

			*a.(*[]Document) = docList

			return err
		}

		if len(tempList) == 0 {
			uadmin.Trail(uadmin.DEBUG, "4")

			*a.(*[]Document) = docList

			uadmin.Trail(uadmin.DEBUG, "a: %#v", a)

			return nil
		}

		for i := range tempList {
			p, _, _, _ := tempList[i].GetPermissions(user)

			if p {
				uadmin.Trail(uadmin.DEBUG, "5")

				docList = append(docList, tempList[i])
			}

			if len(docList) == limit {
				uadmin.Trail(uadmin.DEBUG, "6")

				*a.(*[]Document) = docList

				return nil
			}
		}

		offset += limit
	}
	*a.(*[]Document) = docList

	uadmin.Trail(uadmin.DEBUG, "7")

	return nil
}

// Permissions__Form creates a new field named Permissions !
func (d Document) Permissions__Form() string {
	// Initialize u variable that calls the User model
	u := uadmin.User{}

	// Get the user record based on an assigned ID
	uadmin.Get(&u, "id = ?", 1)

	// Initialize read, add, edit and delete that gets the permission for a
	// specific user based on an assigned ID
	r, a, e, del := d.GetPermissions(u)

	// Returns the permission status
	return fmt.Sprintf("Read: %v Add: %v, Edit: %v, Delete: %v", r, a, e, del)
}

func (DocumentUser) HideInDashboard() bool {
	return true
}
