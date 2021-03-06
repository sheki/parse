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
	// Limit controls the maximum number of objects returned for a query
	Limit int
	// field(s) to order results by
	Order string
}

// QueryClass performs a lookup of objects based on query options.
func (c *Client) QueryClass(className string, options *QueryOptions, destination interface{}) error {
	uri, err := url.Parse(fmt.Sprintf("/1/classes/%s", className))

	if options != nil {
		params := uri.Query()
		if options.Where != "" {
			params.Set("where", options.Where)
		}
		if options.Limit != 0 {
			params.Set("limit", fmt.Sprint(options.Limit))
		}
		if options.Order != "" {
			params.Set("order", fmt.Sprint(options.Order))
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

// Query performs a lookup of objects based on query options. destination must be
// a slice of types satisfying the Object interface.
func (c *Client) Query(options *QueryOptions, destination interface{}) error {
	className, err := objectTypeNameFromSlice(destination)
	if err != nil {
		return err
	}
	return c.QueryClass(className, options, destination)
}
