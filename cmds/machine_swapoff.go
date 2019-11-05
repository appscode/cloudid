/*
Copyright The Pharmer Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func NewCmdSwapoff() *cobra.Command {
	var (
		filename = "/etc/fstab"
		re       = regexp.MustCompile(`^.*\s+swap\s+.*$`)
	)
	cmd := &cobra.Command{
		Use:               "swapoff",
		Short:             `Permanently disabled swap disks on a machine`,
		Long:              `Runs "swapoff -a" and comments out all swap devices on a machine found in /etc/fstab file.`,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			cout, err := exec.Command("swapoff", "-a").CombinedOutput()
			if err != nil {
				fmt.Println(string(cout))
				Fatal(err)
			}
			fmt.Println(string(cout))

			data, err := ioutil.ReadFile(filename)
			if err != nil {
				Fatal(err)
			}
			lines := strings.Split(string(data), "\n")
			out := make([]string, 0, len(lines))
			for _, line := range lines {
				txt := strings.TrimSpace(line)

				if strings.HasPrefix(txt, "#") {
					out = append(out, line) // leave existing comments unchanged
				} else if re.MatchString(txt) {
					out = append(out, "# "+line) // comment out swap devices
				} else {
					out = append(out, line) // leave regular disks unchanged
				}
			}
			err = ioutil.WriteFile(filename, []byte(strings.Join(out, "\n")), 0644)
			if err != nil {
				Fatal(err)
			}
			os.Exit(0)
		},
	}
	cmd.Flags().StringVar(&filename, "filename", filename, "Path to volume list file.")
	return cmd
}
