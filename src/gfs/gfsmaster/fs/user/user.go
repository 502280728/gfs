package user

import ()

//用户
type User struct {
	Name  string //用户名
	Pass  string //密码
	Group Group  //用户所属组
}

type Group string
