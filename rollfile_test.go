package rollfile

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Test_Gzfname(t *testing.T) {
	w := RFWrite{currfname: "1111111.log"}
	_t := time.Now()
	backfile := fmt.Sprintf("%s.%s.%04d.log", w.currfname[0:len(w.currfname)-len(filepath.Ext(w.currfname))],
		_t.Format("20060102150415"), _t.Nanosecond()/1000)

	t.Log(backfile)
}

func Test_NewRFWrite(t *testing.T) {
	cfg := RollCfg{
		FileName:    "rollname.log",
		TmpPath:     os.TempDir() + _tmpdir,
		MaxSize:     10,
		MaxLine:     10000,
		RollPattern: "20060102150405",
		Compress:    true,
	}

	rw, err := NewRFWrite(cfg)
	if err != nil {
		t.Error("NewRFWrite err:", err)
		return
	}

	for i := 0; i < cfg.MaxLine*2; i++ {
		_, err = rw.Write([]byte(fmt.Sprintf("line:%010d\n", i)))
		if err != nil {
			t.Error("write err:", err)
		}
	}

	time.Sleep(20 * time.Second)
}

func BenchmarkRFWrite_Write(b *testing.B) {
	cfg := RollCfg{
		FileName:    "rollname.log",
		TmpPath:     os.TempDir() + _tmpdir,
		MaxSize:     -1,
		MaxLine:     -1,
		RollPattern: "20060102",
		Compress:    true,
	}

	rw, err := NewRFWrite(cfg)
	if err != nil {
		b.Error("NewRFWrite err:", err)
		return
	}

	b.Log("tmpdir:", cfg.TmpPath)
	b.ResetTimer()

	line := []byte(`8912345678912345678912345678912345678912345678912345678912345678912345678989123456789123456789123456789123456789123456789123456789123456789123456789\n`)

	for i := 0; i < b.N; i++ {
		_, err = rw.Write(line)
		if err != nil {
			b.Error("write err:", err)
		}
	}
}
func BenchmarkRFWrite_WriteString(b *testing.B) {
	cfg := RollCfg{
		FileName:    "rollname.log",
		TmpPath:     os.TempDir() + _tmpdir,
		MaxSize:     -1,
		MaxLine:     -1,
		RollPattern: "20060102",
		Compress:    true,
	}

	rw, err := NewRFWrite(cfg)
	if err != nil {
		b.Error("NewRFWrite err:", err)
		return
	}

	b.Log("tmpdir:", cfg.TmpPath)
	b.ResetTimer()

	line := "8912345678912345678912345678912345678912345678912345678912345678912345678989123456789123456789123456789123456789123456789123456789123456789123456789\r\n"

	for i := 0; i < b.N; i++ {
		_, err = rw.WriteString(line)
		if err != nil {
			b.Error("write err:", err)
		}
	}

}
