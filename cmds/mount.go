package cmds

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/util/exec"
	"k8s.io/kubernetes/pkg/util/mount"
)

func NewCmdMount() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "mount",
		Short:             `Mount local disk`,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			t := &MountDiskTask{
				Device:     "/dev/sdb",
				Mountpoint: "/mnt/master-pd",
			}
			renderLocal(t)
		},
	}
	return cmd
}

// MountDiskTask is responsible for mounting a device on a mountpoint
// It will wait for the device to show up, safe_format_and_mount it,
// and then mount it.
type MountDiskTask struct {
	Name string

	Device     string `json:"device"`
	Mountpoint string `json:"mountpoint"`
}

func renderLocal(e *MountDiskTask) error {
	dirMode := os.FileMode(0755)

	// Create the mountpoint
	err := os.MkdirAll(e.Mountpoint, dirMode)
	if err != nil {
		return fmt.Errorf("error creating mountpoint %q: %v", e.Mountpoint, err)
	}

	// Wait for the device to show up
	for {
		_, err := os.Stat(e.Device)
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			return fmt.Errorf("error checking for device %q: %v", e.Device, err)
		}
		fmt.Printf("Waiting for device %q to be attached\n", e.Device)
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Found device %q\n", e.Device)

	// Mount the device
	if e.Mountpoint != "" {
		fmt.Printf("Mounting device %q on %q\n", e.Device, e.Mountpoint)

		mounter := &mount.SafeFormatAndMount{Interface: mount.New(""), Runner: exec.New()}

		fstype := "ext4"
		options := []string{}

		err := mounter.FormatAndMount(e.Device, e.Mountpoint, fstype, options)
		if err != nil {
			return fmt.Errorf("error formatting and mounting disk %q on %q: %v", e.Device, e.Mountpoint, err)
		}
	}

	// TODO: Should we add to /etc/fstab?
	// Mount the master PD as early as possible
	// echo "/dev/xvdb /mnt/master-pd ext4 noatime 0 0" >> /etc/fstab

	return nil
}
