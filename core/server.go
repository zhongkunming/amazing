package core

import (
	"amazing/global"
	"reflect"
	"sync"
)

var services = make([]Loadable, 0)
var lock sync.Mutex

func Server() {
	// 加载需要启动的模块
	for _, service := range services {
		//v := reflect.ValueOf(service)
		//val, ok := v.Interface().(Loadable)

		serviceType := reflect.TypeOf(service)
		serviceTypeName := serviceType.Name()

		global.Log.Infof("%s match to Loadable", serviceTypeName)
		if service.CanLoad() {
			global.Log.Infof("it's loading %s now", serviceTypeName)
			service.Load()
			continue
		}
	}
	select {}
}

func Register(service Loadable) {
	defer lock.Unlock()
	lock.Lock()
	services = append(services, service)
}
