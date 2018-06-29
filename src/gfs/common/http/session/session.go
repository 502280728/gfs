// session
package session

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/satori/uuid"
)

//基本参照 https://github.com/astaxie/beego/blob/master/session/session.go 实现的
//使用其他自定义session是实现方式的方法
//1.自定义一个Session的实现类
//2.自定义一个SessionProvider的实现类
//3.使用RegisterProvider方法注册该SessionProvider
//4.在需要使用的地方用NewProvider方法获取该SessionProvider
//5.使用 SessionStart方法启动session

// session 接口
type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionId() string
	Flush() error // delete all data
}

type SessionProvider interface {
	SessionInit(config string)
	GetSession(sid string) Session //获得session，如果不存在就新建
	CheckIfSessionExists(sid string) bool
	SessionDestroy(sid string) error
	SessionGC()
}

var providers = make(map[string]SessionProvider)
var cookieId = "sessionid"

func RegisterProvider(name string, provider SessionProvider) {
	if provider == nil {
		panic("Session:can not register an nil provider")
	}
	if _, found := providers[name]; found {
		panic(fmt.Sprintf("Session: duplicate name :'%s' when register provider", name))
	}
	providers[name] = provider
}

type Manager struct {
	provider SessionProvider
	config   string
}

func NewProvider(providerName string, config string) *Manager {
	if p, found := providers[providerName]; !found {
		panic(fmt.Sprintf("Session:can not find provider with name '%s'", providerName))
	} else {
		p.SessionInit(config)
		return &Manager{
			p, config,
		}
	}
}

func (m *Manager) getSid(req *http.Request) (string, error) {
	cookie, errs := req.Cookie(cookieId)
	if errs != nil {
		return "", errors.New(fmt.Sprintf("no '%s' found in cookie", cookieId))
	}
	return cookie.Value, nil
}

//从req中获得获得cookie，得到session
//如果req中没有cookie，那么就新建一个cookie，并且把该cookie写入response
func (m *Manager) SessionStart(w http.ResponseWriter, req *http.Request) (Session, error) {
	sid, err := m.getSid(req)
	//如果req中有sessionid，并且该sessionid对应的session已经存在，直接返回该session
	if err == nil && m.provider.CheckIfSessionExists(sid) {
		return m.provider.GetSession(sid), nil
	}
	//如果session不存在，那就新建一个session
	uuid, _ := uuid.NewV4()
	sid = fmt.Sprintf("%s", uuid)
	session := m.provider.GetSession(sid)
	cookie := &http.Cookie{
		Name:  cookieId,
		Value: sid,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	req.AddCookie(cookie)
	return session, nil
}

func (m *Manager) SessionDestroy(w http.ResponseWriter, req *http.Request) error {
	sid, err := m.getSid(req)
	//如果req中有sessionid，并且该sessionid对应的session已经存在，直接返回该session
	if err == nil && m.provider.CheckIfSessionExists(sid) {
		return m.provider.SessionDestroy(sid)
	} else {
		return nil
	}
}
func (m *Manager) SessionGC() {
	m.provider.SessionGC()
	time.AfterFunc(time.Duration(300)*time.Second, func() { m.SessionGC() })
}
