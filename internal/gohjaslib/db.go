package gohjaslib

/*-----------------------------------------------------------------------------
 **
 ** - GohJasmin -
 **
 ** Copyright 2017 by SwordLord - the coding crew - http://www.swordlord.com
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
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
)

func GetIPFromDomain(domain string) (string, error) {

	dbFile := GetStringFromConfig("db.file")

	LogDebug("Opening DB file.", logrus.Fields{"DB": dbFile})

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		LogError("Opening DB file fails.", logrus.Fields{"error": err})
		return "", err
	}

	if db == nil {
		LogError("Opening DB file fails.", logrus.Fields{"error": "database is nil"})
		return "", errors.New("DB is nil")
	}

	defer db.Close()

	sql := GetStringFromConfig("db.sql.getipfromdomain")
	if len(sql) == 0 {
		LogError("Query DNS Record fails.", logrus.Fields{"error": "db.sql.getipfromdomain is missing from config file"})
		return "", errors.New("db.sql.getipfromdomain is missing from config file")
	}

	var s []byte
	err = db.QueryRow(sql, domain).Scan(&s)

	if err == nil {

		return string(s), nil
	} else {

		return "", err
	}
}

func UpdateDNSRecord(domain string, ip string) (int, error) {

	dbFile := GetStringFromConfig("db.file")

	LogDebug("Opening DB file.", logrus.Fields{"DB": dbFile})

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		LogError("Opening DB file fails.", logrus.Fields{"error": err})
		return 0, err
	}

	if db == nil {
		LogError("Opening DB file fails.", logrus.Fields{"error": "database is nil"})
		return 0, errors.New("DB is nil")
	}

	defer db.Close()

	sql := GetStringFromConfig("db.sql.update")
	if len(sql) == 0 {
		LogError("Update DNS Record fails.", logrus.Fields{"error": "db.sql.update is missing from config file"})
		return 0, errors.New("db.sql.update is missing from config file")
	}

	stmt, err := db.Prepare(sql)
	if err != nil {
		LogError("Update DNS Record fails when preparing statement.", logrus.Fields{"error": err})
		return 0, err
	}

	LogDebug("Update statement.", logrus.Fields{"statement": stmt})

	defer stmt.Close()

	res, err := stmt.Exec(ip, domain)
	if err != nil {
		LogError("Update DNS Record fails while executing statement.", logrus.Fields{"error": err})
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		LogError("Update DNS Record fails.", logrus.Fields{"error": err})
		return 0, err
	}

	return int(rowsAffected), nil
}
