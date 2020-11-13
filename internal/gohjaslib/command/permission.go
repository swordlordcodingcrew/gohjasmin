package command

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
	"github.com/spf13/cobra"
	"gohjaslib/tablemodule"
)

// domainCmd represents the domain command
var permissionCmd = &cobra.Command{
	Use:   "permission",
	Short: "Add, change and manage permissions.",
	Long:  `Add, change and manage permissions. Requires a subcommand.`,
	RunE:  nil,
}

var permissionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all permissions per users.",
	Long:  `List all permissions per users.`,
	RunE:  ListPermission,
}

var permissionAddCmd = &cobra.Command{
	Use:   "add [user] [permission string]",
	Short: "Add new permission for a given user.",
	Long:  `Add new permission for a given user.`,
	RunE:  AddPermission,
}

var permissionDeleteCmd = &cobra.Command{
	Use:   "delete [permission id]",
	Short: "Deletes a permission string.",
	Long:  `Deletes a permission string.`,
	RunE:  DeletePermission,
}

func ListPermission(cmd *cobra.Command, args []string) error {

	tablemodule.ListPermission()

	return nil
}

func AddPermission(cmd *cobra.Command, args []string) error {

	if len(args) < 2 {
		return fmt.Errorf("command 'add' needs a user and a permission string")
	} else {
		tablemodule.AddPermission(args[0], args[1])
	}

	return nil
}

func DeletePermission(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("command 'delete' needs a permission id")
	} else {
		tablemodule.DeletePermission(args[0])
	}

	return nil
}

func init() {

	RootCmd.AddCommand(permissionCmd)

	permissionCmd.AddCommand(permissionListCmd)
	permissionCmd.AddCommand(permissionAddCmd)
	permissionCmd.AddCommand(permissionDeleteCmd)
}
