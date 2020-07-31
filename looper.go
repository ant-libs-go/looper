/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-03-05 16:34:07
# File Name: looper.go
# Description:
####################################################################### */

package looper

import (
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v7"
)

type Entry struct {
	name string
	Spec time.Duration
	Job  func()
}

type Looper struct {
	entries  []*Entry
	running  bool
	redisCli *redis.Client
}

func New(rdsOpt *redis.Options) *Looper {
	o := &Looper{running: true}
	if rdsOpt != nil {
		o.redisCli = redis.NewClient(rdsOpt)
	}

	return o
}

// 简单粗暴，锁定一小时。确保哪怕死锁也就一个小时的事儿
func (this *Looper) Lock(entry *Entry) bool {
	if this.redisCli == nil {
		return true
	}
	key := fmt.Sprintf("looper.%s", entry.name)
	res, _ := this.redisCli.Do("SET", key, "1", "EX", 60*60, "NX").String()
	if res == "OK" {
		return true
	}
	return false
}

func (this *Looper) UnLock(entry *Entry) bool {
	if this.redisCli == nil {
		return true
	}
	key := fmt.Sprintf("looper.%s", entry.name)
	res, _ := this.redisCli.Do("DEL", key).Int64()
	if res == 1 {
		return true
	}
	return false
}

func (this *Looper) AddFunc(name string, spec time.Duration, cmd func()) {
	entry := &Entry{
		name: name,
		Spec: spec,
		Job:  cmd,
	}
	this.entries = append(this.entries, entry)
}

func (this *Looper) Start() {
	for _, entry := range this.entries {
		go func(entry *Entry) {
			for this.running {
				if this.Lock(entry) == true {
					entry.Job()
					this.UnLock(entry)
				}
				time.Sleep(entry.Spec)
			}
		}(entry)
	}
}

func (this *Looper) Stop() {
	this.running = false
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
