# Looper

looper是一个简单的循环调度器

[![License](https://img.shields.io/:license-apache%202-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/ant-libs-go/looper?status.png)](http://godoc.org/github.com/ant-libs-go/looper)
[![Go Report Card](https://goreportcard.com/badge/github.com/ant-libs-go/looper)](https://goreportcard.com/report/github.com/ant-libs-go/looper)

## 特性

* 支持多个任务间隔指定时间运行
* 支持多实例部署情况下，避免并发运行

## 安装

	go get github.com/ant-libs-go/looper

## 快速开始

```golang
l := looper.New()
l.AddFunc("task01", 10 * time.Second, nil, 0, func() { ... })
l.AddFunc("task02", 10 * time.Second, nil, 0, func() { ... })
l.AddFunc("task03", 10 * time.Second, nil, 0, func() { ... })
l.Start()
...
l.Stop()
```
