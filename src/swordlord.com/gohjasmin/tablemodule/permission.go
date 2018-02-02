package tablemodule
/*-----------------------------------------------------------------------------
 **
 ** - GohJasmin -
 **
 ** Copyright 2017-18 by SwordLord - the coding crew - http://www.swordlord.com
 ** and contributing authors
 **
 ** This program is free software; you can redistribute it and/or modify it
 ** under the terms of the GNU Affero General Public License as published by the
 ** Free Software Foundation, either version 3 of the License, or (at your option)
 ** any later version.
 **
 ** This program is distributed in the hope that it will be useful, but WITHOUT
 ** ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 ** FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
 ** for more details.
 **
 ** You should have received a copy of the GNU Affero General Public License
 ** along with this program. If not, see <http://www.gnu.org/licenses/>.
 **
 **-----------------------------------------------------------------------------
 **
 ** Original Authors:
 ** LordEidi@swordlord.com
 ** LordLightningBolt@swordlord.com
 **
-----------------------------------------------------------------------------*/

import (
	"log"
	"fmt"
	"strconv"
	"swordlord.com/gohjasmin"
	"swordlord.com/gohjasmin/model"
	"regexp"
)

func GetPermissionForUser(name string) []string {

	permissions := []string{ }

	var rows []*model.Permission

	db := gohjasmin.GetDB()

	retDB := db.Where("user = ?", name).Find(&rows)

	if retDB.Error != nil {
		log.Printf("Error with Permission for User %q: %s\n", name, retDB.Error )
		log.Fatal(retDB.Error)
		return permissions
	}

	if retDB.RowsAffected <= 0 {
		log.Printf("Permission for User not found: %s\n", name)
		return permissions
	}

	for _, permission := range rows {

		permissions = append(permissions, permission.Permission)
	}

	return permissions
}

func ListPermission() {

	db := gohjasmin.GetDB()

	var rows []*model.Permission

	db.Order("user").Find(&rows)

	// Create
	//db.Create(&model.Permission{Name: "demo", Pwd: "demo"})

	//db.First(&user, "name = ?", "demo") // find product with id 1

	var permissions [][]string

	for _, p := range rows {

		permissions = append(permissions, []string{ strconv.Itoa(p.PermissionId), p.User, p.Permission, p.CrtDat.Format("2006-01-02 15:04:05"), p.UpdDat.Format("2006-01-02 15:04:05")})
	}

	gohjasmin.WriteTable([]string{"Id", "User", "Permission", "CrtDat", "UpdDat"}, permissions)
}

func AddPermission(user string, permission string) {

	db := gohjasmin.GetDB()

	retDB := db.Create(&model.Permission{Permission: permission, User: user})

	if retDB.Error != nil {
		log.Printf("Error when creating a new permission %s for user %s: %s\n", permission, user, retDB.Error )
		log.Fatal(retDB.Error)
		return
	}

	fmt.Printf("Permission %s for user %s added.\n", permission, user)
}

func DeletePermission(permission string) {

	db := gohjasmin.GetDB()

	p := &model.Permission{}

	retDB := db.Where("permission = ?", permission).First(&p)

	if retDB.Error != nil {
		log.Printf("Error with Permission %q: %s\n", permission, retDB.Error )
		log.Fatal(retDB.Error)
		return
	}

	if retDB.RowsAffected <= 0 {
		log.Printf("User not found: %s\n", permission)
		log.Fatal("User not found: " + permission + "\n")
		return
	}

	log.Printf("Deleting permission: %s", permission)

	db.Delete(&p)

	fmt.Printf("Permission %s deleted.\n", permission)
}

func HasPermission(requestedPermission string, permissions []string) bool {

	for _, p := range permissions {

		rp := regexp.MustCompile(p)

		if rp.MatchString(requestedPermission){
			return true
		}
	}

	return false
}
