package gossip

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type Codec interface {
	Encode(msg any) ([]byte, error)
	Decode(data []byte) (any, error)
}

type codec struct {
}

func NewCodec() Codec {
	return &codec{}
}

type envelope struct {
	Type    MessageType
	Payload []byte
}

func (c *codec) Encode(msg any) ([]byte, error) {
	msgType := getMessageType(msg)
	if msgType == 0 {
		return nil, errors.New("unknown message type")
	}

	// Encode message payload
	var payloadBuff bytes.Buffer
	payloadEnc := gob.NewEncoder(&payloadBuff)

	if err := payloadEnc.Encode(msg); err != nil {
		return nil, err
	}

	if len(payloadBuff.Bytes()) == 0 {
		return nil, errors.New("payload cannot be empty")
	}

	// wrap payload in envelope and Encode
	var envBuff bytes.Buffer
	envEnc := gob.NewEncoder(&envBuff)
	env := envelope{
		Type:    msgType,
		Payload: payloadBuff.Bytes(),
	}

	if err := envEnc.Encode(env); err != nil {
		return nil, err
	}

	return envBuff.Bytes(), nil
}

func (c *codec) Decode(data []byte) (any, error) {
	var env envelope
	envDec := gob.NewDecoder(bytes.NewReader(data))
	if err := envDec.Decode(&env); err != nil {
		return nil, err
	}

	msg, err := newMessageByType(env.Type)
	if err != nil {
		return nil, err
	}

	payloadDec := gob.NewDecoder(bytes.NewReader(env.Payload))
	if err := payloadDec.Decode(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func getMessageType(msg any) MessageType {
	switch msg.(type) {
	case *Ping:
		return MessageTypePing
	case *Ack:
		return MessageTypeAck
	case *PingReq:
		return MessageTypePingReq
	default:
		return 0
	}
}

func newMessageByType(messageType MessageType) (any, error) {
	switch messageType {
	case MessageTypePing:
		return &Ping{}, nil
	case MessageTypeAck:
		return &Ack{}, nil
	case MessageTypePingReq:
		return &PingReq{}, nil
	default:
		return nil, errors.New("unknown message type")
	}
}
