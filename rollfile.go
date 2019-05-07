package rollfile

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/atomic"
)

type RollCfg struct {
	FileName    string
	TmpPath     string `json:"tmppath"`
	CommitPath  string `json:"commitpath"`
	MaxSize     int
	MaxLine     int
	RollPattern string //yyyymmdd
	Compress    bool
}

type WriteInfo struct {
	FName       string
	TmpPath     string
	CommitPath  string
	MaxSizeByte int64
	MaxLine     int
	RollPattern string
	Commpres    bool
}

type RFWrite struct {
	WriteInfo
	wsize int64
	wline int

	checkTime *atomic.Int32

	currfname string
	file      *os.File
	lock      *sync.Mutex
	fmode     os.FileMode
	tlayout   string
}

func NewRFWrite(cfg RollCfg) (w *RFWrite, err error) {

	w = &RFWrite{
		checkTime: atomic.NewInt32(0),
		lock:      &sync.Mutex{},
		fmode:     0664,

		WriteInfo: WriteInfo{
			FName:       cfg.FileName,
			TmpPath:     cfg.FileName,
			CommitPath:  cfg.CommitPath,
			MaxLine:     cfg.MaxLine,
			MaxSizeByte: int64(cfg.MaxSize * 1024 * 1024),
			Commpres:    cfg.Compress,
		},
		currfname: fmt.Sprintf("%s/%s", cfg.TmpPath, cfg.FileName),
	}

	if len(cfg.RollPattern) > 0 {
		w.tlayout = "20060102150405"
		w.tlayout = w.tlayout[0:len(cfg.RollPattern)]
	}
	go w.TimeTick()

	if CheckRollFile(w.currfname, w.tlayout, time.Now()) {
		err = w.RollFile(true)
	} else {
		_nfile, err := OpenFile(w.currfname, w.fmode)
		if err != nil {
			return w, err
		}
		//不计算文件行数, 只记录文件的大小.
		_fstat, err := os.Stat(w.currfname)
		if err == nil {
			w.wsize = _fstat.Size()
		}

		w.file = _nfile
	}

	return w, nil
}

func (w *RFWrite) WriteString(str string) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.file == nil {
		return -1, fmt.Errorf("file is nil")
	}
	n, err := w.file.WriteString(str)

	w.wsize += int64(n)
	w.wline++

	w.RollFile(false)
	return n, err
}

func (w *RFWrite) Write(b []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.file == nil {
		return -1, fmt.Errorf("file is nil")
	}

	n, err := w.file.Write(b)
	if err != nil {
		return n, err
	}

	w.wsize += int64(n)
	w.wline++

	w.RollFile(false)

	return n, err
}

func (w *RFWrite) TimeTick() {
	for {
		time.Sleep(time.Minute)

		if CheckRollFile(w.currfname, w.tlayout, time.Now()) {
			w.checkTime.Store(1)
		}
	}
}

func (w *RFWrite) RollFile(foreopt bool) error {
	if w.MaxSizeByte > 0 && w.wsize >= w.MaxSizeByte ||
		w.MaxLine > 0 && w.wline >= w.MaxLine ||
		w.checkTime.Load() == 1 ||
		foreopt {

		w.wline = 0
		w.wsize = 0
		w.checkTime.Store(0)

		_t := time.Now()
		backfile := fmt.Sprintf("%s.%s.log", w.currfname[0:len(w.currfname)-len(filepath.Ext(w.currfname))],
			_t.Format("200160102150415.9999"))

		//当rename异常时, 并不会,并不会关闭原来的文件句柄,保证文件能写.
		if err := os.Rename(w.currfname, backfile); err != nil {
			return err
		}

		//当打开新文件异常时, 原来的文件句柄并没有关闭,还是可以写的.
		_nfile, err := OpenFile(w.currfname, w.fmode)
		if err != nil {
			return err
		}

		if w.file != nil {
			w.file.Close()
		}

		w.file = _nfile

		go func() {
			if w.Commpres {
				if err := CompressFile(backfile, backfile+".gz", w.fmode); err == nil {
					RemoveFile(backfile)
				} else {
					fmt.Printf("compress :%v err:%v", backfile, err)
				}
			}
		}()

	}
	return nil
}
