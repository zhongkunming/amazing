package core

import (
	"mcs/global"
	"reflect"
	"sync"
)

var services = make([]Loadable, 0)
var lock sync.Mutex

func Server() {
	// 加载需要启动的模块
	for _, service := range services {
		v := reflect.ValueOf(service)
		val, ok := v.Interface().(Loadable)

		serviceType := reflect.TypeOf(val)
		serviceTypeName := serviceType.Name()

		if ok {
			global.Log.Infof("%s match to Loadable", serviceTypeName)
			if val.Judge() {
				global.Log.Infof("it's loading %s now", serviceTypeName)
				val.Load()
				continue
			}
		}
		global.Log.Fatalf("%s not match to Loadable", serviceTypeName)
	}
	select {}
}

func Register(service Loadable) {
	defer lock.Unlock()
	lock.Lock()
	services = append(services, service)
}
