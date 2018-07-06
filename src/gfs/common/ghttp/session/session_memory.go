package session

import (
	"container/list"
	"sync"
	"time"
)

type MapSession struct {
	content      map[interface{}]interface{}
	sid          string
	timeAccessed time.Time //访问该sessoin的最后时间
	lock         sync.RWMutex
}

func (s *MapSession) Set(k, v interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.content[k] = v
	return nil
}

func (s *MapSession) Get(k interface{}) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if v, ok := s.content[k]; ok {
		return v
	}
	return nil
}

func (s *MapSession) Delete(k interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.content, k)
	return nil
}

func (s *MapSession) SessionId() string {
	return s.sid
}

func (s *MapSession) Flush() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.content = make(map[interface{}]interface{})
	return nil
}

type MapSessionProvider struct {
	lock        sync.RWMutex
	content     map[string]*list.Element
	list        *list.List
	maxLifeTime int64
}

func (sp *MapSessionProvider) SessionInit(config string) {

}
func (sp *MapSessionProvider) GetSession(sid string) Session {
	sp.lock.Lock()
	if element, ok := sp.content[sid]; ok {
		ms := element.Value.(*MapSession)
		ms.timeAccessed = time.Now()
		sp.list.MoveToFront(element)
		sp.lock.Unlock()
		return ms
	}
	sp.lock.Unlock()
	sp.lock.Lock()
	ms := &MapSession{
		content:      make(map[interface{}]interface{}),
		sid:          sid,
		timeAccessed: time.Now(),
	}
	element := sp.list.PushFront(ms)
	sp.content[sid] = element
	sp.lock.Unlock()
	return ms
}

func (sp *MapSessionProvider) CheckIfSessionExists(sid string) bool {
	sp.lock.Lock()
	defer sp.lock.Unlock()
	_, found := sp.content[sid]
	return found
}
func (sp *MapSessionProvider) SessionDestroy(sid string) error {
	sp.lock.Lock()
	defer sp.lock.Unlock()
	if element, ok := sp.content[sid]; ok {
		delete(sp.content, sid)
		sp.list.Remove(element)
	}
	return nil
}
func (sp *MapSessionProvider) SessionGC() {
	sp.lock.RLock()
	for {
		if element := sp.list.Back(); element != nil {
			ms := element.Value.(*MapSession)
			if ms.timeAccessed.Unix()+sp.maxLifeTime < time.Now().Unix() {
				sp.lock.RUnlock()
				sp.lock.Lock()
				sp.list.Remove(element)
				delete(sp.content, ms.sid)
				sp.lock.Unlock()
				sp.lock.RLock()
			} else {
				break
			}
		} else {
			break
		}

	}
	sp.lock.RUnlock()
}
