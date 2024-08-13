package main

const defaultMsg = "服务器开小差啦,稍后再来试一试"

var msg map[uint32]string

func init() {
	msg = make(map[uint32]string)
	msg[OK] = "SUCCESS"
	msg[SERVER_COMMON_ERROR] = "服务器开小差啦，请稍后再来试一试～"
	msg[REQUEST_PARAM_ERROR] = "请求参数错误哦"
	msg[TOKEN_EXPIRE_ERROR] = "token过期啦，请重新登陆"
	msg[TOKEN_GENERATE_ERROR] = "token生成失败"
	msg[DB_ERROR] = "数据库开小差啦，请稍后再试一试～"
}

// ParseErrMsg 根据错误码，获取错误信息
func ParseErrMsg(ecode uint32) string {
	if v, ok := msg[ecode]; ok {
		return v
	}
	return defaultMsg
}

// IsCodeErr 是否是自定义错误码
func IsCodeErr(ecode uint32) bool {
	_, ok := msg[ecode]
	return ok
}
