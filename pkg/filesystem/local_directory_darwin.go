//go:build darwin
// +build darwin

package filesystem

import (
	"os"

	"github.com/buildbarn/bb-storage/pkg/filesystem/path"

	"golang.org/x/sys/unix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// rawDeviceNumber is the equivalent of POSIX dev_t.
type rawDeviceNumber = int32

// golang.org/x/sys/unix does not define O_SEARCH for darwin. From the macOS
// SDK <sys/fcntl.h>: O_SEARCH == (O_EXEC | O_DIRECTORY) == 0x40000000 | 0x00100000.
const oflagSearch = 0x40100000

func (localDirectory) Mknod(name path.Component, perm os.FileMode, deviceNumber DeviceNumber) error {
	return status.Error(codes.Unimplemented, "Creation of device nodes is not supported on Darwin")
}

func clonefileImpl(oldFD int, oldName string, newFD int, newName string) error {
	return unix.Clonefileat(oldFD, oldName, newFD, newName, unix.CLONE_NOFOLLOW)
}

func (localDirectory) Mount(mountpoint path.Component, source, fstype string) error {
	return status.Error(codes.Unimplemented, "Mount is not supported")
}
