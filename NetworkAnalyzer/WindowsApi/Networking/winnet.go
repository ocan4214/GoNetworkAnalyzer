package windows_networking

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type tcpTableClass int32

const (
	TcpTableBasicListener tcpTableClass = iota
	TcpTableBasicConnections
	TcpTableBasicAll
	TcpTableOwnerPidListener
	TcpTableOwnerPidConnections
	TcpTableOwnerPidAll
	TcpTableOwnerModuleListener
	TcpTableOwnerModuleConnections
	TcpTableOwnerModuleAll
)

func boolToUintPtr(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

// DLL API LINKING
var (
	winiphlpapi             = windows.NewLazySystemDLL("iphlpapi.dll")
	procGetExtendedTCPTable = winiphlpapi.NewProc("GetExtendedTcpTable")
	procGetExtendedUDPTable = winiphlpapi.NewProc("GetExtendedUdpTable")
	procGetIfEntry2         = winiphlpapi.NewProc("GetIfEntry2")
)

func GetExtendedTcpTable(pTcpTable uintptr, pdwSize *uint32, bOrder bool, ulAf uint32, tableClass tcpTableClass, reserved uint32) (errcode error) {
	r1, _, _ := syscall.SyscallN(procGetExtendedTCPTable.Addr(), pTcpTable, uintptr(unsafe.Pointer(pdwSize)), boolToUintPtr(bOrder), uintptr(ulAf), uintptr(tableClass), uintptr(reserved))
	if r1 != 0 {
		errcode = syscall.Errno(r1)
	}
	return
}
