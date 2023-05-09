package server

import (
	"context"
	"fmt"
	"runtime"
	"strings"
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

	query := strings.Split(req.Attributes.Request.Http.Path, "?")

	fmt.Println(query)

	if len(query) > 1 && strings.Contains(strings.ToLower(query[1]), "token") {

		newHeaders := map[string]string{"name": "luckQuery"}
		fmt.Println("query string, 验证通过")
		return OkResponse(&newHeaders, &newHeaders)
	}

	headers := req.Attributes.Request.Http.Headers

	_, ok := headers["token"]

	if !ok {

		msg := "请重新登录"

		fmt.Println("登录失败")

		return DeniedResponse(msg, rpc.UNAUTHENTICATED)

	}

	newHeaders := map[string]string{"name": "luckHeader"}

	fmt.Println("Header, 验证通过")
	return OkResponse(&newHeaders, &headers)
}
