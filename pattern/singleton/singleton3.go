package singleton

import "sync"

type singleton3 struct {
}

func (s3 *singleton3) Work() {

}

var s3 *singleton3
var once sync.Once

func newInstance3() *singleton3 {
	return &singleton3{}
}

func GetInstance3() Instance {
	once.Do(func() {
		s3 = newInstance3()
	})
	return s3
}
