package biz

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/sunmi-OS/apidoc/pkg/ms"
	"github.com/tidwall/gjson"
)

// SyncApiDoc 同步接口文档到文档管理系统
func SyncApiDoc(plateform string) error {
	conf, err := readConfig()
	if err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	filename := pwd + "/docs/swagger.json"
	_, err = os.Stat(filename)
	if err != nil {
		filename = pwd + "/app/docs/swagger.json"
		_, err = os.Stat(filename)
		if err != nil {
			return err
		}
	}
	msConfig := conf.MS
	msClient := ms.NewClient(msConfig.Host, msConfig.AccessKey, msConfig.SecretKey)
	if plateform == "" { // upload to yapi and ms
		err = uploadToYAPI(filename, conf.YAPI)
		if err != nil {
			return err
		}
		err = uploadToMS(filename, msClient, msConfig)
		if err != nil {
			return err
		}
	} else if plateform == "yapi" { // upload to yapi
		err = uploadToYAPI(filename, conf.YAPI)
		if err != nil {
			return err
		}
	} else { // upload to metersphere
		err = uploadToMS(filename, msClient, msConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

// UploadToYAPI upload apidoc to yapi
// filename named file for uploading
func uploadToYAPI(filename string, conf yapiConf) error {
	importUrl := conf.Host + "/api/open/import_data"
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	urlValues := url.Values{}
	urlValues.Set("type", "swagger")
	urlValues.Set("merge", "merge")
	urlValues.Set("token", conf.Token)
	urlValues.Set("json", string(b))
	requestBody := strings.NewReader(urlValues.Encode())
	httpRequst, _ := http.NewRequest("POST", importUrl, requestBody)
	httpRequst.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpRequst)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("faild to upload to yapi")
	}
	response, _ := io.ReadAll(resp.Body)
	// {"errcode":0,"errmsg":"成功导入接口 2 个, 已存在的接口 2 个","data":null}
	// {"errcode":40011,"errmsg":"请登录...","data":null}
	if gjson.Get(string(response), "errcode").Int() != 0 {
		return errors.New("faild to upload to yapi")
	}
	return nil
}

// uploadToMS upload apidoc to metersphere
// filename named file for uploading
func uploadToMS(filename string, client *ms.Client, conf msConf) error {
	workspaceId, err := client.GetWorkspaceId(conf.Workspace)
	if err != nil {
		return err
	}
	projectId, err := client.GetProjectId(workspaceId, conf.Project)
	if err != nil {
		return err
	}
	applicationId, err := client.GetApplicationId(projectId, conf.Application)
	if err != nil {
		return err
	}
	err = client.APIDefinitionImport(filename, projectId, conf.Application, applicationId)
	if err != nil {
		return err
	}
	return nil
}
