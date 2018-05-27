package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//每次gfs最大读取的字节数
var MaxRead = 1024

//每次gfs最大写的字节数
var MaxWrite = 1024

//GFS最重要的读写流
type GFSReader struct {
	TargetNode []*FileLocation
	TargetFile string
	location   int64
}
type GFSWriter struct {
	TargetNode []*FileLocation
	TargetFile string //在整个fs中的文件名
	location   int64
}

func (gfs *GFSWriter) Write(p []byte) (n int, err error) {
	logger.Info("start to write")
	if gfs.TargetNode == nil || len(gfs.TargetNode) == 0 {
		//TODO 联系master获得所有的datanode
	}
	var toIndex = 0
	var end = false
	var todo = []byte{}
	for {
		if (toIndex+1)*MaxWrite >= len(p) {
			end = true
			todo = p[toIndex:len(p)]
		} else {
			todo = p[toIndex : (toIndex+1)*MaxWrite]
		}
		toIndex++
		targetIndex := gfs.location / BlockSize
		a := BlockSize - (gfs.location/BlockSize)*BlockSize
		if int(a) > len(todo) {
			sendDataTo(gfs.TargetNode[targetIndex], gfs.TargetFile, todo[:a], targetIndex)
			sendDataTo(gfs.TargetNode[targetIndex+1], gfs.TargetFile, todo[a:len(todo)], targetIndex+1)
		} else {
			sendDataTo(gfs.TargetNode[targetIndex], gfs.TargetFile, todo, targetIndex)
		}

		gfs.location = gfs.location + int64(len(todo))

		if end {
			break
		}
	}
	return len(p), nil
}

func sendDataTo(fl *FileLocation, file string, bb []byte, block int64) (int, error) {
	var blockchip = &FileBlockChip{Block: int(block), FileName: file, Data: bb, Replica: fl.Replica}
	buf := blockchip.Encode()
	//TODO写
	_, err := sendDataToSpecialUrl(fl.Main, buf)
	if err != nil {
		var suc = false
		for _, url := range fl.Replica {
			_, err2 := sendDataToSpecialUrl(url, buf)
			if err2 == nil {
				suc = true
				break
			}
		}
		if suc {
			return len(bb), nil
		} else {
			return 0, fmt.Errorf("an error occurs %s", "---")
		}
	} else {
		return len(bb), nil
	}

}

func sendDataToSpecialUrl(url string, buf *bytes.Buffer) (int, error) {
	resp, err := http.Post(url+"/data/in", "application/octet-stream", buf)
	defer resp.Body.Close()
	if err != nil {
		return 0, err
	} else {
		return 1, nil
	}
}

func (gfs *GFSReader) Read(p []byte) (n int, err error) {
	if gfs.TargetNode == nil || len(gfs.TargetNode) == 0 {
		//TODO 联系master获得所有的datanode
	}
	if cap(p) < MaxRead {
		return 0, fmt.Errorf("the cap of byte slice is smaller than %d,use a bigger one", MaxRead)
	}
	targetIndex := int(gfs.location / BlockSize)
	if targetIndex == len(gfs.TargetNode) {
		return 0, nil
	}

	bb, err := getDataFrom(gfs.TargetNode[targetIndex], gfs.TargetFile, gfs.location%BlockSize, int64(MaxRead), targetIndex)
	if err != nil {
		return 0, err
	}
	logger.Info(len(bb))
	copy(p, bb)
	var previousLength = int64(len(bb))
	gfs.location = gfs.location + previousLength
	if len(bb) < MaxRead { //获得的字节少于MaxRead，有两种情况：已到整个文件末尾了； 在两个文件块中
		if targetIndex == len(gfs.TargetNode)-1 {
			return len(bb), nil
		} else {
			bb, err = getDataFrom(gfs.TargetNode[targetIndex+1], gfs.TargetFile, gfs.location%BlockSize, int64(MaxRead)-previousLength, targetIndex+1)
			if err != nil {
				return 0, err
			}
			copy(p[previousLength:], bb)
			gfs.location = gfs.location + int64(len(bb))
			return int(previousLength) + len(bb), nil
		}
	} else {
		return len(bb), nil
	}

}

//从目标datanode中，读取limit个字节
//file: 文件名
//offset: 偏移的字节量
//limit: 读取的最大字节数
func getDataFrom(fl *FileLocation, file string, offset int64, limit int64, block int) ([]byte, error) {
	logger.Infof("retrieve block %d of file %s from location %s with limit %d offset %d", block, file, fl.Main, limit, offset)
	var fbc = &FileBlockChip{FileName: file, Limit: limit, Offset: offset, Block: block}
	buf := fbc.Encode()
	res, err := getDataFromSpecialURL(buf, fl.Main)
	//如果主要的datanode出现错误，那就从replica中寻找
	if err != nil {
		var found bool = false
		for _, url := range fl.Replica {
			res, err = getDataFromSpecialURL(buf, url)
			if err == nil {
				found = true
				break
			}
		}
		if found {
			return res, nil
		} else {
			return nil, errors.New("failed to retrieve data from all relative datanodes")
		}
	} else {
		if len(res) == 0 {
			return nil, io.EOF
		} else {
			return res, nil
		}
	}
}

func getDataFromSpecialURL(data *bytes.Buffer, url string) ([]byte, error) {
	var resp *http.Response
	var err error
	logger.Info("http://" + url + "/data/out")
	resp, err = http.Post("http://"+url+"/data/out", "application/octet-stream", data)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	} else {
		if bb, err2 := ioutil.ReadAll(resp.Body); err2 == nil {
			return bb, nil
		} else {
			return nil, err
		}
	}
}

func GetDataFromSpecialURL(data *bytes.Buffer, url string) ([]byte, error) {
	var resp *http.Response
	var err error
	resp, err = http.Post(url+"/data/out", "application/octet-stream", data)
	defer resp.Body.Close()
	if err != nil {

		return nil, err
	} else {

		if bb, err2 := ioutil.ReadAll(resp.Body); err2 == nil {
			logger.Info(string(bb))
			return bb, nil
		} else {

			return nil, err
		}
	}
}
