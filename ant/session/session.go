/*
@Time : 2018/11/9 15:50 
@Author : sky
@Description:
@File : session
@Software: GoLand
*/
package session

import (
	"sync"
	"time"
)

type ISession interface {
	Set(key,value interface{})
	Get(key interface{})interface{}
	Remove(key interface{})error
	GetCurrentId()string
	Count()int
}


type IStorage interface {
	//实例化一个session
	InitSession(sid string,maxAge int64)(ISession,error)
	GetSession(sid string)interface{}
	SetSession(session ISession)error
	DestroySession(sid string)error
	GCSession()
}


//session 存储在内存
type SessionFromMemory struct {
	sid string
	lock sync.Mutex
	lastAccessedTime time.Time
	maxAge int64
	data map[interface{}]interface{}
}

//实例化

func newSessionFromMemory()*SessionFromMemory{
	return &SessionFromMemory{
		data:make(map[interface{}]interface{}),
		maxAge:60*30,//有效时间30分钟
	}
}

//实现接口

func(this *SessionFromMemory)Set(key,value interface{}){
	this.lock.Lock()
	defer this.lock.Unlock()
	this.data[key]=value
}


func(this *SessionFromMemory)Get(key interface{})interface{}{
	if value,ok:=this.data[key];ok{
		return value
	}
	return nil
}
func(this *SessionFromMemory)Remove(key interface{})error{
	if _,ok:=this.data[key];ok{
		defer this.lock.Unlock()
		this.lock.Lock()
		delete(this.data,key)
	}
	return nil
}
func(this *SessionFromMemory)GetCurrentId()string{
	return this.sid
}
func(this *SessionFromMemory)Count()int{
	return len(this.data)
}