package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/helloworlde/cos/tool"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
)

func main() {
	cloudfunction.Start(operate)
}

func operate(request events.APIGatewayRequest) (string, error) {
	requestContent := request.Body
	if requestContent == "" {
		fmt.Println("请求内容为空")
		return fmt.Sprint("请求内容为空"), nil
	}

	req := Request{}
	err := json.Unmarshal([]byte(requestContent), &req)
	if err != nil {
		fmt.Println("反序列化请求 Body 失败: ", err)
		return "反序列化请求 Body 失败", err
	}

	fmt.Println("请求 Body 内容: ", req)
	token := os.Getenv("COS_TOKEN")

	if req.Token != token {
		resp := Response{
			Success: false,
			Message: "UN_AUTHENTICATION",
		}
		fmt.Println("Token 不正确: ", req.Token)
		content, _ := json.Marshal(resp)
		return string(content), errors.New(resp.Message)
	}

	if req.Action != "UPLOAD" && req.Action != "DOWNLOAD" {
		resp := Response{
			Success: false,
			Message: "Action must be UPLOAD or DOWNLOAD",
		}
		fmt.Println("Action 不正确: ", req.Action)
		content, _ := json.Marshal(resp)
		return string(content), errors.New(resp.Message)
	}

	if req.Domain == "" {
		resp := Response{
			Success: false,
			Message: "Invalid content",
		}
		fmt.Println("Domain/Cookie 不正确, Domain: ", req.Domain, " Cookie: ", req.Cookie)
		content, _ := json.Marshal(resp)
		return string(content), errors.New(resp.Message)
	}

	var resp = Response{}

	fileName := getFileName(req.Domain)

	if "DOWNLOAD" == req.Action {
		resp = executeDownload(fileName)
	} else {
		resp = executeUpload(fileName, req.Cookie)
	}

	content, _ := json.Marshal(resp)
	return string(content), nil
}

func executeUpload(name string, cookie string) Response {
	err := tool.Upload(name, cookie)
	if err != nil {
		resp := Response{
			Success: false,
			Message: err.Error(),
		}
		return resp
	}

	resp := Response{
		Success: true,
		Domain:  name,
		Cookie:  cookie,
	}
	return resp
}

func getFileName(domain string) string {
	prefix := os.Getenv("COS_PATH")
	name := fmt.Sprintf("%s/%s", prefix, domain)
	return name
}

func executeDownload(name string) Response {
	result, err := tool.Download(name)
	if err != nil {
		resp := Response{
			Success: false,
			Message: err.Error(),
		}
		return resp
	}

	resp := Response{
		Success: true,
		Domain:  name,
		Cookie:  result,
	}
	return resp
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Domain  string `json:"domain"`
	Cookie  string `json:"cookie"`
}

type Request struct {
	Token  string `json:"token"`
	Action string `json:"action"`
	Domain string `json:"domain"`
	Cookie string `json:"cookie"`
}
