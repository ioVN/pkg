// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package net

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type MultipartForm struct {
	*multipart.Form
}

func (ins *MultipartForm) GetForm(key string) []string {
	var list = make([]string, len(ins.Value[key]))
	copy(list, ins.Value[key])
	return list
}

func (ins *MultipartForm) GetFile(key string) []*FormFile {
	var (
		list = make([]*FormFile, 0)
	)
	for keyName := range ins.File {
		log.Printf("MultipartForm form key available: %s", keyName)
	}
	for _, file := range ins.File[key] {
		var (
			temp = bytes.NewBuffer(nil)
		)
		content, err := file.Open()
		if err != nil {
			log.Printf("MultipartForm open error: %s", err.Error())
			continue
		}
		if _, err := io.Copy(temp, content); err != nil {
			log.Printf("MultipartForm copy error: %s", err.Error())
		}
		_ = content.Close()
		list = append(list, &FormFile{Filename: file.Filename, File: temp.Bytes()})
		log.Printf("MultipartForm get %s to filename: %s", key, file.Filename)
	}
	return list
}

func NewMultipartForm(c *gin.Context) (*MultipartForm, error) {
	if err := c.Request.ParseMultipartForm(0); err != nil {
		return nil, err
	}
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	return &MultipartForm{form}, nil
}

func MakeMultipartForm(c *gin.Context) (*MultipartForm, error) {
	if err := c.Request.ParseMultipartForm(0); err != nil {
		return nil, err
	}
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	return &MultipartForm{form}, nil
}
