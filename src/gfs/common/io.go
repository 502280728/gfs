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
var MaxRead = 1000

//每次gfs最大写的字节数
var MaxWrite = 1000

//GFS最重要的读写流
type GFSReader struct {
	TargetNode []*FileLocation
	TargetFile string
	location   int64
	MaxSize    int64 //该reader最多能够读取的字节数，如果是读文件的话，就是这个文件的大小
}
type GFSWriter struct {
	TargetNode []*FileLocation
	TargetFile string //在整个fs中的文件名
	location   int64  //写了多少字节,从0开始计数。如果写入的字节超过一个Block的大小，那么转写下一个TargetNode。如果写入的字节超过了一次最大写入的字节，那么等待下一轮的写入
	MaxSize    int64  //该writer最多能够写的字节数，如果是写文件的话，就是这个文件的大小
}

func (gfs *GFSWriter) Write(p []byte) (n int, err error) {
	logger.Info("start to write")
	if gfs.TargetNode == nil || len(gfs.TargetNode) == 0 {
		//TODO 联系master获得所有的datanode
	}
	var toIndex = 0     //思路是将切片p分割中大小为MaxWrite的子切片，一次写入MaxWrite大小的数据。toIndex表示第几个子切片
	var end = false     //是否结束循环
	var todo = []byte{} //表示写入targetnode的数据，该切片是p的子切片
	//一次循环写入切片p中最大为MaxWrite的子切片
	for {
		if (toIndex+1)*MaxWrite >= len(p) {
			end = true
			todo = p[toIndex*MaxWrite : len(p)]
		} else {
			todo = p[toIndex*MaxWrite : (toIndex+1)*MaxWrite]
		}
		toIndex++

		targetIndex := gfs.location / BlockSize //表示切片todo里的数据需要写入那个targetnode

		a := BlockSize - gfs.location%BlockSize //计算该block还能接收的字节数
		if int(a) < len(todo) {                 //如果todo的长度大于还能接受的长度，那么要把todo分成两段传输，当然这么做的前提是block的大小不能小于maxwrite
			sendDataTo(gfs.TargetNode[targetIndex], gfs.TargetFile, todo[:a], targetIndex)
			sendDataTo(gfs.TargetNode[targetIndex+1], gfs.TargetFile, todo[a:len(todo)], targetIndex+1)
		} else {
			sendDataTo(gfs.TargetNode[targetIndex], gfs.TargetFile, todo, targetIndex)
		}

		gfs.location = gfs.location + int64(len(todo))

		if end || gfs.location >= gfs.MaxSize { //已经写完切片p或者写完整个文件，就退出
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
	var pIndex int = 0 //指向切片p中开始写入的位置
	for {
		if pIndex >= cap(p) || gfs.location >= gfs.MaxSize { //只要是把切片p写满或者读到了文件末尾，就退出
			break
		}
		targetIndex := int(gfs.location / BlockSize) //从哪个targetnode读取

		//一次可以读取的数据要从切片p剩余的空间、一个block剩余的数据、以及MaxRead中选择最小的那一个;
		//但是就算如此，最终返回的数据量也可能小于该值。因为可能block中没有那么的数据
		remaina := int(BlockSize - gfs.location%BlockSize)
		remainb := int(cap(p) - pIndex)
		remainc := int(MaxRead)

		var remain int
		if remaina <= remainb {
			remain = remaina
		} else {
			remain = remainb
		}
		if remainc < remain {
			remain = remainc
		}

		if bb, err := getDataFrom(gfs.TargetNode[targetIndex], gfs.TargetFile, gfs.location%BlockSize, int64(remain), targetIndex); err == nil {
			copy(p[pIndex:], bb)
			pIndex = pIndex + len(bb)
			gfs.location = gfs.location + int64(len(bb))
			//			if targetIndex == len(gfs.TargetNode)-1 && len(bb) < remain {
			//				break
			//			}
		} /*else {
			break
		}*/

	}
	if pIndex > 0 {
		return pIndex, nil
	} else {
		return 0, io.EOF
	}
}

//从目标datanode中，读取limit个字节
//file: 文件名
//offset: 偏移的字节量
//limit: 读取的最大字节数
func getDataFrom(fl *FileLocation, file string, offset int64, limit int64, block int) ([]byte, error) {
	//logger.Infof("retrieve block %d of file %s from location %s with limit %d offset %d", block, file, fl.Main, limit, offset)
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
	//logger.Info(url + "/data/out")
	resp, err = http.Post(url+"/data/out", "application/octet-stream", data)
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
