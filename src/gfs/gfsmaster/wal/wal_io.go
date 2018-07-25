// wal_write
package wal

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"sync"
	"time"
)

type WALWritor interface {
	Write(wal *WAL)
	Flush()
}

type WALChanWritor struct {
	buf chan *WAL
}

func (c *WALChanWritor) Write(wal *WAL) {
	c.buf <- wal
}
func (c *WALChanWritor) Flush() {
	for {
		//TODO 重试机制
		tmp := <-c.buf
		if f, err := os.OpenFile(walConf.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend); err == nil {
			f.Write(tmp.ToBytes())
			f.Close()
		}
	}
}

type WALBufWritor struct {
	sync.Mutex
	buf []*WAL
}

func (c *WALBufWritor) Write(wal *WAL) {
	c.Lock()
	defer c.Unlock()
	c.buf = append(c.buf, wal)
}

func (c *WALBufWritor) Flush() {
	c.Lock()
	defer c.Unlock()
	if f, err := os.OpenFile(walConf.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend); err == nil {
		for _, tmp := range c.buf {
			f.Write(tmp.ToBytes())
		}
		f.Close()
		c.buf = c.buf[0:0]
	}
	time.AfterFunc(time.Duration(60)*time.Second, func() {
		c.Flush()
	})
}

type WALReader interface {
	Read(tid string, buf []*WAL) (int, error) //读取从tid(不包含)开始往后的所有WAL
}

type WALFileReader struct {
	offset int64
}

func (w *WALFileReader) Read(tid string, buf []*WAL) (int, error) {
	if cap(buf) == 0 {
		return 0, errors.New("capacity of buf can not be zero")
	}
	if len(buf) == 0 {
		return 0, errors.New("length of buf can not be zero,better equals to capacity")
	}
	file, err := os.Open(walConf.File)
	if err != nil {
		return 0, err
	}
	file.Seek(w.offset, 0)
	scanner := bufio.NewScanner(file)
	var begin = false
	var first = false
	var index = 0
	for {
		if !scanner.Scan() || index > len(buf)-1 {
			break
		}
		line := scanner.Text()
		if strings.HasPrefix(line, tid) {
			begin = true
			first = true
		} else {
			first = false
		}
		if begin && !first {
			//line = strings.TrimSpace(line)
			buf[index] = parseWAL(strings.TrimSpace(line))
			w.offset += int64(len(scanner.Bytes()) + 1)
			index++
		}
	}
	return index, nil
}

func parseWAL(line string) *WAL {
	var res WAL
	res.Parse(line)
	return &res
}
