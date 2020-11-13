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
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
)

func InitConfig() {

	// we look in these dirs for the config file
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.gohjasmin")
	viper.AddConfigPath("/etc/gohjasmin")

	// And then register config file name (no extension)
	viper.SetConfigName("gohjasmin.config")
	// Optionally we can set specific config type
	viper.SetConfigType("json")

	// viper allows watching of config files for changes (and potential reloads)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		LogDebug("Config file changed. Reloading config.", logrus.Fields{"file": e.Name})
	})

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {

		// TODO: don't just overwrite, check for existence first,
		// then write a standard config file and move on...
		writeStandardConfig()

		// crash now to have the user check the config file we just wrote out
		LogFatal("Error reading config file. New file dumped.", logrus.Fields{"error": err, "location": "gohjasmin.config.json"})
	}
}

// Helper function
func GetStringFromConfig(key string) string {

	return viper.GetString(key)
}

func GetIntFromConfig(key string) int {

	return viper.GetInt(key)
}

func GetLogLevel() string {

	loglevel := viper.GetString("log.level")
	if loglevel == "" {

		return "warn"
	} else {

		return loglevel
	}
}

func GetSubConfig(key string) *viper.Viper {

	return viper.Sub(key)
}

//
func writeStandardConfig() error {

	err := ioutil.WriteFile("gohjasmin.config.json", defaultConfig, 0700)

	return err
}

//
var defaultConfig = []byte(`
{
  "env": "dev",
  "log.level": "debug",
  "www": {
    "host": "127.0.0.1",
    "port": "8081"
  },
  "db": {
    "file": "gohjasmin.sqlite3",
    "sql.update": "UPDATE records SET content = ?1, change_date = CURRENT_TIMESTAMP WHERE name = ?2 AND type = 'A' AND content != ?1"
  },
  "auth": {
	"file": "gohjasmin.auth"
  },
  "dns": {
    "ttl": "300"
  }
}
`)
