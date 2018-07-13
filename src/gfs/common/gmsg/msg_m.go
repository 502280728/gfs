// 由master往外发出的消息
package gmsg

//master告知sdk每个block的存储位置
type MsgToSDK1 struct {
	Main     string   //主要存储节点
	Replica  []string //备份节点
	FileName string   //文件名称
	Block    int      //数据块id
}
