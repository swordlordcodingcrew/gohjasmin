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
	"github.com/sirupsen/logrus"
	"swordlord.com/gohjasmin"

	//"swordlord.com/gohjasmin/tablemodule"
	"net/http"
	//"fmt"
)

// Legacy URL
// https://username:password@members.dyndns.org/nic/update
// 		?hostname=yourhostname&myip=ipaddress&wildcard=NOCHG&mx=NOCHG&backmx=NOCHG
// TODO get hostname from param, not from login
func LegacyUpdate(c *gin.Context) {

	// ignored -> ?hostname=yourhostname&myip=ipaddress&wildcard=NOCHG&mx=NOCHG&backmx=NOCHG
	//replyNOCHG(c)
	OhJasminUpdate(c)
}

// v3
// https://{user}:{updater client key aka pwd}@members.dyndns.org/v3/update
// 		?hostname={hostname}&myip={IP Address}
// TODO get hostname from param, not from login
func V3Update(c *gin.Context) {

	// check if IP changed.
	// if yes, return good
	// if no, return nochg

	// ignored -> &myip={IP Address}
	OhJasminUpdate(c)
}

func OhJasminUpdate(c *gin.Context) {

	user := c.GetString(AuthUserName)
	ip := c.ClientIP()

	if gohjasmin.HasIPChanged(user, ip) {

		rowsAffected, err := handleIPUpdate(user, ip)
		if err != nil {
			gohjasmin.LogError("DB Update error.", logrus.Fields{"error": err, "ip": ip})
			replyError(c, "There was an error updating your IP")
			return
		}

		// check if IP changed: if yes, return good, if no, return nochg
		if rowsAffected == 0 {

			gohjasmin.LogInfo("DB not updated. Record not changed.", logrus.Fields{"ip": ip})
			replyNOCHG(c)

		} else {

			gohjasmin.LogInfo("DB updated.", logrus.Fields{"rows_affected": rowsAffected, "ip": ip})

			// make sure to set in memory representation of ip correctly
			gohjasmin.ChangeIPInMemory(user, ip)
			replyGOOD(c)
		}

	} else {

		gohjasmin.LogDebug("IP was not changed.", logrus.Fields{"ip": ip})

		replyNOCHG(c)
	}
}

func handleIPUpdate(user string, ip string) (int, error) {

	domain := gohjasmin.GetDomainFromUser(user)

	rowsAffected, err := gohjasmin.UpdateDNSRecord(domain, ip)
	if err != nil {

		return 0, err
	}

	return rowsAffected, nil
}

func replyNOCHG(c *gin.Context) {

	gohjasmin.LogDebug("Reply.", logrus.Fields{"status": "nochg"})

	c.JSON(200, gin.H{
		"status": "nochg",
	})
}

func replyGOOD(c *gin.Context) {

	gohjasmin.LogDebug("Reply.", logrus.Fields{"status": "good"})

	c.JSON(200, gin.H{
		"status": "good",
	})
}

func replyError(c *gin.Context, error string) {

	gohjasmin.LogInfo("Error sent to client.", logrus.Fields{"error": error})

	c.JSON(http.StatusInternalServerError, gin.H{
		"status": "dnserr",
		"error":  error,
	})
}
