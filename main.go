package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/helloworlde/cos/tool"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	events "github.com/tencentyun/scf-go-lib/events"
)

func main() {
	cloudfunction.Start(operate)
}

func operate(request events.APIGatewayRequest) (string, error) {
	requestContent := request.Body
	if requestContent == "" {
		fmt.Println("请求内容为空")
		return generateResponse(false, "请求内容为空", nil), nil
	}

	req := Request{}
	err := json.Unmarshal([]byte(requestContent), &req)
	if err != nil {
		fmt.Println("反序列化请求 Body 失败: ", err)
		return generateResponse(false, "反序列化请求 Body 失败", nil), nil
	}

	fmt.Println("请求 Body 内容: ", req)
	token := os.Getenv("COS_TOKEN")

	if req.Token != token {
		fmt.Println("Token 不正确: ", req.Token)
		return generateResponse(false, "UN_AUTHENTICATION", nil), nil
	}

	if req.Action != "UPLOAD" && req.Action != "DOWNLOAD" {
		fmt.Println("Action 不正确: ", req.Action)
		return generateResponse(false, "Action must be UPLOAD or DOWNLOAD", nil), nil

	}

	if req.Domain == "" {
		fmt.Println("Domain/Cookie 不正确, Domain: ", req.Domain, " Cookie: ", req.Cookie)
		return generateResponse(false, "Domain 不正确", nil), nil

	}

	var resp = Response{}

	fileName := getFileName(req.Domain)

	if "DOWNLOAD" == req.Action {
		resp, err = executeDownload(fileName)
	} else {
		resp, err = executeUpload(fileName, req.Cookie)
	}

	if err != nil {
		return generateResponse(false, err.Error(), nil), nil
	}

	return generateResponse(true, "SUCCESS", &resp), nil

}

func executeUpload(name string, cookie string) (Response, error) {
	err := tool.Upload(name, cookie)
	if err != nil {
		return Response{}, err
	}

	resp := Response{}
	return resp, nil
}

func getFileName(domain string) string {
	prefix := os.Getenv("COS_PATH")
	name := fmt.Sprintf("%s/%s", prefix, domain)
	return name
}

func executeDownload(name string) (Response, error) {
	result, err := tool.Download(name)
	if err != nil {
		return Response{}, err
	}

	resp := Response{
		Domain: name,
		Cookie: result,
	}
	return resp, nil
}

func generateResponse(success bool, message string, body *Response) string {
	var content = ""
	if body == nil {
		body = &Response{}
	}

	body.Success = success
	body.Message = message

	bodyBytes, _ := json.Marshal(body)
	content = string(bodyBytes)

	return content
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
