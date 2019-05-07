package rollfile

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

func CreateDir(path string, fmode os.FileMode) (err error) {

	if err = os.MkdirAll(path, fmode); err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}

	return nil
}

func CompressFile(oldfile string, gzfname string, fmode os.FileMode) (err error) {
	defer func() {
		if err != nil {
			fmt.Printf("compress :%v -> %v err:%v", oldfile, gzfname, err)
		}
	}()

	if err = CreateDir(path.Dir(gzfname), fmode); err != nil {
		return err
	}

	gzfile, err := os.OpenFile(gzfname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fmode)
	if err != nil {
		return err
	}
	defer gzfile.Close()

	gw := gzip.NewWriter(gzfile)
	defer gw.Close()

	infile, err := os.Open(oldfile)
	if err != nil {
		return err
	}
	defer infile.Close()

	if _, err := io.Copy(gw, infile); err != nil {
		return err
	}
	return gw.Flush()
}

func RemoveFile(fname string) (err error) {
	finfo, err := os.Stat(fname)
	if err != nil {
		return err
	}
	if finfo.Mode().IsDir() {
		return fmt.Errorf("fname:%v is dir", fname)
	}
	if !finfo.Mode().IsRegular() {
		return fmt.Errorf("fname:%s ModeType:%v", fname, finfo.Mode().String())
	}
	return os.Remove(fname)
}

func OpenFile(fname string, fmode os.FileMode) (f *os.File, err error) {

	f, err = os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, fmode)
	if err != nil {
		return f, err
	}
	return f, nil
}

//检查文件信息,根据最后一次写时间,判断是否需要更换时间.
//如果文件为大小为0,则不会更换.
func CheckRollFile(fname string, tlayout string, ut time.Time) bool {

	if len(tlayout) == 0 {
		return false
	}
	finfo, err := os.Stat(fname)
	if err != nil {
		return false
	}

	//fmt.Println("CheckRollFile", GetCTime(finfo).Format(tlayout), ut.Format(tlayout))
	if GetCTime(finfo).Format(tlayout) != ut.Format(tlayout) {
		return true
	}

	return false
}
