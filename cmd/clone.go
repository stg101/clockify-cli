// Copyright © 2019 Lucas dos Santos Abreu <lucas.s.abreu@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone <time-entry-id>",
	Short: "Copy a time entry and starts it (use \"last\" to copy the last one)",
	Args:  cobra.ExactArgs(1),
	RunE: withClockifyClient(func(cmd *cobra.Command, args []string, c *api.Client) error {
		var err error

		userId, err := getUserId(c)
		if err != nil {
			return err
		}
		workspace := viper.GetString("workspace")
		tec, err := getTimeEntry(
			args[0],
			workspace,
			userId,
			c,
		)

		if err != nil {
			return err
		}

		if tec.TimeInterval.Start, err = convertToTime(whenString); err != nil {
			return err
		}

		tec.TimeInterval.End = nil
		if whenToCloseString != "" {
			var whenToCloseDate time.Time
			if whenToCloseDate, err = convertToTime(whenToCloseString); err != nil {
				return err
			}
			tec.TimeInterval.End = &whenToCloseDate
		}

		if cmd.Flags().Changed("tags") {
			tec.TagIDs, _ = cmd.Flags().GetStringArray("tags")
		}

		if cmd.Flags().Changed("project") {
			tec.ProjectID, _ = cmd.Flags().GetString("project")
		}

		if cmd.Flags().Changed("description") {
			tec.Description, _ = cmd.Flags().GetString("description")
		}

		format, _ := cmd.Flags().GetString("format")
		noClosing, _ := cmd.Flags().GetBool("no-closing")
		asJSON, _ := cmd.Flags().GetBool("json")
		return newEntry(c, tec, viper.GetBool("interactive"), viper.GetBool("allow-project-name"), !noClosing, format, asJSON)
	}),
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	addTimeEntryFlags(cloneCmd)

	cloneCmd.Flags().StringP("project", "", "", "use this project instead")
	cloneCmd.Flags().StringP("description", "", "", "use this description instead")
	cloneCmd.Flags().BoolP("no-closing", "", false, "don't close any active time entry")
	cloneCmd.Flags().StringP("format", "f", "", "golang text/template format to be applied on each time entry")
	cloneCmd.Flags().BoolP("json", "j", false, "print as json")
}
