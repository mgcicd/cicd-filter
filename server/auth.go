package server

import (
	"context"
	"fmt"
	"runtime"
	"time"

	v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/gogo/googleapis/google/rpc"
	"github.com/mgcicd/cicd-core/logs"
)

const MaxGrpcRequestTimeout = 3000

func (a AuthService) Check(ctx context.Context, req *v3.CheckRequest) (*v3.CheckResponse, error) {

	start := time.Now()
	fmt.Println("begign check :", start)

	defer func() {

		end := time.Now()

		dis := end.Sub(start).Nanoseconds() / 1000000

		if dis >= MaxGrpcRequestTimeout {
			logs.DefaultConsoleLog.Warn("Check", fmt.Sprintf("url : %s --- 本次执行时间 ： %v ms", req.Attributes.Request.Http.Path, dis))
		}

		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			logs.DefaultConsoleLog.Error("Check", fmt.Sprintf("authcheck stack: %s err : %v", string(buf[:n]), err))
		}
	}()

	headers := req.Attributes.Request.Http.Headers

	fmt.Println(headers)

	_, ok := headers["token"]

	if !ok {

		msg := "请重新登录"

		fmt.Println("登录失败")

		return DeniedResponse(msg, rpc.UNAUTHENTICATED)

	}

	newHeaders := map[string]string{"name": "luck"}

	fmt.Println("登录成功")
	return OkResponse(&newHeaders, &headers)
}
