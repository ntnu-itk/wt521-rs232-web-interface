package main

import "fmt"

type WindAngle int
type WindAngleReference string

const (
	Relative WindAngleReference = "R"
	True     WindAngleReference = "T"
)

type WindSpeed float64
type WindSpeedUnit string

const (
	KilometersPerHour WindSpeedUnit = "K"
	MetersPerSecond   WindSpeedUnit = "M"
	Knots             WindSpeedUnit = "N"
)

type MessageValidity string

const (
	Valid   MessageValidity = "A"
	Invalid MessageValidity = "V"
)

type Checksum string

type MWVMessage struct {
	dir WindAngle
	ref WindAngleReference
	spd WindSpeed
	uni WindSpeedUnit
	sta MessageValidity
	chk Checksum
}

func (wa *WindAngle) parse(str string) error {
	n, err := fmt.Sscanf(str, "%d", wa)

	if n != 1 {
		return NewError(fmt.Sprintf("Could not parse wind angle from string '%s'", str))
	}

	return err
}

func (war *WindAngleReference) parse(str string) error {
	n, err := fmt.Sscanf(str, "%1s", war)

	if n != 1 {
		return NewError(fmt.Sprintf("Could not parse wind angle reference from string '%s'", str))
	}

	return err
}

func (ws *WindSpeed) parse(str string) error {
	n, err := fmt.Sscanf(str, "%f", ws)

	if n != 1 {
		return NewError(fmt.Sprintf("Could not parse wind speed from string '%s'", str))
	}

	return err
}

func (wsu *WindSpeedUnit) parse(str string) error {
	n, err := fmt.Sscanf(str, "%1s", wsu)

	if n != 1 {
		return NewError(fmt.Sprintf("Could not parse wind speed unit from string '%s'", str))
	}

	return err
}
