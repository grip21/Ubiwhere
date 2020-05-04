package main

import (
	"fmt"
	"testing"
)

func TestGetCpuMem(t *testing.T) {
	var output, output1 float64
	limitMin := 0.0
	limitMax := 100.0
	output, output1 = getCpuMem()
	if output < limitMin || output > limitMax {
		t.Errorf("incorrect output for the operation: expected bigger than `%f` and lower than `%f`", limitMin, limitMax)
	} else if output1 < limitMin || output1 > limitMax {
		t.Errorf("incorrect output for the operation: expected bigger than `%f` and lower than `%f`", limitMin, limitMax)
	}
}

func TestType1(t *testing.T) {
	query1 := fmt.Sprintf("SELECT voltage, ac, luminosity, wind FROM sensors ORDER BY times DESC LIMIT 4")
	query2 := fmt.Sprintf("SELECT cpu, ram FROM cpuram ORDER BY timest DESC LIMIT 4")
	result1, result := type1(query1, query2)

	if len(result) != 4 || len(result1) != 4 {
		t.Error("Number of returned variables is incorrect")
	} else {
		t.Error("CORRETO")
	}
}

func TestType2CR2S4(t *testing.T) {
	query1 := fmt.Sprintf("SELECT cpu,ram FROM cpuram ORDER BY timest DESC LIMIT 2")
	query2 := fmt.Sprintf("SELECT voltage, ac, luminosity, wind FROM sensors ORDER BY times DESC LIMIT 2")
	query3 := fmt.Sprintf("SELECT cpu,ram FROM cpuram")
	query4 := fmt.Sprintf("SELECT voltage, ac, luminosity, wind FROM sensors")
	var inputType int
	result1, result2 := type2CR2S4(query1, query2, "ram", "cpu", "voltage", "ac", "luminosity", "wind", inputType)
	result3, result4 := type2CR2S4(query3, query4, "ram", "cpu", "voltage", "ac", "luminosity", "wind", inputType)

	if inputType == 2 {
		if len(result1) != 2 && len(result2) != 2 {
			t.Error("ERROR")
		}
	}

	if inputType == 3 {
		if len(result3) != 2 && len(result4) != 4 {
			t.Error("ERROR")
		}
	}
}
