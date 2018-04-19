package main

import (
	"fmt"
	"testing"
)

type Expectations map[interface{}]ExpectedResult

type ExpectedResult struct {
	err   bool
	value interface{}
}

const (
	Error = true
	OK    = false
)

func TestWindAngleParse(t *testing.T) {
	var subject WindAngle

	valuesAndExpectations := Expectations{
		"0":   ExpectedResult{OK, WindAngle(0)},
		"360": ExpectedResult{OK, WindAngle(0)},
		"342": ExpectedResult{OK, WindAngle(342)}}
	for value, result := range valuesAndExpectations {
		err := subject.parse(value.(string))
		if err != nil && result.err {
			t.Error(fmt.Sprintf("Expected different error result after parsing '%v'. Expected error? %t Got error: %s", value, result.err, err))
		}
	}
}

func TestWindAngleReferenceParse(t *testing.T) {
	var subject WindAngleReference

	valuesAndExpectations := Expectations{
		"A": ExpectedResult{OK, WindAngleReference("A")},
		"B": ExpectedResult{OK, WindAngleReference("B")}}
	for value, result := range valuesAndExpectations {
		err := subject.parse(value.(string))
		if (err != nil) != result.err {
			t.Error(fmt.Sprintf("Expected different error result after parsing '%v'. Expected error? %t Got error: %s", value, result.err, err))
		}
	}
}

func TestWindSpeedParse(t *testing.T) {
	var subject WindSpeed

	valuesAndExpectations := Expectations{
		"0.0": ExpectedResult{OK, WindSpeed(0.0)},
		"5.0": ExpectedResult{OK, WindSpeed(5.0)},
		"1.2": ExpectedResult{OK, WindSpeed(1.2)}}
	for value, result := range valuesAndExpectations {
		err := subject.parse(value.(string))
		if (err != nil) != result.err {
			t.Error(fmt.Sprintf("Expected different error result after parsing '%v'. Expected error? %t Got error: %s", value, result.err, err))
		}
	}
}

func TestWindSpeedUnitParse(t *testing.T) {
	var subject WindSpeedUnit

	valuesAndExpectations := Expectations{
		"A": ExpectedResult{OK, WindSpeedUnit("A")},
		"B": ExpectedResult{OK, WindSpeedUnit("B")}}
	for value, result := range valuesAndExpectations {
		err := subject.parse(value.(string))
		if (err != nil) != result.err {
			t.Error(fmt.Sprintf("Expected different error result after parsing '%v'. Expected error? %t Got error: %s", value, result.err, err))
		}
	}
}
