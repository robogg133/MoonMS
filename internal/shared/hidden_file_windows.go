//go:build windows

package shared

import "syscall"

func SetHidden(path string) {
	filenameW, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		panic(err)
	}

	err = syscall.SetFileAttributes(filenameW, syscall.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		panic(err)
	}

}
