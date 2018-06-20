// Copyright © 2018 EOS Canada <alex@eoscanada.com>
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
	"fmt"

	"os"

	"encoding/json"

	"github.com/spf13/cobra"
)

var getInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Retrieve blockchain infos, like head block, chain ID, etc..",
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		info, err := api.GetInfo()
		if err != nil {
			fmt.Println("Error: get info:", err)
			os.Exit(1)
		}

		data, err := json.MarshalIndent(info, "", "  ")

		if err != nil {
			fmt.Println("Error: json conversion:", err)
			os.Exit(1)
		}

		fmt.Println(string(data))
	},
}

func init() {
	getCmd.AddCommand(getInfoCmd)
}
