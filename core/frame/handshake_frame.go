package frame

import (
	"github.com/yomorun/y3"
)

// HandshakeFrame is a Y3 encoded.
type HandshakeFrame struct {
	// Name is client name
	Name string
	// ClientType represents client type (source or sfn)
	ClientType byte
	// auth
	authType    byte
	authPayload []byte
}

// NewHandshakeFrame creates a new HandshakeFrame.
func NewHandshakeFrame(name string, clientType byte, authType byte, authPayload []byte) *HandshakeFrame {
	return &HandshakeFrame{
		Name:        name,
		ClientType:  clientType,
		authType:    authType,
		authPayload: authPayload,
	}
}

// Type gets the type of Frame.
func (h *HandshakeFrame) Type() Type {
	return TagOfHandshakeFrame
}

// Encode to Y3 encoding.
func (h *HandshakeFrame) Encode() []byte {
	// name
	nameBlock := y3.NewPrimitivePacketEncoder(byte(TagOfHandshakeName))
	nameBlock.SetStringValue(h.Name)
	// type
	typeBlock := y3.NewPrimitivePacketEncoder(byte(TagOfHandshakeType))
	typeBlock.SetBytesValue([]byte{h.ClientType})
	// auth
	authTypeBlock := y3.NewPrimitivePacketEncoder(byte(TagOfHandshakeAuthType))
	authTypeBlock.SetBytesValue([]byte{h.authType})
	authPayloadBlock := y3.NewPrimitivePacketEncoder(byte(TagOfHandshakeAuthPayload))
	authPayloadBlock.SetBytesValue(h.authPayload)

	// handshake frame
	handshake := y3.NewNodePacketEncoder(byte(h.Type()))
	handshake.AddPrimitivePacket(nameBlock)
	handshake.AddPrimitivePacket(typeBlock)
	handshake.AddPrimitivePacket(authTypeBlock)
	handshake.AddPrimitivePacket(authPayloadBlock)

	return handshake.Encode()
}

// DecodeToHandshakeFrame decodes Y3 encoded bytes to HandshakeFrame.
func DecodeToHandshakeFrame(buf []byte) (*HandshakeFrame, error) {
	node := y3.NodePacket{}
	_, err := y3.DecodeToNodePacket(buf, &node)
	if err != nil {
		return nil, err
	}

	handshake := &HandshakeFrame{}
	// name
	if nameBlock, ok := node.PrimitivePackets[byte(TagOfHandshakeName)]; ok {
		name, err := nameBlock.ToUTF8String()
		if err != nil {
			return nil, err
		}
		handshake.Name = name
	}
	// type
	if typeBlock, ok := node.PrimitivePackets[byte(TagOfHandshakeType)]; ok {
		clientType := typeBlock.ToBytes()
		handshake.ClientType = clientType[0]
	}
	// auth type
	if authTypeBlock, ok := node.PrimitivePackets[byte(TagOfHandshakeAuthType)]; ok {
		authType := authTypeBlock.ToBytes()
		handshake.authType = authType[0]
	}
	// auth payload
	if authPayloadBlock, ok := node.PrimitivePackets[byte(TagOfHandshakeAuthPayload)]; ok {
		authPayload := authPayloadBlock.ToBytes()
		handshake.authPayload = authPayload
	}

	return handshake, nil
}

func (h *HandshakeFrame) AuthType() byte {
	return h.authType
}

func (h *HandshakeFrame) AuthPayload() []byte {
	return h.authPayload
}
