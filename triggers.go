package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type TriggerFunction struct {
	ClassName    string `json:"className,omitempty"`
	FunctionName string `json:"functionName,omitempty"`
	URL          string `json:"url,omitempty"`
}

func (c *Client) GetTriggerFunctions() ([]*TriggerFunction, error) {
	uri := fmt.Sprintf("/1/triggers/functions")
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
		TriggerFunctions []*TriggerFunction `json:"results"`
	}
	c.trace("GetTriggerFunctions", uri, string(body))
	return result.TriggerFunctions, json.Unmarshal(body, &result)
}

func (c *Client) CreateTriggerFunction(fn *TriggerFunction) error {
	payload, err := json.Marshal(fn)
	c.trace("CreateTriggerFunction >", "/1/triggers/functions", string(payload))
	resp, err := c.doWithBody("POST", "/1/triggers/functions", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var response interface{}
	err = json.Unmarshal(body, &response)
	c.trace("CreateUser <", "/1/triggers/functions", string(body))
	return err
}
