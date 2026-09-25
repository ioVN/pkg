// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package options

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type Option struct {
	Client  *http.Client
	Method  string // http.MethodX
	Queries map[string]string
	Forms   map[string][]interface{}
	Header  http.Header
	Body    io.Reader

	raw []byte
}

func (op *Option) SetMethod(method string) *Option {
	switch method {
	case http.MethodConnect, http.MethodHead, http.MethodOptions,
		http.MethodDelete, http.MethodPatch, http.MethodPut, http.MethodTrace,
		http.MethodGet, http.MethodPost:
		op.Method = method
	}
	return op
}

func (op *Option) SetBasicAuth(username, password string) *Option {
	b := bytes.NewBufferString(username + ":" + password)
	value := base64.StdEncoding.EncodeToString(b.Bytes())
	op.Header.Set("Authorization", "Basic "+value)
	return op
}

func (op *Option) SetBearerAuth(token string) *Option {
	op.Header.Set("Authorization", "Bearer "+token)
	return op
}

func (op *Option) Query(key, value string) *Option {
	op.Queries[key] = value
	return op
}

func (op *Option) AddFormValue(key string, value string) *Option {
	if values := op.Forms[key]; len(values) == 0 {
		op.Forms[key] = make([]interface{}, 0)
	}
	op.Forms[key] = append(op.Forms[key], value)
	return op
}

func (op *Option) AddFormFile(key string, value *FormFile) *Option {
	if forms := op.Forms[key]; len(forms) == 0 {
		op.Forms[key] = make([]interface{}, 0)
	}
	op.Forms[key] = append(op.Forms[key], value)
	return op
}

func (op *Option) SubmitFormData() *Option {
	var (
		payload    = bytes.NewBuffer(nil)
		writer     = multipart.NewWriter(payload)
		fieldCount = 0
	)
	defer func(writer *multipart.Writer) {
		_ = writer.Close()
	}(writer)
	if len(op.Forms) == 0 {
		log.Printf("Log-Debug:\nForms empty\n")
		return op
	}
	for key, values := range op.Forms {
		for _, value := range values {
			switch data := value.(type) {
			case string:
				if w, _ := writer.CreateFormField(key); w != nil {
					_, _ = w.Write(bytes.NewBufferString(data).Bytes())
					fieldCount++
				}
			case *FormFile:
				if w, _ := writer.CreateFormFile(key, data.Filename); w != nil {
					_, _ = w.Write(data.File)
					fieldCount++
				}
			default:
				println("Type: Unknown", "; Key:", key, "; Value:", fmt.Sprintf("%+v", data))
			}
		}
	}
	if fieldCount > 0 {
		op.Header.Set("Content-Type", writer.FormDataContentType())
	}
	op.Method = http.MethodPost
	op.Body = payload
	return op
}

func (op *Option) SetData(body []byte) *Option {
	if len(body) != 0 {
		op.raw = make([]byte, len(body))
		copy(op.raw, body)
		op.Body = bytes.NewReader(op.raw)
	}
	return op
}

func (op *Option) SetJSON(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	op.Header.Set("Content-Type", "application/json")
	op.SetData(b)
	return nil
}

func CurlOption() *Option {
	return &Option{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		Method:  http.MethodGet,
		Queries: make(map[string]string),
		Forms:   make(map[string][]interface{}),
		Header:  http.Header{},
	}
}
