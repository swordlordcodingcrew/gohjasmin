package gohjaslib

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
 ** Parts of this source (WatchAuthenticationConfig) are
 ** Copyright © 2014 Steve Francia <spf@spf13.com>.
 **
-----------------------------------------------------------------------------*/

import (
	"encoding/base64"
	"encoding/csv"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// demo.domain.com,demo,$2a$10$v5ZTJXisKF3b83K.XpKoZOYpOoi77qDkq.hLa0b6.8A/sc85aq2F2,127.0.0.1
type record struct {
	Domain     string
	Username   string
	PasswdHash string
	IP         string
}

var auth map[string]record

func InitAuth() {

	loadAuth()
	watchAuthenticationConfig()
}

func loadAuth() {

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

	LogInfo("Authentication file loaded.", logrus.Fields{"file": sFile})

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

func watchAuthenticationConfig() {
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			LogError("Can not init watcher", logrus.Fields{"error": err})
			return
		}
		defer watcher.Close()

		fileAuth := GetStringFromConfig("auth.file")

		// we have to watch the entire directory to pick up renames/atomic saves in a cross-platform way
		file := filepath.Clean(fileAuth)
		dir, _ := filepath.Split(file)

		done := make(chan bool)
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if filepath.Clean(event.Name) == file {

						if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {

							LogInfo("Authentication file changed", logrus.Fields{"file": event.Name, "event": event.Op})
							loadAuth()
						}
					}
				case err := <-watcher.Errors:
					LogDebug("error", logrus.Fields{"error": err})
				}
			}
		}()

		watcher.Add(dir)
		<-done
	}()
}

const AuthUserName = "user"
const AuthIsAuthenticated = "isauthenticated"

// BasicAuth returns a Basic HTTP Authorization middleware.
// (see http://tools.ietf.org/html/rfc2617#section-1.2)
func BasicAuth() gin.HandlerFunc {

	realm := "Basic realm=" + strconv.Quote("Oh Jasmin")

	return func(c *gin.Context) {

		authHeader := c.Request.Header.Get("Authorization")

		if len(authHeader) == 0 {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		auth := strings.SplitN(authHeader, " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			c.AbortWithStatus(500)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// TODO clean username and password before use
		username := pair[0]
		password := pair[1]

		// TODO Check user and Password against the database
		isAuthenticated := ValidateUser(username, password)
		if isAuthenticated {

			// The user credentials was found, set user's id to key AuthUserKey in this context, the userId can be read later using
			// c.MustGet(gin.AuthUserKey)
			c.Set(AuthUserName, username)
			c.Set(AuthIsAuthenticated, true)
			return

		} else {

			c.Set(AuthIsAuthenticated, false)
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "dnserr",
				"error":  "not authorised",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
