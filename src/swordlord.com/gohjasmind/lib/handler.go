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
	"github.com/gin-gonic/gin"
	"swordlord.com/gohjasmin/tablemodule"
	"net/http"
	"strings"
	"fmt"
)

// Legacy URL
// https://username:password@members.dyndns.org/nic/update
// 		?hostname=yourhostname&myip=ipaddress&wildcard=NOCHG&mx=NOCHG&backmx=NOCHG
func LegacyUpdate(c *gin.Context) {

	// ignored -> ?hostname=yourhostname&myip=ipaddress&wildcard=NOCHG&mx=NOCHG&backmx=NOCHG
	replyNOCHG(c)
}

// v3
// https://{user}:{updater client key aka pwd}@members.dyndns.org/v3/update
// 		?hostname={hostname}&myip={IP Address}

func V3Update(c *gin.Context) {

	// check if IP changed.
	// if yes, return good
	// if no, return nochg

	// ignored -> &myip={IP Address}
	replyNOCHG(c)
}

// *********************************
// OhJasmin Return Codes (as JSON)
// { status: good }
// { status: nochg }
// { status: dnserr, error: description }
func DDNSUpdate(c *gin.Context) {

	permissions, exists := c.Get(AuthPermissions)
	if !exists {

		replyError(c, "You have no permission to update DNS records")
		return
	}

	remoteAddr := c.Request.RemoteAddr

	// check first...
	ip := strings.Split(remoteAddr, ":")[0]

	// TODO clean Param domain before use
	domain := c.Param("domain")

	neededPermission := fmt.Sprintf("ddns.update.%s", domain)

	hasPermission := tablemodule.HasPermission(neededPermission, permissions.([]string))
	if !hasPermission {

		replyError(c, "You have no permission to update DNS records")
		return
	}

	handleIPUpdate(domain, ip, c)
}

func handleIPUpdate(domain string, ip string, c *gin.Context) {

	rowsAffected, err := tablemodule.UpdateDomain(domain, ip)
	if err != nil {

		replyError(c, "There was an error updating your IP")
		return
	}

	// check if IP changed.
	// if yes, return good
	// if no, return nochg
	if rowsAffected == 0 {

		replyNOCHG(c)

	} else {

		replyGOOD(c)
	}
}

func replyNOCHG(c *gin.Context) {

	c.JSON(200, gin.H{
		"status": "nochg",
	})
}

func replyGOOD(c *gin.Context) {

	c.JSON(200, gin.H{
		"status": "good",
	})
}

func replyError(c *gin.Context, error string) {

	c.JSON(http.StatusInternalServerError, gin.H{
		"status": "dnserr",
		"error": error,
	})
}


/*
func DDNSUpdateTest(c *gin.Context) {

	user := c.MustGet(AuthUserName).(string)

	//isAuthenticated := c.MustGet(lib.AuthIsAuthenticated).(bool)

	permitted := hasPermission("add.user", c.MustGet(AuthPermissions).([]string))

	fmt.Printf("User %s has permission: %v\n", user, permitted)

	if permitted {

		ip := c.Request.RemoteAddr

		// TODO clean Param domain before use
		domain := c.Param("domain")

		c.JSON(http.StatusOK, gin.H{
			"status": "nochg",
			"message": "ddns updated from: " + user + " to ip: " + ip + " for domain: " + domain,
		})

	} else {

		c.JSON(http.StatusForbidden, gin.H{
			"message": "not authorised",
		})
	}
}

func ACMEUpdate(c *gin.Context) {
	user := c.MustGet(AuthUserName).(string)
	//isAuthenticated := c.MustGet(lib.AuthIsAuthenticated).(bool)

	ip := c.Request.RemoteAddr

	c.JSON(200, gin.H{
		"message": "acme pushed from: " + user + " at: " + ip,
	})
}
*/
