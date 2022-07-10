package provider

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const HostURL string = "https://businessapi.mosyle.com/v1"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Auth       AuthStruct
	Version    string
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func MosyleClient(version string, username, password, token *string) (*Client, error) {
	c := Client{
		HTTPClient: http.DefaultClient,
		// Default Hashicups URL
		HostURL: HostURL,
		Version: version,
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
		Token:    *token,
	}

	return &c, nil
}

func (c *Client) doBaseRequest(req *http.Request) ([]byte, error) {
	auth := c.Auth

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accesstoken", auth.Token)
	req.Header.Add("Authorization", "Basic "+auth.getAuth())
	req.Header.Add("User-Agent", "terraform-provider-mosyle "+c.Version)

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, errors.New("API response code " + fmt.Sprint(response.StatusCode) + " indicates failure")
	}

	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *Client) doRequest(req *http.Request) (ListResponse, error) {
	bytes, err := c.doBaseRequest(req)
	if err != nil {
		return ListResponse{}, err
	}

	response_obj := ListResponse{}
	err = json.Unmarshal(bytes, &response_obj)
	if err != nil {
		return ListResponse{}, err
	}

	if response_obj.Status != "OK" {
		return ListResponse{}, errors.New("Non succesful API call")
	}

	return response_obj, nil
}

func (a *AuthStruct) getAuth() string {
	base := a.Username + ":" + a.Password
	return base64.StdEncoding.EncodeToString([]byte(base))
}

type ListResponse struct {
	Status   string `json:"status"`
	Response []struct {
		Devices    []map[string]interface{} `json:"devices,omitempty"`
		Users      []map[string]interface{} `json:"users,omitempty"`
		UserGroups []map[string]interface{} `json:"usergroups,omitempty"`
		Rows       int                      `json:"rows"`
		PageSize   int                      `json:"page_size"`
		Page       int                      `json:"page"`
	} `json:"response"`
}

type DeviceGroupListResponse struct {
	Status   string `json:"status"`
	Response struct {
		DeviceGroups []map[string]interface{} `json:"devicegroups,omitempty"`
		Rows         int                      `json:"rows"`
		PageSize     int                      `json:"page_size"`
		Page         int                      `json:"page"`
	} `json:"response"`
}

type ListPostBody struct {
	Operation string                 `json:"operation"`
	Options   map[string]interface{} `json:"options"`
}

type ListUserBody struct {
	Operation string              `json:"operation"`
	Options   map[string][]string `json:"options"`
}
