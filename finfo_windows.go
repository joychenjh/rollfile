package rollfile

import (
	"os"
	"syscall"
	"time"
)

func GetCTime(finfo os.FileInfo) time.Time {
	if _state, ok := finfo.Sys().(*syscall.Win32FileAttributeData); ok {
		return time.Unix(0, _state.CreationTime.Nanoseconds())
	}
	return time.Now()
}
