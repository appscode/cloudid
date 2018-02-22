package cmds

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pharmer/pre-k/lib"
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
			provider = lib.DetectCloudProvider()
			var path string
			var mount bool
			mountDirectories()
			switch provider {
			case "aws":
				path, mount = findAWSMasterPd()
				break
			case "gce":
				path, mount = findGCEMasterPd()
				break
			default:
				Fatal(fmt.Errorf("pd not supported"))
			}
			if mount {
				//https://github.com/kubernetes/kops/blob/0ef2fde69d076e9717839a587a57aec8562327e1/protokube/pkg/protokube/volume_mounter.go#L93
				fmt.Println("Doing safe-format-and-mount path= ", path)
				safeFormatAndMount(path)
				mountDirectories()
			}
		},
	}
	//	cmd.Flags().StringVar(&diskId, "disk-id", "", "Persistant disk id")
	return cmd
}

func safeFormatAndMount(devicePath string) {
	safeFormatAndMount := &mount.SafeFormatAndMount{}
	safeFormatAndMount.Interface = mount.New("")
	safeFormatAndMount.Exec = mount.NewOsExec()

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
		fmt.Println("Mounting device")
		if err = safeFormatAndMount.FormatAndMount(devicePath, mountPath, "", []string{}); err != nil {
			Fatal(err)
		}
	}
}

func findAWSMasterPd() (string, bool) {
	// ref: https://github.com/kubernetes/kubernetes/blob/fe18055adc9ce44a907999f99d239ead47345d5d/cluster/aws/templates/configure-vm-aws.sh#L53
	// $ grep "/mnt/disks/master-pd" /proc/mounts
	file, err := os.Open("/proc/mounts")
	if err != nil {
		Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), mountPath) {
			Fatal(fmt.Errorf("Master PD already mounted; won't remount"))
		}
	}
	if err := scanner.Err(); err != nil {
		Fatal(err)
	}

	// Waiting for master pd to be attached
	attempt := 0
	for true {
		fmt.Println("Attempt %v to check for /dev/xvdb", attempt)
		if _, err := os.Stat("/dev/xvdb"); err == nil {
			fmt.Println("Found /dev/xvdb")
			break
		}
		attempt += 1
		time.Sleep(1 * time.Second)
	}

	// Mount the master PD as early as possible
	//script.AddLine("/etc/fstab", "/dev/xvdb /mnt/master-pd ext4 noatime 0 0")
	return "/dev/xvdb", true
}

func findGCEMasterPd() (string, bool) {
	devicepath := "/dev/disk/by-id/google-master-pd"
	if _, err := os.Stat(devicepath); os.IsNotExist(err) {
		fmt.Println(devicepath + " does not exist")
		// path does not exist
		return "", false
	}

	outBytes, _ := lib.ExecCommand("ls", "-l", devicepath).Output()
	out := string(outBytes)
	out = strings.TrimSpace(out)
	relativePath := out[strings.LastIndex(out, " ")+1:]
	return "/dev/disk/by-id/" + relativePath, true
}

//https://github.com/kubernetes/kubernetes/blob/a6b8c06380526d9631f1f965a61cd1aae2a5832d/cluster/gce/gci/configure-helper.sh#L387
func mountDirectories() {
	//etcd
	fmt.Println("createing etcd symlink")
	os.MkdirAll(mountPath+"/var/lib/", 0700)
	lib.Run("ln", "-sf", "/var/lib/etcd", mountPath+"/var/lib/")

	//kubernetes
	fmt.Println("createing kubernetes symlink")

	os.MkdirAll(mountPath+"/etc", 0777)
	lib.Run("ln", "-sf", "/etc/kubernetes", mountPath+"/etc/")

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
