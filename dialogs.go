package main

import (
	"syscall"
	"unsafe"
)

const (
	MB_OK                = 0x00000000
	MB_OKCANCEL          = 0x00000001
	MB_ABORTRETRYIGNORE  = 0x00000002
	MB_YESNOCANCEL       = 0x00000003
	MB_YESNO             = 0x00000004
	MB_RETRYCANCEL       = 0x00000005
	MB_CANCELTRYCONTINUE = 0x00000006
	MB_ICONHAND          = 0x00000010
	MB_ICONQUESTION      = 0x00000020
	MB_ICONEXCLAMATION   = 0x00000030
	MB_ICONASTERISK      = 0x00000040
	MB_USERICON          = 0x00000080
	MB_ICONWARNING       = MB_ICONEXCLAMATION
	MB_ICONERROR         = MB_ICONHAND
	MB_ICONINFORMATION   = MB_ICONASTERISK
	MB_ICONSTOP          = MB_ICONHAND

	MB_DEFBUTTON1 = 0x00000000
	MB_DEFBUTTON2 = 0x00000100
	MB_DEFBUTTON3 = 0x00000200
	MB_DEFBUTTON4 = 0x00000300
)

const (
	IDOK       = 1
	IDCANCEL   = 2
	IDABORT    = 3
	IDRETRY    = 4
	IDIGNORE   = 5
	IDYES      = 6
	IDNO       = 7
	IDCLOSE    = 8
	IDHELP     = 9
	IDTRYAGAIN = 10
	IDCONTINUE = 11
	IDTIMEOUT  = 32000
)

func showMessage(title string, text string, error bool) int {

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms645505(v=vs.85).aspx

	user32 := syscall.MustLoadDLL("user32.dll")
	msgboxf := user32.MustFindProc("MessageBoxW")

	icon := MB_ICONINFORMATION
	if error {
		icon = MB_ICONERROR
	}

	//noinspection GoDeprecation
	ret, _, _ := msgboxf.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(icon))

	return int(ret)

}

func askMessage(title string, text string) (yes bool, res int) {

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms645505(v=vs.85).aspx

	user32 := syscall.MustLoadDLL("user32.dll")
	msgboxf := user32.MustFindProc("MessageBoxW")

	//noinspection GoDeprecation
	ret, _, _ := msgboxf.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(MB_ICONQUESTION|MB_YESNO))

	return ret == IDYES, int(ret)

}
