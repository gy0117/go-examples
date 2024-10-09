package singleton

import "sync"

type singleton2 struct {
}

var s2 *singleton2
var mux sync.Mutex

func newInstance2() *singleton2 {
	return &singleton2{}
}

// 存在并发问题
func GetInstance2() Instance {
	//mux.Lock()
	//defer mux.Unlock()
	//if s2 == nil {
	//	s2 = newInstance2()
	//}
	//return s2
	if s2 != nil {
		return s2
	}
	mux.Lock()
	defer mux.Unlock()
	s2 = newInstance2()
	return s2
}

func (s2 *singleton2) Work() {

}
