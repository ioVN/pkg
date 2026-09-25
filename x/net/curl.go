// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package net

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ioVN/pkg/x/net/options"
)

const headerClientUnixTime = "X-Client-Unix"

// Curl - call to dst url
func Curl(ctx context.Context, dstUrl string, option *options.Option) (*Response, error) {
	if option == nil {
		option = options.CurlOption()
	}
	if strings.HasPrefix(dstUrl, "ws") {
		return nil, errors.New("the `websocket` protocol is not supported")
	} else {
		dstUrl = UrlPrettyParse(dstUrl)
	}
	for key, value := range option.Queries {
		if !strings.Contains(dstUrl, "?") {
			dstUrl = dstUrl + "?" + key + "=" + value
		} else {
			dstUrl = dstUrl + "&" + key + "=" + value
		}
	}
	request, err := http.NewRequestWithContext(ctx, option.Method, dstUrl, option.Body)
	if err != nil {
		return nil, err
	}
	for key, arr := range option.Header {
		for i, value := range arr {
			if i == 0 {
				request.Header.Set(key, value)
			} else {
				request.Header.Add(key, value)
			}
		}
	}
	request.Header.Set(
		headerClientUnixTime,
		fmt.Sprintf("%d", time.Now().Unix()),
	)
	defer func(t0 time.Time) {
		log.Printf("DEBUG: Request time of %s is %s", dstUrl, time.Since(t0).String())
	}(time.Now())
	response, err := option.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(r io.ReadCloser) {
		_ = r.Close()
	}(response.Body)
	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: response.StatusCode,
		Header:     response.Header.Clone(),
		body:       buffer,
	}, nil
}

type Response struct {
	StatusCode int
	Header     http.Header
	body       []byte
}

func (ins *Response) GetBody() ([]byte, error) {
	if len(ins.body) > 0 {
		return ins.body, nil
	}
	return nil, errors.New("body invalid")
}

func (ins *Response) Load(v interface{}) error {
	if len(ins.body) > 0 {
		return json.Unmarshal(ins.body, v)
	}
	return errors.New("body invalid")
}

func (ins *Response) String() string {
	return bytes.NewBuffer(ins.body).String()
}
