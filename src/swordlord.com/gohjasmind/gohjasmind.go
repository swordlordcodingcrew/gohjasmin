package main

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
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/http/pprof"
	"swordlord.com/gohjasmin"
	"swordlord.com/gohjasmind/lib"
)

func main() {

	//
	gohjasmin.InitConfig()

	// only log database actions when env is set to "dev"
	env := gohjasmin.GetStringFromConfig("env")
	bIsDevMode := env == "dev"

	//gohjasmin.InitDatabase(bIsDevMode)
	//defer gohjasmin.CloseDB()

	gohjasmin.InitLog()

	if bIsDevMode {

		gin.SetMode(gin.ReleaseMode)
	}

	gohjasmin.LoadAuth()

	r := gin.Default()

	// Group using gin.BasicAuth() middleware
	authorized := r.Group("/", lib.BasicAuth())

	// *********************************
	// Legacy URL
	// https://username:password@members.dyndns.org/nic/update?hostname=yourhostname&myip=ipaddress&wildcard=NOCHG&mx=NOCHG&backmx=NOCHG
	authorized.GET("/nic/update", lib.LegacyUpdate)
	authorized.GET("/nic/update/", lib.LegacyUpdate)

	// *********************************
	// v3
	// https://{user}:{updater client key aka pwd}@members.dyndns.org/v3/update?hostname={hostname}&myip={IP Address}
	authorized.GET("/v3/update", lib.V3Update)

	// *********************************
	// v3
	// https://{user}:{updater client key aka pwd}@members.dyndns.org/ohjasmin
	authorized.GET("/ohjasmin", lib.OhJasminUpdate)

	// *********************************
	// Legacy and v3 Return codes
	// https://help.dyn.com/remote-access-api/return-codes/
	// basically good and nochg, as well as dnserr for errors

	// Debugging in Dev mode only
	if bIsDevMode {

		r.GET("/debug/pprof/block", pprofHandler(pprof.Index))
		r.GET("/debug/pprof/heap", pprofHandler(pprof.Index))
		r.GET("/debug/pprof/profile", pprofHandler(pprof.Profile))
		r.POST("/debug/pprof/symbol", pprofHandler(pprof.Symbol))
		r.GET("/debug/pprof/symbol", pprofHandler(pprof.Symbol))
		r.GET("/debug/pprof/trace", pprofHandler(pprof.Trace))
	}

	host := gohjasmin.GetStringFromConfig("www.host")
	port := gohjasmin.GetStringFromConfig("www.port")

	gohjasmin.LogInfoFmt("gohjasmind running on %v:%v\n", host, port)

	if bIsDevMode {
		gohjasmin.LogDebugFmt("try: curl http://demouser:pwd@%s:%s/ddns/update/my.dyndns.domain\n", host, port)
		gohjasmin.LogDebugFmt("try: curl http://demouser:pwd@%s:%s/v3/update\n", host, port)
		gohjasmin.LogDebugFmt("try: curl http://demouser:pwd@%s:%s/ohjasmin\n", host, port)
	}

	r.Run(host + ":" + port) // listen and serve
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {

	handler := http.HandlerFunc(h)

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
