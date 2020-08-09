package iso8583server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/jattento/go-iso8583/pkg/mti"
)

type Server struct {
	config       configuration
	handlerRules []handlerRule
}

type configuration struct {
	Listener       net.Listener
	LogInfo        LogFunc
	LogErr         LogFunc
	NetRead        NetReadFunc
	ReadMTI        ReadMTIFunc
	UnknownHandler HandlerFunc

	// Timeout is the max time the server waits for new messages before closing the connection
	Timeout time.Duration
}

type handlerRule struct {
	rule    Rule
	handler HandlerFunc
}

type ReadMTIFunc = func([]byte) (mti.MTI, error)

type LogFunc = func(v ...interface{})

type HandlerFunc = func(connection io.WriteCloser, message []byte)

type NetReadFunc = func(io.Reader) (chan []byte, chan error)

type Rule = func(mti.MTI) bool

func New(options ...Option) (*Server, error) {
	config := configuration{}
	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}

	if config.Listener == nil {
		l, err := net.Listen("tcp", ":8080")
		if err != nil {
			return nil, err
		}

		config.Listener = l
	}

	if config.Timeout == 0 {
		config.Timeout = 3 * time.Minute
	}

	if config.NetRead == nil {
		config.NetRead = DefaultReader{
			SizeChunkLength: 2,
			SizeReader:      func(b []byte) int { return int(binary.BigEndian.Uint16(b)) },
		}.Read
	}

	if config.ReadMTI == nil {
		config.ReadMTI = DefaultReadMTI
	}

	if config.UnknownHandler == nil {
		// nop handler
		config.UnknownHandler = func(connection io.WriteCloser, message []byte) { return }
	}

	if config.LogInfo == nil {
		config.LogInfo = log.New(os.Stdout, "iso8583server_info ", log.LstdFlags).Print
	}

	if config.LogErr == nil {
		config.LogErr = log.New(os.Stderr, "iso8583server_error ", log.LstdFlags).Print
	}

	return &Server{config: config}, nil
}

func (server *Server) handleConnection(c net.Conn) {
	extendReadDeadline := func() {
		if err := c.SetReadDeadline(time.Now().Add(server.config.Timeout)); err != nil {
			server.config.LogErr(fmt.Sprintf("failed setting read deadline: %v", handleErrorNet(err, c)))
			return
		}
	}

	server.config.LogInfo(fmt.Sprintf("serving %s\n", c.RemoteAddr().String()))

	extendReadDeadline()

	messageChan, errChan := server.config.NetRead(c)

	for {
		select {
		case message := <-messageChan:
			func(){
				extendReadDeadline()

				mtiValue, err := server.config.ReadMTI(message)
				if err != nil {
					server.config.LogErr("cant read MTI from message: ", err)
					go server.config.UnknownHandler(c, message)
					return
				}

				for _, hr := range server.handlerRules {
					if hr.rule(mtiValue) {
						go hr.handler(c, message)
						return
					}
				}

				server.config.LogErr("no handler defined for MTI ", mtiValue, " sending to unknown manager")
				go server.config.UnknownHandler(c, message)
			}()
		case err := <-errChan:
			server.config.LogErr(c.RemoteAddr().String(),": ",handleErrorNet(err, c))
			return
		}
	}
}

func (server *Server) Close() error {
	return server.config.Listener.Close()
}

func (server *Server) AddBottomPriorityHandler(handler HandlerFunc, rules ...Rule) {
	if handler == nil {
		return
	}

	for _, rule := range rules {
		if rule != nil {
			server.handlerRules = append(server.handlerRules, handlerRule{rule: rule, handler: handler})
		}
	}
}

func (server *Server) AddTopPriorityHandler(handler HandlerFunc, rules ...Rule) {
	if handler == nil {
		return
	}

	for _, rule := range rules {
		if rule != nil {
			server.handlerRules = append([]handlerRule{{rule: rule, handler: handler}}, server.handlerRules...)
		}
	}
}

func (server *Server) Start() {
	for range time.NewTicker(time.Second).C {
		c, err := server.config.Listener.Accept()
		if err != nil {
			server.config.LogErr(fmt.Sprintf("cant accept incoming connection %v", err))
			continue
		}

		go server.handleConnection(c)
	}
}

func DefaultReadMTI(b []byte) (mti.MTI, error) {
	if len(b) < 4 {
		return 0, errors.New("message to short to read MTI")
	}

	iMTI, err := strconv.Atoi(string(b[:4]))
	if err != nil {
		return 0, fmt.Errorf("first 4 characters arent numbers: %v", err)
	}

	return mti.MTI(iMTI), nil
}

type DefaultReader struct {
	// SizeChunkLength are the amount of bytes that must be read at message start to determine length
	SizeChunkLength int

	// SizeReader gets as input SizeChunkLength bytes and must return the represented length.
	SizeReader func([]byte) int
}

func (reader DefaultReader) Read(net io.Reader) (chan []byte, chan error) {
	messageChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		var sizeBuffer []byte
		var messageBuffer []byte

		upToFillSizeBuffer := func() int { return reader.SizeChunkLength - len(sizeBuffer) }
		upToFillMessageBuffer := func() int { return reader.SizeReader(sizeBuffer) - len(messageBuffer) }

		for {
			if upToFillSizeBuffer() > 0 {
				sizeChunk := make([]byte, upToFillSizeBuffer())

				n, err := net.Read(sizeChunk)
				if err != nil {
					errChan <- err
					return
				}

				if n > 0 {
					sizeBuffer = append(sizeBuffer, sizeChunk[:n]...)
				}

				continue
			}

			if upToFillMessageBuffer() > 0 {
				messageChunk := make([]byte, upToFillMessageBuffer())

				n, err := net.Read(messageChunk)
				if err != nil {
					errChan <- err
					return
				}

				if n > 0 {
					messageBuffer = append(messageBuffer, messageChunk[:n]...)
				}

				continue
			}

			message := make([]byte, len(messageBuffer))
			copy(message, messageBuffer)

			messageChan <- message

			// Restart all buffer
			sizeBuffer = nil
			messageBuffer = nil
		}
	}()

	return messageChan, errChan
}

func handleErrorNet(err error, c net.Conn) error {
	if clErr := c.Close(); clErr != nil {
		return fmt.Errorf("%v and failed closing connection: %v", err, clErr)
	}
	return err
}
