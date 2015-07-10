package parse

import (
	"errors"
	"fmt"
)

// Zero time in the format parse understands
const zeroTime = "0000-00-00T00:00:00.000Z"

// NewIter returns an iterator to iterate over a pars e class
func NewIter(client *Client, className string) (*Iter, error) {
	return &Iter{client: client, currentTime: zeroTime, class: className}, nil
}

type Iter struct {
	client       *Client
	currentBatch []interface{}
	index        int
	currentTime  string
	class        string
	lastErr      error
}

func (i *Iter) fetchCurrent() ([]interface{}, error) {
	whereStr := fmt.Sprintf(`{"createdAt" : { "$gt" : "%s" } }`, i.currentTime)
	where := QueryOptions{Where: whereStr, Limit: 1000, Order: "createdAt"}
	var currentBatch []interface{}
	err := i.client.QueryClass(i.class, &where, &currentBatch)
	return currentBatch, err
}

// Next returns the object and the continue bool. If continue is false
// the interface should fail
func (i *Iter) Iter() (interface{}, bool) {
	if i.index == 0 || i.index >= len(i.currentBatch) {
		var err error
		i.currentBatch, err = i.fetchCurrent()
		if err != nil {
			i.lastErr = err
			return nil, false
		}
		i.index = 0
	}

	current := i.currentBatch[i.index]
	cmap, ok := current.(map[string]interface{})
	if !ok {
		i.lastErr = errors.New("parse returned non map object")
		return nil, false
	}
	createdAt, ok := cmap["createdAt"]
	if !ok {
		i.lastErr = errors.New("parse returned non object without id")
		return nil, false
	}
	createdAtStr, ok := createdAt.(string)
	if !ok {
		i.lastErr = errors.New("parse returned non object without string id")
		return nil, false
	}
	i.currentTime = createdAtStr
	i.index++
	i.lastErr = nil
	return current, true
}

// Error returns an error if there was an error in the last iteration.
// After Next returns false, err should be checked
func (i *Iter) Err() error {
	return i.lastErr
}
