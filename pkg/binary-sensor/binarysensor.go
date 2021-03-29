package binarysensor

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spritkopf/esb-bridge/pkg/client"
	"github.com/spritkopf/esb-bridge/pkg/esbbridge"
)

const (
	cmdStateNotification = 0x91
	cmdGetChannel        = 0x92
	cmdSetChannel        = 0x93
	valueFalse           = 0x00
	valueTrue            = 0x01
)

type BinarySensor struct {
	esbClient     client.EsbClient
	esbAddress    []byte
	subscriptions []subscription
}

// Open opens a connection to the esb-bridge RPC server.
// Params:
//   - deviceAddress: ESB device pipeline address in format xx:xx:xx:xx:xx
//   - serverAddress: IP address of the esb-bridge server (e.g. "10.32.2.100", or "localhost")
//   - serverPort : Port of the esb-bridge server (e.g. 9815)
func (b *BinarySensor) Open(deviceAddress string, serverAddress string, serverPort uint) error {

	// decode target address to bytes
	addrBytes, err := hex.DecodeString(strings.ReplaceAll(deviceAddress, ":", ""))

	if err != nil {
		return fmt.Errorf("invalid format for deviceAddress: %v", err)
	}
	if len(addrBytes) != 5 {
		return fmt.Errorf("invalid length for deviceAddress: need 5, got %v", len(addrBytes))
	}
	b.esbAddress = addrBytes

	err = b.esbClient.Connect(fmt.Sprintf("%v:%v", serverAddress, serverPort))

	return err
}

// Close closes the connection to the esb-bridge RPC server
func (b *BinarySensor) Close() error {
	err := b.esbClient.Disconnect()

	return err
}

// SetValue sets the value of a specific channel
func (b *BinarySensor) SetValue(channel uint8, value bool) error {

	newVal := byte(valueFalse)

	if value {
		newVal = valueTrue
	}
	answerMsg, err := b.esbClient.Transfer(esbbridge.EsbMessage{Address: b.esbAddress, Cmd: cmdSetChannel, Payload: []byte{byte(channel), newVal}})

	if err != nil {
		return fmt.Errorf("ESB Transfer error: %v", err)
	}
	if answerMsg.Error != 0 {
		return fmt.Errorf("ESB answer has error: %v", answerMsg.Error)
	}
	return nil
}

// GetValue reads the current value of a specific channel
func (b *BinarySensor) GetValue(channel uint8) (bool, error) {

	answerMsg, err := b.esbClient.Transfer(esbbridge.EsbMessage{Address: b.esbAddress, Cmd: cmdGetChannel, Payload: []byte{byte(channel)}})

	if err != nil {
		return false, fmt.Errorf("ESB Transfer error: %v", err)
	}
	if answerMsg.Error != 0 {
		return false, fmt.Errorf("ESB answer has error: %v", answerMsg.Error)
	}
	channelVal := false

	if answerMsg.Payload[0] == valueTrue {
		channelVal = true
	}

	return channelVal, nil
}
