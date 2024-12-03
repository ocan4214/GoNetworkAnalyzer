package network_communicator

import (
	CommonUtils "CommonUtilities"
	windows_networking "WindowsApiGo/Networking"
	"syscall"
	"unsafe"
)

const anySize = 1

/* windows type references*/
type MibTCPRowOwnerPid struct {
	DwState      uint32
	DwLocalAddr  uint32
	DwLocalPort  uint32
	DwRemoteAddr uint32
	DwRemotePort uint32
	DwOwningPid  uint32
}

type MibTCPTableOwnerPid struct {
	DwNumEntries uint32
	Table        [anySize]MibTCPRowOwnerPid
}

type (
	PmibTCPTableOwnerPidAll *MibTCPTableOwnerPid
)

/*Communicator Api References*/

type ConnectionStat struct {
	Fd     uint32
	Family uint32
	Type   uint32
	Laddr  Addr
	Raddr  Addr
	Status string
	Uids   []int32
	Pid    int32
}

type Addr struct {
	IP   string
	Port string
}

func (m *MibTCPRowOwnerPid) ConvertToConnectionStats() ConnectionStat {

	ns := ConnectionStat{
		Family: syscall.AF_INET,
		Type:   syscall.SOCK_STREAM,
		Laddr: Addr{
			IP:   CommonUtils.DecodeIP(m.DwLocalAddr),
			Port: CommonUtils.DecodePort(m.DwLocalPort),
		},
		Raddr: Addr{
			IP:   CommonUtils.DecodeIP(m.DwRemoteAddr),
			Port: CommonUtils.DecodePort(m.DwRemotePort),
		},
		Pid:    int32(m.DwOwningPid),
		Status: windows_networking.ConnectionStateMap[int(m.DwState)],
	}

	return ns

}
func GetTableInfo(table interface{}) (index, step, length int) {
	index = int(unsafe.Sizeof(table.(PmibTCPTableOwnerPidAll).DwNumEntries))
	step = int(unsafe.Sizeof(table.(PmibTCPTableOwnerPidAll).Table))
	length = int(table.(PmibTCPTableOwnerPidAll).DwNumEntries)
	return
}

func GetConnectionsByProcessId(processid uint32) []ConnectionStat {

	var (
		p            uintptr
		size         uint32
		buf          []byte
		pmibTCPTable PmibTCPTableOwnerPidAll
	)
	for {
		if len(buf) > 0 {
			pmibTCPTable = (*MibTCPTableOwnerPid)(unsafe.Pointer(&buf[0]))
			p = uintptr(unsafe.Pointer(pmibTCPTable))
		} else {
			p = uintptr(unsafe.Pointer(pmibTCPTable))
		}

		err := windows_networking.GetExtendedTcpTable(
			p,
			&size,
			true,
			syscall.AF_INET, //TCPIPV4
			windows_networking.TcpTableOwnerPidAll,
			0)

		if err == nil {
			break
		}
		buf = make([]byte, size)
	}

	offset := int(unsafe.Sizeof(pmibTCPTable.DwNumEntries))
	step := int(unsafe.Sizeof(pmibTCPTable.Table))
	tableLength := int(pmibTCPTable.DwNumEntries)

	stats := make([]ConnectionStat, 0)

	for i := 0; i < tableLength; i++ {

		item := (*MibTCPRowOwnerPid)(unsafe.Pointer(&buf[offset]))
		stats = append(stats, item.ConvertToConnectionStats())

		offset += step
	}

	return stats
}
