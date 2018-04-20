package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type WindAngle int

const InvalidWindAngle WindAngle = 0

type WindAngleReference string

const InvalidWindAngleReference WindAngleReference = "X"

const (
	Relative WindAngleReference = "R"
	True     WindAngleReference = "T"
)

type WindSpeed float64

const InvalidWindSpeed WindSpeed = 0.0

type WindSpeedUnit string

const InvalidWindSpeedUnit WindSpeedUnit = "X"

const (
	KilometersPerHour WindSpeedUnit = "K"
	MetersPerSecond   WindSpeedUnit = "M"
	Knots             WindSpeedUnit = "N"
)

type MessageValidity string

const InvalidMessageValidity MessageValidity = "X"

const (
	Valid   MessageValidity = "A"
	Invalid MessageValidity = "V"
)

type Checksum byte

const InvalidChecksum Checksum = 0x00

type MWVMessage struct {
	dir                WindAngle
	ref                WindAngleReference
	spd                WindSpeed
	uni                WindSpeedUnit
	sta                MessageValidity
	chk                Checksum
	source             string
	calculatedChecksum Checksum
}

var IncorrectMessageHeader = errors.New("Incorrect message header")

func (wa *WindAngle) parse(str string) error {
	n, err := fmt.Sscanf(str, "%d", wa)

	if n != 1 {
		*wa = InvalidWindAngle
		return NewError(fmt.Sprintf("Could not parse wind angle from string '%s'", str))
	}

	*wa = *wa % 360

	if *wa < 0 {
		*wa = 360 + *wa
	}

	return err
}

func (war *WindAngleReference) parse(str string) error {
	n, err := fmt.Sscanf(str, "%1s", war)

	if n != 1 {
		return NewError(fmt.Sprintf("Could not parse wind angle reference from string '%s'", str))
	}

	switch *war {
	case Relative, True:
	default:
		err = NewError(fmt.Sprintf("Invalid angle reference '%s'", *war))
		*war = InvalidWindAngleReference
		return err
	}

	return err
}

func (ws *WindSpeed) parse(str string) error {
	n, err := fmt.Sscanf(str, "%f", ws)

	if n != 1 {
		return NewError(fmt.Sprintf("Could not parse wind speed from string '%s'", str))
	}

	if *ws < 0.0 {
		err = NewError(fmt.Sprintf("Parse('%s') resulted in a negative wind speed of %f", str, *ws))
		*ws = InvalidWindSpeed
		return err
	}

	return err
}

func (wsu *WindSpeedUnit) parse(str string) error {
	n, err := fmt.Sscanf(str, "%1s", wsu)

	if n != 1 {
		return NewError(fmt.Sprintf("Could not parse wind speed unit from string '%s'", str))
	}

	switch *wsu {
	case KilometersPerHour, MetersPerSecond, Knots:
	default:
		err = NewError(fmt.Sprintf("Parse('%s') => '%s': not a valid wind speed unit", str, *wsu))
		*wsu = InvalidWindSpeedUnit
		return err
	}

	return err
}

func (mv *MessageValidity) parse(str string) error {
	n, err := fmt.Sscanf(str, "%1s", mv)

	if n != 1 {
		return NewError(fmt.Sprintf("Could not parse message validity from string '%s'", str))
	}

	switch *mv {
	case Valid, Invalid:
	default:
		err = NewError(fmt.Sprintf("Parse('%s') => '%s': not a valid message validity", str, *mv))
		*mv = InvalidMessageValidity
		return err
	}

	return err
}

func (c *Checksum) parse(str string) error {
	if len(str) != 2 {
		*c = InvalidChecksum
		return NewError(fmt.Sprintf("String too long or short to be a checksum (want 2 characters) '%s'", str))
	}

	n, err := fmt.Sscanf(str, "%02X", c)

	if n != 1 {
		*c = InvalidChecksum
		return NewError(fmt.Sprintf("Could not parse checksum byte from string '%s'", str))
	}

	if *c < 0x00 || *c > 0xFF {
		err = NewError(fmt.Sprintf("Parse('%s') => '%X': not a valid message validity", str, *c))
		*c = InvalidChecksum
		return err
	}

	return err
}

func checksumFromString(str string) (c Checksum) {
	b := []byte(str)
	for i := 0; i < len(b); i++ {
		c ^= Checksum(b[i])
	}
	return
}

func (msg *MWVMessage) parse(str string) (err error) {
	msg.source = str

	fields := strings.Split(msg.source, ",")

	if len(fields) != 6 {
		return NewError(fmt.Sprintf("Not correct number of fields in MWV-message; want %d but got %d in string '%s'", 5, len(fields), str))
	}

	if fields[0] != "$WIMWV" {
		return IncorrectMessageHeader
	}

	validationParts := strings.Split(fields[5], "*")

	err = msg.dir.parse(fields[1])          // WindAngle
	err = msg.ref.parse(fields[2])          // WindAngleReference
	err = msg.spd.parse(fields[3])          // WindSpeed
	err = msg.uni.parse(fields[4])          // WindSpeedUnit
	err = msg.sta.parse(validationParts[0]) // MessageValidity
	err = msg.chk.parse(validationParts[1]) // Checksum

	if err != nil {
		return
	}

	var stringToChecksum string
	stringToChecksum = strings.TrimLeft(msg.source, "$")
	stringToChecksum = strings.Split(stringToChecksum, "*")[0]
	msg.calculatedChecksum = checksumFromString(stringToChecksum)

	if msg.calculatedChecksum != msg.chk {
		return NewError(fmt.Sprintf("Parsed message could not be validated; checksum 0x%02X is not 0x%02X\n    %v", msg.chk, msg.calculatedChecksum, *msg))
	}

	return
}

func MWVMessageConinuousScan(byteChannel <-chan byte, messageChannel chan<- MWVMessage) {
	buf := make([]byte, 64)
	var bytesRead int
	var msg MWVMessage
	for {
		buf[0] = 0x00
		for string(buf[0]) != "$" {
			buf[0] = <-byteChannel
		}
		for bytesRead = 1; buf[bytesRead-1] != 0x0D; bytesRead++ {
			buf[bytesRead] = <-byteChannel
		}
		msg.parse(string(buf[:bytesRead-1]))
		if flagVerbose {
			log.Printf("Sending parsed message %v to message channel", msg)
		}
		messageChannel <- msg
		if flagVerbose {
			log.Printf("Message %v was received", msg)
		}
	}
}
