// Package ambient is a client library for Ambient (http://ambidata.io)

package ambient

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/franela/goreq"
)

type Client struct {
	ChannelID int
	WriteKey  string
	UserKey   string
	ReadKey   string

	endpoint string
}

type ClientOption func(*Client) error

// UserKey specifies user key.
func UserKey(v string) ClientOption {
	return func(c *Client) error {
		c.UserKey = v
		return nil
	}
}

// ReadKey specifies read key.
func ReadKey(v string) ClientOption {
	return func(c *Client) error {
		c.ReadKey = v
		return nil
	}
}

// NewClient creates Client for Ambient.
func NewClient(channelId int, writeKey string, opts ...ClientOption) *Client {

	client := &Client{
		ChannelID: channelId,
		WriteKey:  writeKey,
	}

	client.endpoint = fmt.Sprintf("https://ambidata.io/api/v2/channels/%v", client.ChannelID)

	for _, opt := range opts {
		opt(client)
	}
	return client
}

// DataPoint represents a point in series.
type DataPoint map[string]interface{}

// NewDataPoint creates a data point. If t is provided, 'created' field will be sent.
func NewDataPoint(t ...time.Time) DataPoint {

	dp := DataPoint{}
	if len(t) > 0 {
		dp["created"] = t[0].Unix() * 1000
	}
	return dp
}

// Send sends multiple data points to Ambient.
func (c *Client) Send(points ...DataPoint) error {

	args := map[string]interface{}{}
	args["writeKey"] = c.WriteKey
	args["data"] = points

	req := goreq.Request{
		Uri:         c.endpoint + "/dataarray",
		Method:      "POST",
		ContentType: "application/json",
		Body:        args,
	}
	if res, err := req.Do(); err != nil {
		return err
	} else {
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf(res.Status)
		}
		return nil
	}
}

type ReadArgument struct {
	query url.Values
}

type ReadOption func(*ReadArgument)

// Date adds 'date' field to read request.
func Date(t time.Time) ReadOption {
	return func(arg *ReadArgument) {
		arg.query.Add("date", t.Format("2006-01-02"))
	}
}

// Range adds 'start' and 'end' fields to read request
func Range(start, end time.Time) ReadOption {
	return func(arg *ReadArgument) {
		arg.query.Add("start", start.Format("2006-01-02 15:04:05"))
		arg.query.Add("end", end.Format("2006-01-02 15:04:05"))
	}
}

// Count adds 'n' field to read request.
func Count(n int) ReadOption {
	return func(arg *ReadArgument) {
		arg.query.Add("n", strconv.Itoa(n))
	}
}

// Skip adds 'skip' field to read request.
func Skip(skip int) ReadOption {
	return func(arg *ReadArgument) {
		arg.query.Add("skip", strconv.Itoa(skip))
	}
}

// Read returns data points from Ambient.
func (c *Client) Read(opts ...ReadOption) ([]map[string]interface{}, error) {

	arg := &ReadArgument{query: make(url.Values)}
	if c.ReadKey != "" {
		arg.query.Add("readKey", c.ReadKey)
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(arg)
	}

	req := goreq.Request{
		Uri:         c.endpoint + "/data",
		QueryString: arg.query,
	}

	if res, err := req.Do(); err != nil {
		return nil, err
	} else {

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf(res.Status)
		}

		var values []map[string]interface{}
		res.Body.FromJsonTo(&values)

		// just reverse slice
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
		return values, nil
	}
}

// GetProp reaturns properties of a channel.
func (c *Client) GetProp() (map[string]interface{}, error) {

	query := url.Values{}
	if c.ReadKey != "" {
		query.Add("readKey", c.ReadKey)
	}

	req := goreq.Request{
		Uri:         c.endpoint,
		QueryString: query,
	}

	if res, err := req.Do(); err != nil {
		return nil, err
	} else {
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf(res.Status)
		}

		var values map[string]interface{}
		res.Body.FromJsonTo(&values)
		return values, nil
	}
}
