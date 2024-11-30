package windows_networking_test

import (
	windows_networking "WindowsApiGo/Networking"
	"syscall"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

const anySize = 1

type mibTCPRowOwnerPid struct {
	DwState      uint32
	DwLocalAddr  uint32
	DwLocalPort  uint32
	DwRemoteAddr uint32
	DwRemotePort uint32
	DwOwningPid  uint32
}

type mibTCPTableOwnerPid struct {
	DwNumEntries uint32
	Table        [anySize]mibTCPRowOwnerPid
}

type (
	pmibTCPTableOwnerPidAll *mibTCPTableOwnerPid
)

func TestGetExtendedTcpTable(t *testing.T) {

	var (
		p            uintptr
		size         uint32
		buf          []byte
		pmibTCPTable pmibTCPTableOwnerPidAll
	)
	for {

		if len(buf) > 0 {
			pmibTCPTable = (*mibTCPTableOwnerPid)(unsafe.Pointer(&buf[0]))
			p = uintptr(unsafe.Pointer(pmibTCPTable))
		} else {
			p = uintptr(unsafe.Pointer(pmibTCPTable))
		}

		err := windows_networking.GetExtendedTcpTable(
			p,
			&size,
			true,
			syscall.AF_INET,
			windows_networking.TcpTableOwnerPidAll,
			0)

		if err == nil {
			break
		}
		if err != windows.ERROR_INSUFFICIENT_BUFFER {
			t.Fatalf(`Error occured %v use this size %v`, err.Error(), size)
		}
		buf = make([]byte, size)
	}

}
