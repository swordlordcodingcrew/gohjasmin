package main

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
	"fmt"
	"os"

	"internal/gohjaslib"
	"internal/gohjaslib/command"
)

func main() {

	//
	gohjaslib.InitConfig()

	// only log database actions when env is set to "dev"
	env := gohjaslib.GetStringFromConfig("env")
	bLog := env == "dev"

	gohjaslib.InitDatabase(bLog)
	defer gohjaslib.CloseDB()

	if err := command.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
