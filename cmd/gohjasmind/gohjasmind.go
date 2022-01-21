package main

/*-----------------------------------------------------------------------------
 **
 ** - GohJasmin -
 **
 ** Copyright 2017-22 by SwordLord - the coding crew - http://www.swordlord.com
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
	"internal/gohjaslib"
	"net/http"
	"net/http/pprof"
)

func main() {

	//
	gohjaslib.InitConfig()

	// only log database actions when env is set to "dev"
	env := gohjaslib.GetStringFromConfig("env")
	bIsDevMode := env == "debug"

	gohjaslib.InitLog()

	if !bIsDevMode {
		gin.SetMode(gin.ReleaseMode)
	}

	gohjaslib.InitAuth()

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Group using gin.BasicAuth() middleware
	authorized := r.Group("/", gohjaslib.BasicAuth())

	// *********************************
	// Legacy URL
	// https://username:password@members.dyndns.org/nic/update?hostname=yourhostname&myip=ipaddress&wildcard=NOCHG&mx=NOCHG&backmx=NOCHG
	authorized.GET("/nic/update", gohjaslib.LegacyUpdate)
	authorized.GET("/nic/update/", gohjaslib.LegacyUpdate)

	// *********************************
	// v3
	// https://{user}:{updater client key aka pwd}@members.dyndns.org/v3/update?hostname={hostname}&myip={IP Address}
	authorized.GET("/v3/update", gohjaslib.V3Update)

	// *********************************
	// v3
	// https://{user}:{updater client key aka pwd}@members.dyndns.org/ohjasmin
	authorized.GET("/ohjasmin", gohjaslib.GOhJasminUpdate)

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

	host := gohjaslib.GetStringFromConfig("www.host")
	port := gohjaslib.GetStringFromConfig("www.port")

	gohjaslib.LogInfoFmt("gohjasmind running on %v:%v", host, port)

	if bIsDevMode {
		gohjaslib.LogDebugFmt("try: curl http://demouser:pwd@%s:%s/ddns/update/my.dyndns.domain", host, port)
		gohjaslib.LogDebugFmt("try: curl http://demouser:pwd@%s:%s/v3/update", host, port)
		gohjaslib.LogDebugFmt("try: curl http://demouser:pwd@%s:%s/ohjasmin", host, port)
	}

	ret, err := gohjaslib.GetIPFromDomain("demo.domain.com")
	if err != nil {
		gohjaslib.LogDebugFmt("%s:%s", err, ret)
	} else {
		gohjaslib.LogDebugFmt("%s", ret)
	}

	r.Run(host + ":" + port) // listen and serve
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {

	handler := http.HandlerFunc(h)

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
