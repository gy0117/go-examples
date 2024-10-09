package singleton

type Instance interface {
	Work()
}

type singleton struct {
}

var s *singleton

func init() {
	s = newInstance()
}

func newInstance() *singleton {
	return &singleton{}
}

// 这个单例是不推荐的，s是不对外暴露的，就算外界调了GetInstance()方法，
// 返回了不可导出类型
//func GetInstance() *singleton {
//	return s
//}

func GetInstance() Instance {
	return s
}

func (s *singleton) Work() {
}
