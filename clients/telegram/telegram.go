package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/untibullet/dailyhelper/tools/elog"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func NewClient(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) GetUpdates(offset int, limit int) (updates []Update, err error) {
	q := url.Values{}
	// (Itoa) Interger to ASCII
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	defer func() {
		err = elog.WrapIfErr("can`t get updates", err)
	}()

	data, err := c.sendRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chat_id int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chat_id))
	q.Add("text", text)

	_, err := c.sendRequest(sendMessageMethod, q)
	if err != nil {
		return elog.Wrap("can`t send message", err)
	}

	return nil
}

func (c *Client) sendRequest(method string, query url.Values) (data []byte, err error) {
	defer func() {
		err = elog.WrapIfErr("can`t send request", err)
	}()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
