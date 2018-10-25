package lib

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
	//"crypto/subtle"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"swordlord.com/gohjasmin"
	//"swordlord.com/gohjasmin/tablemodule"
)

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
		isAuthenticated := gohjasmin.ValidateUser(username, password)
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
