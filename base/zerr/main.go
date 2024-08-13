package main

import (
	"github.com/pkg/errors"
	"log"
)

var ErrUserRegistered = NewMsgErr("用户注册失败！！！")

func main() {
	phone := "123456789"
	err := errors.New("hah")

	err = errors.Wrapf(ErrUserRegistered, "用户注册失败 phone: %s", phone)

	if err != nil {
		causeErr := errors.Cause(err)           // err类型
		if e, ok := causeErr.(*CodeError); ok { //自定义错误类型
			log.Printf("【RPC-SRV-ERR-1】 %+v\n", err)
			log.Printf("【RPC-SRV-ERR-1-1】 %+v\n", e)
			//转成grpc err
			//err = status.Error(codes.Code(e.GetErrCode()), e.GetErrMsg())
		} else {
			log.Printf("【RPC-SRV-ERR-2】 %+v", err)
		}

	}
}
