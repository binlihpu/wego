package core

import (
	"fmt"
	"math/rand"
	"time"
)

//CacheServer cache server interface
type CacheServer interface {
	Get() (string, error)            //获取token
	Set(string, time.Duration) error //设置token
	Lock() error                     //锁
	Unlock()                         //解锁
}

// DefaultMemoryCacheServer cache server using memory
type DefaultMemoryCacheServer struct {
	token     string
	expiresAt time.Time
}

// Get  set token
func (s *DefaultMemoryCacheServer) Get() (string, error) {
	token := s.token

	if s.expiresAt.Before(time.Now()) {
		return token, fmt.Errorf("token expires")
	}
	if token == "" {
		return token, fmt.Errorf("empty token")
	}
	return token, nil
}

// Set  get token
func (s *DefaultMemoryCacheServer) Set(token string, expires time.Duration) error {
	s.token = token
	s.expiresAt = time.Now().Add(expires)
	return nil
}

//Lock just implement the interface
func (s *DefaultMemoryCacheServer) Lock() error {
	return nil
}

//Unlock just implement the interface
func (s *DefaultMemoryCacheServer) Unlock() {
}

//CacheTokenServer token server use cache server
type CacheTokenServer struct {
	TokenFunc           func() (*AccessToken, error)
	CacheServer         CacheServer
	RefreshTokenReqChan chan *AccessToken // chan
}

func abs(x time.Duration) time.Duration {
	if x >= 0 {
		return x
	}
	return -x
}

func retryToken(attempts int, sleep time.Duration, fn func() (string, error)) (string, error) {
	token, err := fn()
	attempts--
	if nil != err && attempts > 0 {
		time.Sleep(sleep)
		return retryToken(attempts, 2*sleep, fn)
	}
	return token, err
}

//NewCacheTokenServer get a new cache token server
func NewCacheTokenServer(cacheServer CacheServer, tokenFunc func() (*AccessToken, error)) *CacheTokenServer {
	srv := &CacheTokenServer{
		TokenFunc:           tokenFunc,
		CacheServer:         cacheServer,
		RefreshTokenReqChan: make(chan *AccessToken),
	}

	go srv.tokenUpdateDaemon(time.Hour * 2 * time.Duration(100+rand.Int63n(200)))
	return srv
}

func (cts *CacheTokenServer) tokenUpdateDaemon(initTickDuration time.Duration) {
	tickDuration := initTickDuration

NEW_TICK_DURATION:
	ticker := time.NewTicker(tickDuration)
	for {
		select {
		case accessToken := <-cts.RefreshTokenReqChan:
			//保存到缓存中
			_ = cts.CacheServer.Set(accessToken.Token, time.Duration(accessToken.ExpiresIn)*time.Second)
			tickDuration = time.Duration(accessToken.ExpiresIn) * time.Second
			ticker.Stop()
			goto NEW_TICK_DURATION
		case <-ticker.C:
			err := cts.CacheServer.Lock()
			if nil == err {
				accessToken, err := cts.TokenFunc()
				if nil == err {
					//保存到缓存中
					_ = cts.CacheServer.Set(accessToken.Token, time.Duration(accessToken.ExpiresIn)*time.Second)
					newTickDuration := time.Duration(accessToken.ExpiresIn) * time.Second
					if abs(tickDuration-newTickDuration) > time.Second*5 {
						tickDuration = newTickDuration
						ticker.Stop()
						goto NEW_TICK_DURATION
					}
				}
			}
			cts.CacheServer.Unlock()
		}
	}
}

//Token get token from CacheTokenServer
func (cts *CacheTokenServer) Token() (string, error) {
	token, err := cts.CacheServer.Get()
	if nil != err || token == "" {
		return cts.RefreshToken()
	}
	return token, err
}

//RefreshToken refresh token from CacheTokenServer
func (cts *CacheTokenServer) RefreshToken() (string, error) {
	err := cts.CacheServer.Lock()
	if nil != err {
		//retry
		// cts.CacheServer.Unlock()
		return retryToken(3, 300*time.Millisecond, cts.CacheServer.Get)
	}
	accessToken, err := cts.TokenFunc()
	//保存到缓存中
	if nil != err {
		return "", err
	}
	cts.RefreshTokenReqChan <- accessToken
	return accessToken.Token, nil
}
