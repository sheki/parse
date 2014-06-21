package parse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

// QueryOptions represents the parameters to a Parse query.
type QueryOptions struct {
	Where string
}

// Query performs a lookup of objects based on query options. destination must be
// a slice of types satisfying the Object interface. 
func (c *Client) Query(options *QueryOptions, destination interface{}) error {
	className, err := objectTypeNameFromSlice(destination)
	if err != nil {
		return err
	}
	uri, err := url.Parse(fmt.Sprintf("/1/classes/%s", className))

	if options != nil {
		params := uri.Query()
		if options.Where != "" {
			params.Set("where", options.Where)
		}
		uri.RawQuery = params.Encode()
	}

	resp, err := c.doSimple("GET", uri.String())
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// delay parsing of results
	c.trace("Query", uri, string(body))
	results := struct {
		Results json.RawMessage `json:"results"`
	}{}
	// first pass
	err = json.Unmarshal(body, &results)
	if err != nil {
		return err
	}
	return json.Unmarshal(results.Results, destination)
}
