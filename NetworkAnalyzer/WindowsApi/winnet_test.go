package network_communicator_test

import (
	CommonUtils "CommonUtilities"
	network_communicator "WindowsApiGo"
	windows_networking "WindowsApiGo/Networking"
	"fmt"
	"strings"
	"syscall"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

const anySize = 1

// func (m *mibTCPRowOwnerPid) ConvertToConnectionStats() ConnectionStat {

// 	ns := ConnectionStat{
// 		Family: syscall.AF_INET,
// 		Type:   syscall.SOCK_STREAM,
// 		Laddr: Addr{
// 			IP:   parseIPv4HexString(m.DwLocalAddr),
// 			Port: uint32(decodePort(m.DwLocalPort)),
// 		},
// 		Raddr: Addr{
// 			IP:   parseIPv4HexString(m.DwRemoteAddr),
// 			Port: uint32(decodePort(m.DwRemotePort)),
// 		},
// 		Pid:    int32(m.DwOwningPid),
// 		Status: tcpStatuses[mibTCPState(m.DwState)],
// 	}

// 	return ns

// }

// Readable ip
func TestParseIPv4HexString(t *testing.T) {

	var addr uint32 = 0x100007f

	fmt.Printf(fmt.Sprintf("%d.%d.%d.%d", addr&255, addr>>8&255, addr>>16&255, addr>>24&255))
}

func TestDecodePort(t *testing.T) {
	var port uint32 = 0x34c7
	res := ""
	fmt.Println(port)
	if CommonUtils.IsLittleEndian() {
		res = strings.Join(CommonUtils.ConvertU16ToLittleEndianStringSlice(port), "")
	} else {
		res = strings.Join(CommonUtils.ConvertU16ToBigEndianStringSlice(port), "")
	}

	fmt.Println(res)
}

func GetTableInfo(table interface{}) (index, step, length int) {
	index = int(unsafe.Sizeof(table.(network_communicator.PmibTCPTableOwnerPidAll).DwNumEntries))
	step = int(unsafe.Sizeof(table.(network_communicator.PmibTCPTableOwnerPidAll).Table))
	length = int(table.(network_communicator.PmibTCPTableOwnerPidAll).DwNumEntries)
	return
}

func TestGetExtendedTcpTable(t *testing.T) {

	var (
		p            uintptr
		size         uint32
		buf          []byte
		pmibTCPTable network_communicator.PmibTCPTableOwnerPidAll
	)
	for {

		if len(buf) > 0 {
			pmibTCPTable = (*network_communicator.MibTCPTableOwnerPid)(unsafe.Pointer(&buf[0]))
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

	offset := int(unsafe.Sizeof(pmibTCPTable.DwNumEntries))
	step := int(unsafe.Sizeof(pmibTCPTable.Table))
	tableLength := int(pmibTCPTable.DwNumEntries)

	stats := make([]network_communicator.ConnectionStat, 0)

	for i := 0; i < tableLength; i++ {

		item := (*network_communicator.MibTCPRowOwnerPid)(unsafe.Pointer(&buf[offset]))
		stats = append(stats, item.ConvertToConnectionStats())

		offset += step
	}

}
