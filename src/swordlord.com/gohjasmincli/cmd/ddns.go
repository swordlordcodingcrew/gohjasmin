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
	"os"
	"swordlord.com/gohjasmin/tablemodule"
	"swordlord.com/gohjasmin"
	"fmt"
)

var (
	ttl	int
)

// domainCmd represents the domain command
var ddnsCmd = &cobra.Command{
	Use:   "ddns",
	Short: "Add, change and manage your dyndns domains.",
	Long: `Add, change and manage your dyndns domains. Requires a subcommand.`,
	RunE: nil,
}

var ddnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed domains.",
	Long: `List all managed domains.`,
	RunE: ListDDNS,
}

var ddnsExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all managed domains.",
	Long: `Export all managed domains in specified format. If format is left out, format is tinydns.`,
	RunE: ExportDDNS,
}

var ddnsAddCmd = &cobra.Command{
	Use:   "add [domain] [password] [optional ip address]",
	Short: "Add new domain to be managed by this instance of gohjasmin.",
	Long: `Add new domain to be managed by this instance of gohjasmin. If IP argument is left out, 127.0.0.1 will be used.`,
	RunE: AddDDNS,
}

var ddnsUpdateCmd = &cobra.Command{
	Use:   "update [domain] [password] [ip address]",
	Short: "Update the ip address of a domain.",
	Long: `Update the ip address of a domain.`,
	RunE: UpdateDDNS,
}

var ddnsDeleteCmd = &cobra.Command{
	Use:   "delete [domain]",
	Short: "Deletes a domain and removes management by this instance of gohjasmin.",
	Long: `Deletes a domain and removes management by this instance of gohjasmin.`,
	RunE: DeleteDDNS,
}

func ListDDNS(cmd *cobra.Command, args []string) error {

	tablemodule.ListDomains()

	return nil
}

func ExportDDNS(cmd *cobra.Command, args []string) error {

	if ttl == 0 {
		ttl = gohjasmin.GetIntFromConfig("ttl")
	}

	exportPath := cmd.Flags().Lookup("file").Value.String()
	if exportPath == "" {

		tablemodule.ExportDomains(os.Stdout, ttl)

	} else {

		exportFile, err := os.Create(exportPath)
		if err != nil {
			return err
		}

		defer exportFile.Close()

		tablemodule.ExportDomains(exportFile, ttl)
	}

	return nil
}

func AddDDNS(cmd *cobra.Command, args []string) error {

	if len(args) < 2 {
		return fmt.Errorf("command 'add' needs a domain and an optional ip address to be added")
	} else {
		tablemodule.AddDomain(args[0], args[1])
	}

	return nil
}

func UpdateDDNS(cmd *cobra.Command, args []string) error {

	if len(args) < 2 {
		return fmt.Errorf("command 'update' needs a domain and an ip address to be changed to")
	} else {
		tablemodule.UpdateDomain(args[0], args[1])
	}

	return nil
}

func DeleteDDNS(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("command 'delete' needs a domain")
	} else {
		tablemodule.DeleteDomain(args[0])
	}

	return nil
}

func init() {
	RootCmd.AddCommand(ddnsCmd)

	ddnsCmd.AddCommand(ddnsListCmd)

	ddnsCmd.AddCommand(ddnsExportCmd)
	ddnsExportCmd.Flags().StringP("file", "f",  "", "Export ddns records to this file")
	ddnsExportCmd.Flags().IntVar(&ttl, "ttl",  0, "Default ttl for dns records. Uses 300 (5 minutes) if none is given and none or value zero is in the database")

	ddnsCmd.AddCommand(ddnsAddCmd)
	ddnsCmd.AddCommand(ddnsUpdateCmd)
	ddnsCmd.AddCommand(ddnsDeleteCmd)
}
