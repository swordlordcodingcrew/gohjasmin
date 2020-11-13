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
-----------------------------------------------------------------------------*/

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type LegacyParams struct {
	Hostname string `form:"hostname" query:"hostname"`
	MyIP     string `form:"myip" query:"myip"`
	Wildcard string `form:"wildcard" query:"wildcard"`
	MX       string `form:"mx" query:"mx"`
	BackMX   string `form:"backmx" query:"backmx"`
}

type V3Params struct {
	Hostname string `form:"hostname" query:"hostname"`
	MyIP     string `form:"myip" query:"myip"`
}

// Legacy URL
// https://username:password@members.dyndns.org/nic/update
// 		?hostname=yourhostname&myip=ipaddress&wildcard=NOCHG&mx=NOCHG&backmx=NOCHG
// TODO get hostname from param, not from login
func LegacyUpdate(c *gin.Context) {

	var params LegacyParams

	// ?hostname=yourhostname&myip=ipaddress&wildcard=NOCHG&mx=NOCHG&backmx=NOCHG
	err := c.ShouldBind(&params)
	if err != nil {
		LogError("Legacy parsing error.", logrus.Fields{"error": err, "url": c.Request.URL})
		replyError(c, "Bad formed request")
		return
	}

	user := c.GetString(AuthUserName)

	ohJasminUpdate(c, user, params.Hostname, params.MyIP)
}

// v3
// https://{user}:{updater client key aka pwd}@members.dyndns.org/v3/update
// 		?hostname={hostname}&myip={IP Address}
// TODO get hostname from param, not from login
func V3Update(c *gin.Context) {

	var params V3Params

	// check if IP changed.
	// if yes, return good
	// if no, return nochg

	// ?hostname={hostname}&myip={IP Address}
	err := c.ShouldBind(&params)
	if err != nil {
		LogError("V3 parsing error.", logrus.Fields{"error": err, "url": c.Request.URL})
		replyError(c, "Bad formed request")
		return
	}

	user := c.GetString(AuthUserName)

	ohJasminUpdate(c, user, params.Hostname, params.MyIP)
}

func GOhJasminUpdate(c *gin.Context) {

	user := c.GetString(AuthUserName)
	domain := GetDomainFromUser(user)
	ip := c.ClientIP()

	ohJasminUpdate(c, user, domain, ip)
}

func ohJasminUpdate(c *gin.Context, user string, domain string, ip string) {

	// TODO: get user from hostname
	// validate that user is allowed to change said hostname

	if HasIPChanged(user, ip) {

		rowsAffected, err := handleIPUpdate(user, domain, ip)
		if err != nil {
			LogError("DB Update error.", logrus.Fields{"error": err, "ip": ip})
			replyError(c, "There was an error updating your IP")
			return
		}

		// check if IP changed: if yes, return good, if no, return nochg
		if rowsAffected == 0 {

			LogInfo("DB not updated. Record not changed.", logrus.Fields{"ip": ip})
			replyNOCHG(c)

		} else {

			LogInfo("DB updated.", logrus.Fields{"rows_affected": rowsAffected, "ip": ip})

			// make sure to set in memory representation of ip correctly
			ChangeIPInMemory(user, ip)
			replyGOOD(c)
		}

	} else {

		LogDebug("IP was not changed.", logrus.Fields{"ip": ip})

		replyNOCHG(c)
	}
}

func handleIPUpdate(user string, domain string, ip string) (int, error) {

	rowsAffected, err := UpdateDNSRecord(domain, ip)
	if err != nil {

		return 0, err
	}

	return rowsAffected, nil
}

func replyNOCHG(c *gin.Context) {

	LogDebug("Reply.", logrus.Fields{"status": "nochg"})

	c.JSON(200, gin.H{
		"status": "nochg",
	})
}

func replyGOOD(c *gin.Context) {

	LogDebug("Reply.", logrus.Fields{"status": "good"})

	c.JSON(200, gin.H{
		"status": "good",
	})
}

func replyError(c *gin.Context, error string) {

	LogInfo("Error sent to client.", logrus.Fields{"error": error})

	c.JSON(http.StatusInternalServerError, gin.H{
		"status": "dnserr",
		"error":  error,
	})
}
