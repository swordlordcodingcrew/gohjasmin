package tablemodule

/*-----------------------------------------------------------------------------
 **
 ** - GohJasmin -
 **
 ** Copyright 2017-20 by SwordLord - the coding crew - http://www.swordlord.com
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
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gohjasmin"
	"gohjasmin/model"
	"io"
	"log"
	"strconv"
)

func getDomains() [][]string {

	db := gohjasmin.GetDB()

	var rows []*model.Domain

	db.Order("domain").Find(&rows)

	// Create
	//db.Create(&model.Domain{Domain: "demo", Pwd: "demo"})

	//db.First(&domain, "name = ?", "demo") // find product with id 1

	var domains [][]string

	for _, domain := range rows {

		domains = append(domains, []string{domain.Domain, domain.IP, domain.CrtDat.Format("2006-01-02 15:04:05"), domain.UpdDat.Format("2006-01-02 15:04:05")})
	}

	return domains
}

func ListDomains() {

	domains := getDomains()

	gohjasmin.WriteTable([]string{"Domain", "IP", "CrtDat", "UpdDat"}, domains)
}

// write list of domains to w (usually either os.Stdout or a file
func ExportDomains(w io.Writer, ttl int) {

	domains := getDomains()

	for _, domain := range domains {

		ttlS := strconv.Itoa(ttl)

		fmt.Fprintf(w, "+%s:%s:%s\n", domain[0], domain[1], ttlS)
	}
}

func AddDomain(domain string, ip string) {

	db := gohjasmin.GetDB()

	retDB := db.Create(&model.Domain{Domain: domain, IP: ip})

	if retDB.Error != nil {
		log.Printf("Error with Domain %q: %s\n", domain, retDB.Error)
		log.Fatal(retDB.Error)
		return
	}

	fmt.Printf("Domain %s added.\n", domain)
}

func UpdateDomain(domain string, ip string) (int, error) {

	db := gohjasmin.GetDB()

	retDB := db.Model(&model.Domain{}).Where("domain=? AND ip<>?", domain, ip).Update("ip", ip)

	if retDB.Error != nil {
		log.Printf("Error with Domain %q: %s\n", domain, retDB.Error)
		return 0, retDB.Error
	}

	fmt.Printf("Domain %s updated with IP %s.\n", domain, ip)

	return int(retDB.RowsAffected), nil
}

func DeleteDomain(domain string) {

	db := gohjasmin.GetDB()

	d := &model.Domain{}

	retDB := db.Where("name = ?", domain).First(&d)

	if retDB.Error != nil {
		log.Printf("Error with Domain %q: %s\n", domain, retDB.Error)
		log.Fatal(retDB.Error)
		return
	}

	if retDB.RowsAffected <= 0 {
		log.Printf("Domain not found: %s\n", domain)
		log.Fatal("Domain not found: " + domain + "\n")
		return
	}

	log.Printf("Deleting Domain: %s", &d.Domain)

	db.Delete(&d)

	fmt.Printf("Domain %s deleted.\n", domain)
}
