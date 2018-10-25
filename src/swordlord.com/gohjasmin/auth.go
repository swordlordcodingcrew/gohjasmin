package gohjasmin

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
	"encoding/csv"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
)

// demo.domain.com,demo,$2a$10$v5ZTJXisKF3b83K.XpKoZOYpOoi77qDkq.hLa0b6.8A/sc85aq2F2,127.0.0.1
type record struct {
	Domain     string
	Username   string
	PasswdHash string
	IP         string
}

var auth map[string]record

func LoadAuth() {

	auth = make(map[string]record)

	sFile := GetStringFromConfig("auth.file")

	file, err := os.Open(sFile)
	if err != nil {
		LogFatal("Error reading authentication file", logrus.Fields{"file": sFile, "error": err})
		return
	}

	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ','
	r.Comment = '#'

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			LogFatal("Failing when reading authentication file.", logrus.Fields{"file": sFile, "error": err})
			return
		}

		rec := record{
			Domain:     line[0],
			Username:   line[1],
			PasswdHash: line[2],
			IP:         line[3],
		}

		auth[rec.Username] = rec
	}
}

func ValidateUser(user string, password string) bool {

	//log.Printf("REMOVE ME: User: %s Pwd: %s\n", user, password)

	rec, ok := auth[user]
	if !ok {
		LogInfo("User unknown.", logrus.Fields{"user": user})
		return false
	}

	err := checkHashedPassword(rec.PasswdHash, password)
	if err != nil {
		LogDebug("Passwordhash missmatch for user.", logrus.Fields{"user": user, "error": err})
		return false
	}

	return true
}

func HasIPChanged(user string, IP string) bool {

	rec, ok := auth[user]
	if !ok {
		LogInfo("User unknown.", logrus.Fields{"user": user})
		return false
	}

	LogDebug("Checking if IP has changed.", logrus.Fields{"ip_old": rec.IP, "ip_new": IP})

	hasChanged := rec.IP != IP

	return hasChanged
}

func GetDomainFromUser(user string) string {

	rec, ok := auth[user]
	if !ok {
		LogInfo("User unknown.", logrus.Fields{"user": user})
		return user
	}

	return rec.Domain
}

func ChangeIPInMemory(user string, IP string) {

	rec, ok := auth[user]
	if !ok {
		LogInfo("User unknown.", logrus.Fields{"user": user})
		return
	}

	rec.IP = IP

	auth[user] = rec
}

func hashPassword(pwd string) (string, error) {

	password := []byte(pwd)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func checkHashedPassword(hashedPassword string, password string) error {

	pwd := []byte(password)
	hashedPwd := []byte(hashedPassword)

	// Comparing the password with the hash
	err := bcrypt.CompareHashAndPassword(hashedPwd, pwd)

	// nil means it is a match
	return err
}
