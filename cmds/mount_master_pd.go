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
	"os"
	"strings"
	"time"

	"pharmer.dev/pre-k/lib"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/util/mount"
)

const mountPath = "/mnt/master-pd"

//https://github.com/kubernetes/kubernetes/blob/a6b8c06380526d9631f1f965a61cd1aae2a5832d/cluster/gce/gci/configure-helper.sh#L366
func NewCmdMountMasterPD() *cobra.Command {
	var provider string
	cmd := &cobra.Command{
		Use:               "mount-master-pd",
		Short:             "Mount persistant disk to master",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if provider != "aws" && provider != "gce" {
				Fatal(fmt.Errorf("master pd not supported"))
			}

			safeFormatAndMount := &mount.SafeFormatAndMount{}
			safeFormatAndMount.Interface = mount.New("")
			safeFormatAndMount.Exec = mount.NewOSExec()

			mounts, err := safeFormatAndMount.List()
			if err != nil {
				Fatal(err)
			}
			var existing []*mount.MountPoint
			for _, m := range mounts {
				if m.Path == mountPath {
					existing = append(existing, &m)
				}
			}

			if len(existing) == 0 {
				fmt.Println("Creating mount directory ", mountPath)
				if err = os.MkdirAll(mountPath, 0750); err != nil {
					Fatal(err)
				}
				var device string
				switch provider {
				case "aws":
					device = findAWSMasterPd()
				case "gce":
					device = findGCEMasterPd()
				}
				fmt.Println("Mounting device", device)
				if err = safeFormatAndMount.FormatAndMount(device, mountPath, "", []string{}); err != nil {
					Fatal(err)
				}
			} else {
				glog.Infof("Device already mounted on %q, verifying it is our device", mountPath)

				if len(existing) != 1 {
					glog.Infof("Existing mounts unexpected")

					for i := range mounts {
						m := &mounts[i]
						glog.Infof("%s\t%s", m.Device, m.Path)
					}
					Fatal(fmt.Errorf("found multiple existing mounts at %q", mountPath))
				} else {
					glog.Infof("Found existing mount at %q", mountPath)
				}
			}
			createSymlinks()
		},
	}
	cmd.Flags().StringVar(&provider, "provider", "", "Name of cloud provider")
	return cmd
}

func findAWSMasterPd() string {
	// ref: https://github.com/kubernetes/kubernetes/blob/fe18055adc9ce44a907999f99d239ead47345d5d/cluster/aws/templates/configure-vm-aws.sh#L53
	// ref: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/device_naming.html#available-ec2-device-names
	// Waiting for master pd to be attached
	attempt := 0
	for {
		fmt.Printf("attempt %v to check for /dev/xvdb\n", attempt)
		if _, err := os.Stat("/dev/xvdb"); err == nil {
			fmt.Println("Found /dev/xvdb")
			break
		}
		attempt += 1
		time.Sleep(1 * time.Second)
	}

	// Mount the master PD as early as possible
	//script.AddLine("/etc/fstab", "/dev/xvdb /mnt/master-pd ext4 noatime 0 0")
	return "/dev/xvdb"
}

func findGCEMasterPd() string {
	device := "/dev/disk/by-id/google-master-pd"
	if _, err := os.Stat(device); os.IsNotExist(err) {
		Fatal(fmt.Errorf("device %q not found", device))
	}

	outBytes, _ := lib.ExecCommand("ls", "-l", device).Output()
	out := string(outBytes)
	out = strings.TrimSpace(out)
	relativePath := out[strings.LastIndex(out, " ")+1:]
	return "/dev/disk/by-id/" + relativePath
}

//https://github.com/kubernetes/kubernetes/blob/a6b8c06380526d9631f1f965a61cd1aae2a5832d/cluster/gce/gci/configure-helper.sh#L387
func createSymlinks() {
	//etcd
	fmt.Println("createing etcd symlink")
	if err := os.MkdirAll(mountPath+"/var/lib/", 0700); err != nil {
		Fatal(err)
	}
	if err := lib.Run("ln", "-sf", "/var/lib/etcd", mountPath+"/var/lib/"); err != nil {
		Fatal(err)
	}

	//kubernetes
	fmt.Println("createing kubernetes symlink")

	if err := os.MkdirAll(mountPath+"/etc", 0777); err != nil {
		Fatal(err)
	}
	if err := lib.Run("ln", "-sf", "/etc/kubernetes", mountPath+"/etc/"); err != nil {
		Fatal(err)
	}

	if !lib.UserExists("etcd") {
		if err := lib.Run("useradd", "-s", "/sbin/nologin", "-d", "/var/etcd", "etcd"); err != nil {
			Fatal(err)
		}
	}
	if err := lib.Run("chown", "-R", "etcd", mountPath+"/var/lib/etcd"); err != nil {
		Fatal(err)
	}
	if err := lib.Run("chgrp", "-R", "etcd", mountPath+"/var/lib/etcd"); err != nil {
		Fatal(err)
	}
}
