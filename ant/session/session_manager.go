/*
@Time : 2018/11/9 17:08 
@Author : sky
@Description:
@File : session_manager
@Software: GoLand
*/
package session

import (
	"sync"
	"net/http"
	"net/url"
	"time"
	"fmt"
	"io"
	"encoding/base64"
	"crypto/rand"
)

type SessionManager struct {
	cookieName string
	storage    IStorage
	maxAge     int64
	lock       sync.Mutex
}

func NewSessionManager() *SessionManager {
	sessionManager := &SessionManager{
		cookieName: "ant-cookie",
		storage:    newMemorySession(),
		maxAge:     60 * 30, //默认30分钟
	}
	go sessionManager.storage.GCSession()
	return sessionManager
}

func (m *SessionManager) GetCookieN() string {
	return m.cookieName
}
func (m *SessionManager) BeginSession(w http.ResponseWriter, r *http.Request) ISession {
	m.lock.Lock()
	defer m.lock.Unlock()
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		sid := m.randomId()
		session, _ := m.storage.InitSession(sid, m.maxAge)
		maxAge := m.maxAge
		if maxAge == 0 {
			maxAge = session.(*SessionFromMemory).maxAge
		}
		cookie := http.Cookie{
			Name:     m.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true, MaxAge: int(maxAge),
			Expires:  time.Now().Add(time.Duration(maxAge)),
		}
		http.SetCookie(w, &cookie) //设置到响应中
		return session
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		//反转义特殊符号
		session := m.storage.(*MemorySession).sessions[sid]
		// 从保存session介质中获取
		fmt.Println("session --------->", session)
		if session == nil {
			fmt.Println("-----------> current session is nil")
			newSession, _ := m.storage.InitSession(sid, m.maxAge) //该方法有自己的锁，多处调用到
			maxAge := m.maxAge
			if maxAge == 0 {
				maxAge = newSession.(*SessionFromMemory).maxAge
			}
			newCookie := http.Cookie{
				Name:     m.cookieName,
				Value:    url.QueryEscape(sid), //转义特殊符号@#￥%+*-等
				Path:     "/",
				HttpOnly: true,
				MaxAge:   int(maxAge),
				Expires:  time.Now().Add(time.Duration(maxAge)),
			}
			http.SetCookie(w, &newCookie) //设置到响应中
			return newSession
		}
		fmt.Println("-----------> current session exists")
		return session
	}
}

//更新超时
func (m *SessionManager) Update(w http.ResponseWriter, r *http.Request) {
	m.lock.Lock()
	defer m.lock.Unlock()
	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return
	}
	t := time.Now()
	sid, _ := url.QueryUnescape(cookie.Value)
	sessions := m.storage.(*MemorySession).sessions
	session := sessions[sid].(*SessionFromMemory)
	session.lastAccessedTime = t
	sessions[sid] = session
	if m.maxAge != 0 {
		cookie.MaxAge = int(m.maxAge)
	} else {
		cookie.MaxAge = int(session.maxAge)
	}
	http.SetCookie(w, cookie)
}

//手动销毁session，同时删除cookie
func (m *SessionManager) Destroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		m.lock.Lock()
		defer m.lock.Unlock()
		sid, _ := url.QueryUnescape(cookie.Value)
		m.storage.DestroySession(sid)
		cookie2 := http.Cookie{
			MaxAge:  0,
			Name:    m.cookieName,
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(time.Duration(0)),
		}
		http.SetCookie(w, &cookie2)
	}
}
func (m *SessionManager) CookieIsExists(r *http.Request) bool {
	_, err := r.Cookie(m.cookieName)
	if err != nil {
		return false
	}
	return true
}

//开启每个会话，同时定时调用该方法
//到达session最大生命时，且超时时。回收它
func (m *SessionManager) GC() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.storage.GCSession()
	//在多长时间后执行匿名函数，这里指在某个时间后执行GC
	time.AfterFunc(time.Duration(m.maxAge*10), func() {
		m.GC()
	})
}

//是否将session放入内存（操作内存）默认是操作内存
func (m *SessionManager) IsFromMemory() {
	m.storage = newMemorySession()
} //是否将session放入数据库（操作数据库）
func (m *SessionManager) IsFromDB() {
	//TODO //关于存数据库暂未实现
}
func (m *SessionManager) SetMaxAge(t int64) {
	m.maxAge = t
}

//如果你自己实现保存session的方式，可以调该函数进行定义
func (m *SessionManager) SetSessionFrom(storage IStorage) {
	m.storage = storage
} //生成一定长度的随机数
func (m *SessionManager) randomId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	//加密
	return base64.URLEncoding.EncodeToString(b)
}
