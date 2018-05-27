package session

import (
	"fmt"
)

type MapSession struct {
	content map[interface{}]interface{}
	sid     string
}

type MapSessionProvider struct {
	content map[string]MapSession
}

func (s MapSession) Set(k, v interface{}) error {
	s.content[k] = v
	return nil
}

func (s MapSession) Get(k interface{}) interface{} {
	return s.content[k]
}

func (s MapSession) Delete(k interface{}) error {
	delete(s.content, k)
	return nil
}

func (s MapSession) SessionId() string {
	return s.sid
}

func (s MapSession) Flush() error {
	s.content = make(map[interface{}]interface{})
	return nil
}

func (sp MapSessionProvider) SessionInit(config string) {

}
func (sp MapSessionProvider) GetSession(sid string) Session {
	if !sp.CheckIfSessionExists(sid) {
		sp.content[sid] = MapSession{
			content: make(map[interface{}]interface{}),
			sid:     sid,
		}
	}
	return sp.content[sid]

}

func (sp MapSessionProvider) CheckIfSessionExists(sid string) bool {
	_, found := sp.content[sid]
	return found
}
func (sp MapSessionProvider) SessionDestroy(sid string) error {
	delete(sp.content, sid)
	return nil
}
func (sp MapSessionProvider) SessionGC() {
	fmt.Println("session gc")
}
