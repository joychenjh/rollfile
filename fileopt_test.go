package rollfile

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var _tmpdir = "roofile_dir"

func Test_CreateDir(t *testing.T) {
	dir := os.TempDir() + "_tmpdir"

	t.Log("tmpdir:", dir)
	for i := 0; i < 10; i++ {
		if err := CreateDir(dir, os.ModePerm); err != nil {
			t.Errorf("CreateDir dir:%v err:%v", dir, err)
		}
	}
}

func Test_RemoveFile(t *testing.T) {
	dir := os.TempDir() + "_tmpdir"
	CreateDir(dir, os.ModePerm)

	if err := RemoveFile(dir); err == nil {
		t.Errorf("RemoveFile dir must err")
	} else {
		t.Logf("RemoveFile dir err:%v", err)
	}

	_tfname := dir + "/111111.txt"
	_tfile, err := OpenFile(_tfname, os.ModePerm)
	if err != nil {
		t.Error("OpenFile err:", err)
		return
	}
	_tfile.Close()

	if err = RemoveFile(_tfname); err != nil {
		t.Errorf("RemoveFile err:%v", err)
	}

}

func Test_CheckRollFile(t *testing.T) {

	type _rcfg struct {
		mtime   time.Time
		tlayout string
		Roll    bool
		msg     string
	}

	_tfile := os.TempDir() + _tmpdir + "/" + time.Now().Format("20060102150405")

	CreateDir(os.TempDir()+_tmpdir, os.ModePerm)
	_f, err := os.Create(_tfile)
	if err != nil {
		t.Errorf("create file:%v err:%v", _tfile, err)
		return
	}
	defer _f.Close()
	_f.WriteString("HIIHI")
	defer os.Remove(_tfile)
	for _, _v := range []_rcfg{
		{mtime: time.Now().Add(time.Minute), tlayout: "200601021504", Roll: true, msg: "按分钟,1分钟之前的文件."},
		{mtime: time.Now().Add(time.Hour), tlayout: "2006010215", Roll: true, msg: "按小时,1小时之前的文件."},
		{mtime: time.Now().Add(24 * time.Hour), tlayout: "20060102", Roll: true, msg: "按天,1天之前的文件."},
		{mtime: time.Now().Add(24 * 31 * time.Hour), tlayout: "200601", Roll: true, msg: "按月,31天之前的文件."},
		{mtime: time.Now().Add(24 * 366 * time.Hour), tlayout: "2006", Roll: true, msg: "按年,366天之前的文件."},

		{mtime: time.Now(), tlayout: "200601021504", Roll: false, msg: "按分钟,当前时间的."},
		{mtime: time.Now().Add(-1 * time.Minute), tlayout: "2006010215", Roll: false, msg: "按小时.1分钟之前的文件."},
	} {

		if CheckRollFile(_tfile, _v.tlayout, _v.mtime) != _v.Roll {
			t.Errorf("file msg:%v mtime:%v now:%v, tlayout:%v  Roll:%v err", _v.msg, _v.mtime.Format("20060102150405"),
				time.Now().Format("20060102150405"), _v.tlayout, _v.Roll)
		} else {
			t.Logf("file msg:%v mtime:%v now:%v, tlayout:%v Roll:%v ok", _v.msg, _v.mtime.Format("20060102150405"),
				time.Now().Format("20060102150405"), _v.tlayout, _v.Roll)
		}

	}
}

func Test_CompressFile(t *testing.T) {
	_tfile := os.TempDir() + _tmpdir + "/" + time.Now().Format("20060102150405") + ".log"
	_gzfile := _tfile + ".gz"

	CreateDir(os.TempDir()+_tmpdir, os.ModePerm)

	defer os.Remove(_tfile)
	defer os.Remove(_gzfile)

	t.Log("gzfile:", _gzfile)
	buf := []byte(`msg:按小时.1分钟之前的文件. mtime:20190507185031 now:20190507185131, tlayout:2006010215 Roll:false`)

	for i := 0; i < 10; i++ {
		buf = append(buf, buf...)
	}
	if err := ioutil.WriteFile(_tfile, buf, os.ModePerm); err != nil {
		t.Error("WriteFile err:", err)
		return
	}

	if err := CompressFile(_tfile, _gzfile, os.ModePerm); err != nil {
		t.Error("CompressFile err:", err)
		return
	}
}
