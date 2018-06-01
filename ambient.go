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

func UserKey(v string) ClientOption {
	return func(c *Client) error {
		c.UserKey = v
		return nil
	}
}

func ReadKey(v string) ClientOption {
	return func(c *Client) error {
		c.ReadKey = v
		return nil
	}
}

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

type DataPoint map[string]interface{}

func NewDataPoint(t ...time.Time) DataPoint {

	dp := DataPoint{}
	if len(t) > 0 {
		dp["created"] = t[0].Unix() * 1000
	}
	return dp
}

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

func Date(t time.Time) ReadOption {
	return func(arg *ReadArgument) {
		arg.query.Add("date", t.Format("2006-01-02"))
	}
}

func Range(start, end time.Time) ReadOption {
	return func(arg *ReadArgument) {
		arg.query.Add("start", start.Format("2006-01-02 15:04:05"))
		arg.query.Add("end", end.Format("2006-01-02 15:04:05"))
	}
}

func Count(n int) ReadOption {
	return func(arg *ReadArgument) {
		arg.query.Add("n", strconv.Itoa(n))
	}
}

func Skip(skip int) ReadOption {
	return func(arg *ReadArgument) {
		arg.query.Add("skip", strconv.Itoa(skip))
	}
}

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
