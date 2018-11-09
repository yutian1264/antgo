/*
@Time : 2018/11/9 16:11 
@Author : sky
@Description:
@File : memory_session
@Software: GoLand
*/
package session

import (
	"sync"
	"time"
)

type MemorySession struct {
	lock sync.Mutex
	sessions map[string]ISession
}

func newMemorySession()*MemorySession{

	return &MemorySession{
		sessions:make(map[string]ISession,0),
	}
}

//实现存储方式

func(this *MemorySession)InitSession(sid string,maxAge int64)(ISession,error){
	this.lock.Lock()
	defer this.lock.Unlock()
	//实例化一个sesionFromMemory
	sfm:=newSessionFromMemory()
	sfm.sid=sid
	if maxAge!=0{
		sfm.maxAge=maxAge
	}
	sfm.lastAccessedTime=time.Now()
	this.sessions[sid]=sfm
	return sfm,nil
}

func(this *MemorySession)SetSession(session ISession)error{
	this.lock.Lock()
	defer this.lock.Unlock()
	this.sessions[session.GetCurrentId()] = session
	return nil
}
func(this *MemorySession)GetSession(sid string )interface{}{
	this.lock.Lock()
	defer this.lock.Unlock()
	if value,ok:=this.sessions[sid];ok{
		return value
	}
	return nil
}
func(this *MemorySession)DestroySession(sid string)error{
	this.lock.Lock()
	defer this.lock.Unlock()
	if _,ok:=this.sessions[sid];ok{
		delete(this.sessions,sid)
	}
	return nil
}
func(this *MemorySession)GCSession(){

	sessions:=this.sessions
	if len(sessions)<1{
		return
	}
	for k,v:=range sessions{
		t:=(v.(*SessionFromMemory).lastAccessedTime.Unix())+ (v.(*SessionFromMemory).maxAge)

		if t<time.Now().Unix(){
			//超时了
			this.lock.Lock()
			defer this.lock.Unlock()
			delete(this.sessions, k)
		}
	}
}
