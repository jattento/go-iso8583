package iso8583server_test

import (
	"encoding/binary"
	"fmt"
	"github.com/jattento/go-iso8583/pkg/iso8583server"
	"github.com/jattento/go-iso8583/pkg/mti"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"testing"
	"time"
)

func TestSimpleHandler(t *testing.T) {
	go func() {
		time.Sleep(time.Minute)
		t.Fatal("test timeout")
	}()

	sv, err := iso8583server.New()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	expectedCalls := 3
	sv.AddTopPriorityHandler(func(connection io.WriteCloser, message []byte) { expectedCalls-- }, func(mtiValue mti.MTI) bool {
		return mtiValue.Version() == mti.Version8583To1993
	})

	go sv.Start()
	defer func() {
		assert.Nil(t, sv.Close())
	}()

	cn, err := net.Dial("tcp", "localhost:8080")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	defer func() {
		assert.Nil(t, cn.Close())
	}()

	bLen := make([]byte, 2)
	binary.BigEndian.PutUint16(bLen, 10)
	_, err = cn.Write(bLen)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	_,err =cn.Write([]byte("1000111111"))
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	_, err = cn.Write(bLen)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	_,err =cn.Write([]byte("1000111111"))
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	time.Sleep(5 * time.Second)

	_, err = cn.Write(bLen)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	_,err =cn.Write([]byte("1000111111"))
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	for expectedCalls != 0{
		// Wait...
	}
	fmt.Println(expectedCalls)
}
