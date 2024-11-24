package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/sys/windows"
)

type SafeProcessIDList struct {
	mutex   sync.Mutex
	idSlice []uint32
}

/* GO Process Start */

func GetProcessID(pname string) (*SafeProcessIDList, error) {

	cmd := exec.Command("tasklist")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var processIDList SafeProcessIDList = SafeProcessIDList{}
	pidArr := make([]uint32, 0)
	// Search for the process name in the output
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pname) {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				pid, _ := (strconv.Atoi(fields[1]))
				pidArr = append(pidArr, uint32(pid))
			}
		}
	}

	processIDList.idSlice = pidArr
	return &processIDList, fmt.Errorf("process %s not found", pname)
}

func IsProcessRunningStatus(pID uint32) (bool, error) {

	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, pID)

	if err != nil {
		return false, err // Process not found or error opening the process
	}
	defer windows.CloseHandle(handle)

	var exitCode uint32
	err = windows.GetExitCodeProcess(handle, &exitCode)
	if err != nil {
		return false, err
	}

	//https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getexitcodeprocess
	// Use numeric value 259 (0x103) for STILL_ACTIVE
	if exitCode == 259 { // 259 is the equivalent of STILL_ACTIVE
		return true, nil // Process is still running
	} else {
		return false, nil // Process has terminated
	}

}

/* GO Process End */

func main() {

	//Get PID first

	processIDList, err := GetProcessID("Spotify.exe")

	if err != nil || processIDList.idSlice == nil || len(processIDList.idSlice) == 0 {
		fmt.Println(err.Error())
	}

	var waitGroup sync.WaitGroup

	for {

		for i, v := range processIDList.idSlice {
			waitGroup.Add(1)
			go func() {
				processIDList.mutex.Lock()

				processIDList.mutex.Unlock()
				waitGroup.Done()
			}()
		}

		waitGroup.Wait()
	}
	fmt.Println(processID)
}

// package main

// import "fmt"

// func fibonacci(c, quit chan int) {
// 	x, y := 0, 1

// 	for {
// 		select {
// 		case c <- x:
// 			x, y = y, x+y
// 		case <-quit:
// 			fmt.Println("quit")
// 			return
// 		}
// 	}
// }

// func main() {
// 	c := make(chan int)
// 	quit := make(chan int)
// 	go func() {
// 		for i := 0; i < 10; i++ {
// 			fmt.Println(<-c)
// 		}
// 		quit <- 0
// 	}()
// 	fibonacci(c, quit)
// }
