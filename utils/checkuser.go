package utils

import (
	"fmt"
	"syscall"
	"unsafe"
)

func CheckHighPriv() bool {
	token, err := syscall.OpenCurrentProcessToken()
	defer token.Close()
	if err != nil {
		fmt.Printf("open current process token failed: %v\n", err)
		return false
	}
	var isElevated uint32
	var outLen uint32
	err = syscall.GetTokenInformation(token, syscall.TokenElevation, (*byte)(unsafe.Pointer(&isElevated)), uint32(unsafe.Sizeof(isElevated)), &outLen)
	if err != nil {
		return false
	}
	return outLen == uint32(unsafe.Sizeof(isElevated)) && isElevated != 0
}
