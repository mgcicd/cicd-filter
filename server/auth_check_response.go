package server

import (
	"errors"
	"strings"

	_ "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	coreV3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authV3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typeV3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/mgcicd/cicd-core/util"
	"google.golang.org/genproto/googleapis/rpc/status"
)

type apiResponse struct {
	Error   int
	Message string
	Data    interface{}
}

type ErrorCode int

const (
	Normal      ErrorCode = iota
	NoLogin     ErrorCode = -2
	SystemError ErrorCode = 999
)

func DeniedResponse(msg string, code rpc.Code) (*authV3.CheckResponse, error) {
	if msg == "" {
		return nil, errors.New("msg is nil")
	}

	errCode := SystemError

	if code == rpc.UNAUTHENTICATED {
		errCode = NoLogin
	}

	return &authV3.CheckResponse{
		Status: &status.Status{
			Code: int32(code),
		},
		HttpResponse: &authV3.CheckResponse_DeniedResponse{
			DeniedResponse: &authV3.DeniedHttpResponse{
				Status: &typeV3.HttpStatus{
					Code: typeV3.StatusCode_OK,
				},
				Headers: []*coreV3.HeaderValueOption{
					{
						Header: &coreV3.HeaderValue{
							Key:   "content-type",
							Value: "application/json; charset=utf-8",
						},
					},
				},
				Body: util.StructToJson(&apiResponse{
					Error:   int(errCode),
					Message: msg,
					Data:    nil,
				}),
			},
		},
	}, nil
}

func OkResponse(addheader *map[string]string, headers *map[string]string) (*authV3.CheckResponse, error) {

	newHeader := make(map[string]string, 0)

	for key, value := range *headers {

		add := false

		for k, v := range *addheader {
			if strings.ToLower(key) == strings.ToLower(k) {
				newHeader[key] = v
				add = true
				break
			}
		}

		if add {
			continue
		}

		newHeader[key] = value
	}

	hos := make([]*coreV3.HeaderValueOption, 0)

	for key, value := range newHeader {
		hvo := &coreV3.HeaderValueOption{}

		hvo.Header = &coreV3.HeaderValue{
			Key:   key,
			Value: value,
		}

		hvo.Append = &wrappers.BoolValue{Value: false}
		hos = append(hos, hvo)
	}

	return &authV3.CheckResponse{
		Status: &status.Status{
			Code: int32(rpc.OK),
		},
		HttpResponse: &authV3.CheckResponse_OkResponse{
			OkResponse: &authV3.OkHttpResponse{
				Headers: hos,
			},
		},
	}, nil
}
