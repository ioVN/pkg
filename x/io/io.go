// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package io

import (
	"bytes"
)

type io struct {
	r []byte
}

func (ins *io) Set(src []byte) {
	ins.r = make([]byte, len(src))
	copy(ins.r, src)
}

func (ins *io) Get() []byte {
	dst := make([]byte, len(ins.r))
	copy(dst, ins.r)
	return dst
}

func (ins *io) String() string {
	dst := make([]byte, len(ins.r))
	copy(dst, ins.r)
	return bytes.NewBuffer(dst).String()
}

type IO interface {
	Set(src []byte)
	Get() []byte
	String() string
}

func New() IO {
	return &io{}
}
