package rollfile

import (
	"os"
	"syscall"
	"time"
)

func GetCTime(finfo os.FileInfo) time.Time {
	if _state, ok := finfo.Sys().(*syscall.Stat_t); ok {
		return time.Unix(_state.Ctim.Sec, _state.Ctim.Nsec)
	}
	return time.Now()
}
