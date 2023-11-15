package ms

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
	"time"

	"github.com/sunmi-OS/gocore/v2/utils/cryption/aes"
	"github.com/tidwall/gjson"
)

type Client struct {
	host      string
	accessKey string
	secretKey string
}

func NewClient(host, accesKey, secretKey string) *Client {
	c := &Client{
		host:      host,
		accessKey: accesKey,
		secretKey: secretKey,
	}
	return c
}

// APIDefinitionImport import api doc
func (c *Client) APIDefinitionImport(filename, projectId, application, applicationId string) error {
	url := c.host + "/api/definition/import"
	signature, err := c.signature()
	if err != nil {
		return err
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			"file", "doc.json"))
	h.Set("Content-Type", "application/json")
	part, _ := writer.CreatePart(h)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	h2 := make(textproto.MIMEHeader)
	h2.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			"request", "blob"))
	h2.Set("Content-Type", "application/json")
	part2, _ := writer.CreatePart(h2)
	jsonBody := fmt.Sprintf(
		`{"modeId":"fullCoverage","moduleId":"%s","coverModule":false,"modulePath":"/%s","platform":"Swagger2","saved":true,"model":"definition","projectId":"%s","protocol":"HTTP"}`,
		applicationId,
		application,
		projectId,
	)
	_, err = part2.Write([]byte(jsonBody))
	if err != nil {
		return err
	}
	_ = writer.Close()
	httpRequst, _ := http.NewRequest("POST", url, body)
	httpRequst.Header.Add("Content-Type", writer.FormDataContentType())
	httpRequst.Header.Add("accessKey", c.accessKey)
	httpRequst.Header.Add("signature", signature)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpRequst)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("request ms error")
	}
	return nil
}

// Signature get signature
func (c *Client) signature() (string, error) {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	b, err := aes.EncryptUseCBC([]byte(c.accessKey+"|"+ts), []byte(c.secretKey), []byte(c.accessKey))
	if err != nil {
		return "", err
	}
	sign := base64.StdEncoding.EncodeToString(b)
	return sign, nil
}

// GetWorkspaceId get worksapce ID by name
func (c *Client) GetWorkspaceId(workspaceName string) (string, error) {
	url := c.host + "/workspace/list/userworkspace"
	signature, err := c.signature()
	if err != nil {
		return "", err
	}
	httpRequst, _ := http.NewRequest("GET", url, nil)
	httpRequst.Header.Add("accessKey", c.accessKey)
	httpRequst.Header.Add("signature", signature)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpRequst)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("request ms error")
	}
	response, _ := io.ReadAll(resp.Body)
	workspaceId := ""
	gjson.Get(string(response), "data").ForEach(func(key, value gjson.Result) bool {
		if value.Get("name").String() == workspaceName {
			workspaceId = value.Get("id").String()
			return false
		}
		return true
	})
	if workspaceId == "" {
		return "", errors.New("workspace not found")
	}
	return workspaceId, nil
}

// GetProjectId get project ID by workspace ID and project name
func (c *Client) GetProjectId(workspaceId, projectName string) (string, error) {
	url := c.host + "/project/listAll/" + workspaceId
	signature, err := c.signature()
	if err != nil {
		return "", err
	}
	httpRequst, _ := http.NewRequest("GET", url, nil)
	httpRequst.Header.Add("accessKey", c.accessKey)
	httpRequst.Header.Add("signature", signature)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpRequst)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("request ms error")
	}
	response, _ := io.ReadAll(resp.Body)
	projectId := ""
	gjson.Get(string(response), "data").ForEach(func(key, value gjson.Result) bool {
		if value.Get("name").String() == projectName {
			projectId = value.Get("id").String()
			return false
		}
		return true
	})
	if projectId == "" {
		return "", errors.New("project not found")
	}
	return projectId, nil
}

// GetApplicationId get application(module in ms) ID by project ID and application name
func (c *Client) GetApplicationId(projectId, applicationName string) (string, error) {
	url := c.host + fmt.Sprintf("/api/module/list/%s/HTTP", projectId)
	signature, err := c.signature()
	if err != nil {
		return "", err
	}
	httpRequst, _ := http.NewRequest("GET", url, nil)
	httpRequst.Header.Add("accessKey", c.accessKey)
	httpRequst.Header.Add("signature", signature)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpRequst)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("request ms error")
	}
	response, _ := io.ReadAll(resp.Body)
	applicationId := ""
	gjson.Get(string(response), "data").ForEach(func(key, value gjson.Result) bool {
		if value.Get("name").String() == applicationName {
			applicationId = value.Get("id").String()
			return false
		}
		return true
	})
	if applicationId == "" {
		return "", errors.New("application not found")
	}
	return applicationId, nil
}
