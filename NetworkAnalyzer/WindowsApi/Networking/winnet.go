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

var ConnectionStateMap map[int]string = map[int]string{
	1:  "CLOSED",
	2:  "LISTEN",
	3:  "SYN_SENT",
	4:  "SYN_RECEIVED",
	5:  "ESTABLISHED",
	6:  "FIN_WAIT_1",
	7:  "FIN_WAIT_2",
	8:  "CLOSE_WAIT",
	9:  "CLOSING",
	10: "LAST_ACK",
	11: "TIME_WAIT",
	12: "DELETE",
}

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
