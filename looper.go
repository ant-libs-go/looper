/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-03-05 16:34:07
# File Name: looper.go
# Description:
####################################################################### */

package looper

import (
	"time"

	"github.com/ant-libs-go/redis/lock"
)

type Entry struct {
	name             string
	Spec             time.Duration
	Job              func()
	rdsLock          *lock.Lock
	lockAliveSeconds int64
}

type Looper struct {
	entries []*Entry
	running bool
}

func New() *Looper {
	o := &Looper{running: true}
	return o
}

func (this *Looper) Lock(entry *Entry) bool {
	if entry.rdsLock == nil {
		return true
	}
	if entry.rdsLock.WaitAndLock(entry.lockAliveSeconds) == nil {
		return true
	}
	//key := fmt.Sprintf("looper.%s", entry.name)
	return false
}

func (this *Looper) UnLock(entry *Entry) bool {
	if entry.rdsLock == nil {
		return true
	}
	if entry.rdsLock.Release() == nil {
		return true
	}
	return false
}

func (this *Looper) AddFunc(name string, spec time.Duration, rdsLock *lock.Lock, lockAliveSeconds int64, cmd func()) {
	entry := &Entry{
		name:             name,
		Spec:             spec,
		Job:              cmd,
		rdsLock:          rdsLock,
		lockAliveSeconds: lockAliveSeconds,
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
