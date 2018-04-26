package cache

//datanode本地的cache，存储文件的保存位置
import (
	"gfsnode/data"
)

type Cache map[string]*data.Data
