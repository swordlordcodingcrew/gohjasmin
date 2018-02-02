package cmd
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
	"github.com/spf13/cobra"
	"swordlord.com/gohjasmin/tablemodule"
	"fmt"
)

// domainCmd represents the domain command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Add, change and manage your users.",
	Long: `Add, change and manage your users. Requires a subcommand.`,
	RunE: nil,
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users.",
	Long: `List all users.`,
	RunE: ListUser,
}

var userAddCmd = &cobra.Command{
	Use:   "add [user] [password]",
	Short: "Add new user to this instance of gohjasmin.",
	Long: `Add new user to this instance of gohjasmin.`,
	RunE: AddUser,
}

var userUpdateCmd = &cobra.Command{
	Use:   "update [user] [password]",
	Short: "Update the password of the user.",
	Long: `Update the password of the user.`,
	RunE: UpdateUser,
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete [user]",
	Short: "Deletes a user.",
	Long: `Deletes a user.`,
	RunE: DeleteUser,
}

func ListUser(cmd *cobra.Command, args []string) error {

	tablemodule.ListUser()

	return nil
}

func AddUser(cmd *cobra.Command, args []string) error {

	if len(args) < 2 {
		return fmt.Errorf("command 'add' needs a user and a password")
	} else {
		tablemodule.AddUser(args[0], args[1])
	}

	return nil
}

func UpdateUser(cmd *cobra.Command, args []string) error {

	if len(args) < 2 {
		return fmt.Errorf("command 'update' needs a user and a new password")
	} else {
		tablemodule.UpdateUser(args[0], args[1])
	}

	return nil
}

func DeleteUser(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("command 'delete' needs a user")
	} else {
		tablemodule.DeleteUser(args[0])
	}

	return nil
}

func init() {
	RootCmd.AddCommand(userCmd)

	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userUpdateCmd)
	userCmd.AddCommand(userDeleteCmd)
}
