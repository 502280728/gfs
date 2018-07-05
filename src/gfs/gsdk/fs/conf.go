// conf
package fs

type GFSClientConf struct {
	Master   string //必须是完成的包含protocal://host:port的url
	UserName string
	UserPass string
}

var Configuration *GFSClientConf
