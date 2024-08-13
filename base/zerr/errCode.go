package main

const OK uint32 = 200

// 前三位表示业务，后三位表示功能
// 全局错误码
const SERVER_COMMON_ERROR uint32 = 100001
const REQUEST_PARAM_ERROR uint32 = 100002
const TOKEN_EXPIRE_ERROR uint32 = 100003
const TOKEN_GENERATE_ERROR uint32 = 100004
const DB_ERROR uint32 = 100005
