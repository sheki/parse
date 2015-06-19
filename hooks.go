package parse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type HookFunction struct {
	FunctionName string `json:"functionName,omitempty"`
	URL          string `json:"url,omitempty"`
}

func (c *Client) GetHookFunctions() ([]*HookFunction, error) {
	uri := fmt.Sprintf("/1/hooks/functions")
	resp, err := c.doSimple("GET", uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result struct {
		HookFunctions []*HookFunction `json:"results"`
	}
	c.trace("GetHookFunctions", uri, string(body))
	return result.HookFunctions, json.Unmarshal(body, &result)
}
