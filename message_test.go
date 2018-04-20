package main

import "testing"

type Expectations map[interface{}]ExpectedResult

type ExpectedResult struct {
	err   bool
	value interface{}
}

const (
	Error = true
	OK    = false
)

const (
	AnyValue = -999999
)

func TestWindAngleParse(t *testing.T) {
	var subject WindAngle

	valuesAndExpectations := Expectations{
		"0":   ExpectedResult{OK, WindAngle(0)},
		"360": ExpectedResult{OK, WindAngle(0)},
		"342": ExpectedResult{OK, WindAngle(342)},
		"-1":  ExpectedResult{OK, WindAngle(359)},
		"":    ExpectedResult{Error, InvalidWindAngle},
		"s":   ExpectedResult{Error, InvalidWindAngle}}
	for value, result := range valuesAndExpectations {
		err := subject.parse(value.(string))
		if (err != nil) != result.err {
			if err != nil {
				t.Errorf("parse('%v') should succeed but gave error: %s", value, err)
			} else {
				t.Errorf("parse('%v') should not have succeeded. Got '%v'", value, subject)
			}
		}
		if result.value != subject && (result.value != AnyValue) {
			t.Errorf("Wrong value from parse('%v'):\n    '%v' is not the expected\n    '%v'", value, subject, result.value)
		}
	}
}

func TestWindAngleReferenceParse(t *testing.T) {
	var subject WindAngleReference

	valuesAndExpectations := Expectations{
		"A": ExpectedResult{Error, InvalidWindAngleReference},
		"B": ExpectedResult{Error, InvalidWindAngleReference},
		"R": ExpectedResult{OK, Relative},
		"T": ExpectedResult{OK, True}}
	for value, result := range valuesAndExpectations {
		err := subject.parse(value.(string))
		if (err != nil) != result.err {
			if err != nil {
				t.Errorf("parse('%v') should succeed but gave error: %s", value, err)
			} else {
				t.Errorf("parse('%v') should not have succeeded. Got '%v'", value, subject)
			}
		}
		if result.value != subject && (result.value != AnyValue) {
			t.Errorf("Wrong value from parse('%v'):\n    '%v' is not the expected\n    '%v'", value, subject, result.value)
		}
	}
}

func TestWindSpeedParse(t *testing.T) {
	var subject WindSpeed

	valuesAndExpectations := Expectations{
		"0.0":  ExpectedResult{OK, WindSpeed(0.0)},
		"-3.2": ExpectedResult{Error, InvalidWindSpeed},
		"5.0":  ExpectedResult{OK, WindSpeed(5.0)},
		"1.2":  ExpectedResult{OK, WindSpeed(1.2)}}
	for value, result := range valuesAndExpectations {
		err := subject.parse(value.(string))
		if (err != nil) != result.err {
			if err != nil {
				t.Errorf("parse('%v') should succeed but gave error: %s", value, err)
			} else {
				t.Errorf("parse('%v') should not have succeeded. Got '%v'", value, subject)
			}
		} else if result.value != subject && (result.value != AnyValue) {
			t.Errorf("Wrong value from parse('%v'):\n    '%v' is not the expected\n    '%v'", value, subject, result.value)
		}
	}
}

func TestWindSpeedUnitParse(t *testing.T) {
	var subject WindSpeedUnit

	valuesAndExpectations := Expectations{
		"A": ExpectedResult{Error, InvalidWindSpeedUnit},
		"B": ExpectedResult{Error, InvalidWindSpeedUnit},
		"K": ExpectedResult{OK, KilometersPerHour},
		"M": ExpectedResult{OK, MetersPerSecond},
		"N": ExpectedResult{OK, Knots}}
	for value, result := range valuesAndExpectations {
		err := subject.parse(value.(string))
		if (err != nil) != result.err {
			if err != nil {
				t.Errorf("parse('%v') should succeed but gave error: %s", value, err)
			} else {
				t.Errorf("parse('%v') should not have succeeded. Got '%v'", value, subject)
			}
		}
		if result.value != subject && (result.value != AnyValue) {
			t.Errorf("Wrong value from parse('%v'):\n    '%v' is not the expected\n    '%v'", value, subject, result.value)
		}
	}
}

func TestChecksumParse(t *testing.T) {
	var subject Checksum

	valuesAndExpectations := Expectations{
		"A":  ExpectedResult{Error, InvalidChecksum},
		"":   ExpectedResult{Error, InvalidChecksum},
		"10": ExpectedResult{OK, Checksum(0x10)},
		"0F": ExpectedResult{OK, Checksum(0x0F)}}
	for value, result := range valuesAndExpectations {
		err := subject.parse(value.(string))
		if (err != nil) != result.err {
			if err != nil {
				t.Errorf("parse('%v') should succeed but gave error: %s", value, err)
			} else {
				t.Errorf("parse('%v') should not have succeeded. Got '%v'", value, subject)
			}
		}
		if result.value != subject && (result.value != AnyValue) {
			t.Errorf("Wrong value from parse('%v'):\n    '%v' is not the expected\n    '%v'", value, subject, result.value)
		}
	}
}

func TestMWVMessageParse(t *testing.T) {
	var subject MWVMessage

	valuesAndExpectations := Expectations{
		"$WIMWV,315,R,0.0,M,A*39": ExpectedResult{OK, MWVMessage{
			source:             "$WIMWV,315,R,0.0,M,A*39",
			dir:                WindAngle(315),
			ref:                Relative,
			spd:                WindSpeed(0.0),
			uni:                MetersPerSecond,
			sta:                Valid,
			chk:                0x39,
			calculatedChecksum: 0x39}},
		"$WIMWV,170,R,4.8,M,A*34": ExpectedResult{OK, MWVMessage{
			source:             "$WIMWV,170,R,4.8,M,A*34",
			dir:                WindAngle(170),
			ref:                Relative,
			spd:                WindSpeed(4.8),
			uni:                MetersPerSecond,
			sta:                Valid,
			chk:                0x34,
			calculatedChecksum: 0x34}},
		"$WIMWV,170,R,4.8,M,A*35": ExpectedResult{Error, MWVMessage{
			source:             "$WIMWV,170,R,4.8,M,A*35",
			dir:                WindAngle(170),
			ref:                Relative,
			spd:                WindSpeed(4.8),
			uni:                MetersPerSecond,
			sta:                Valid,
			chk:                0x35,
			calculatedChecksum: 0x34}},
		"B": ExpectedResult{Error, AnyValue}}
	for value, result := range valuesAndExpectations {
		err := subject.parse(value.(string))
		if (err != nil) != result.err {
			if err != nil {
				t.Errorf("parse('%v') should succeed but gave error: %s", value, err)
			} else {
				t.Errorf("parse('%v') should not have succeeded. Got '%v'", value, subject)
			}
		}
		if result.value != subject && (result.value != AnyValue) {
			t.Errorf("Wrong value from parse('%v'):\n    '%v' is not the expected\n    '%v'", value, subject, result.value)
		}
	}

	err := subject.parse("$WIMWV,315,R,0.0,M,A*39")

	if err != nil {
		t.Error(err)
	}
}
