package core

import (
	"reflect"
	"service-hub/global"
	"sync"
)

var services = make([]Service, 0)
var lock sync.Mutex

func Server() {
	// 加载需要启动的模块
	for _, service := range services {
		v := reflect.ValueOf(service)
		val, ok := v.Interface().(Service)
		if ok {
			global.Log.Infof("符合加载标准，加载 -> %s进程", val.Description())
			val.Load()
		} else {
			global.Log.Fatalf("不符合加载标准，加载失败 -> %s", val)
		}
	}

	select {}
}

func Register(service Service) {
	defer lock.Unlock()
	lock.Lock()
	services = append(services, service)
}
