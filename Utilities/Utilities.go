package CommonUtils

import (
	"fmt"
	"time"
)

func GetCurrentDate() string {
	d := time.Now()

	return fmt.Sprintln("%v/%v/%v %v:%v:%v", d.Day(), d.Month(), d.Year(), d.Hour(), d.Minute(), d.Second())
}
