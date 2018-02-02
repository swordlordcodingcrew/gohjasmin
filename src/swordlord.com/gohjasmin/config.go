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
	"github.com/spf13/viper"
	"log"
	"io/ioutil"
)

func InitConfig() {

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/gohjasmin")
	// And then register config file name (no extension)
	viper.SetConfigName("gohjasmin.config")
	// Optionally we can set specific config type
	viper.SetConfigType("json")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {

		// TODO: don't just overwrite, check for existence first,
		// then write a standard config file and move on...
		writeStandardConfig()

		if err := viper.ReadInConfig(); err != nil {
			// we tried it once, crash now
			log.Fatalf("Error reading config file: %s", err)
		}
	}
}

// Helper function
func GetStringFromConfig(key string) string {

	return viper.GetString(key)
}

func GetIntFromConfig(key string) int {

	return viper.GetInt(key)
}

//
func writeStandardConfig() (error) {

	err := ioutil.WriteFile("ohjasmin.config.json", defaultConfig, 0700)

	return err
}

//
var defaultConfig = []byte(`
{
  "env": "prod",
  "ttl": "300",
  "www": {
      "host": "127.0.0.1",
      "port": "8081"
  },
  "db": {
    "dialect": "sqlite3",
    "args": "ohjasmin.db"
  }
}
`)

