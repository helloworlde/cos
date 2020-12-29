package tool

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var client = initClient()

func Upload(name string, content string) error {
	fmt.Println("开始上传文件: ", name, " 内容: ", content)
	f := strings.NewReader(content)

	_, err := client.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		fmt.Println("name: ", name, "上传文件失败")
	}
	return err
}

func Download(name string) (string, error) {
	fmt.Println("开始下载文件: ", name)
	resp, err := client.Object.Get(context.Background(), name, nil)
	if err != nil {
		fmt.Println("name: ", name, "获取文件失败")
		return "", nil
	}
	defer resp.Body.Close()

	contentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("name: ", name, "读取内容失败")
		return "", err
	}
	return string(contentBytes), nil
}

func initClient() *cos.Client {
	cosUrl := os.Getenv("COS_URL")
	secretId := os.Getenv("COS_SECRET_ID")
	secretKey := os.Getenv("COS_SECRET_KEY")

	u, _ := url.Parse(cosUrl)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})
	return c
}
