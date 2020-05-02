package main

import (
	"fmt"
	"testing"
)

func TestType1(t *testing.T) {
	query1 := fmt.Sprintf("SELECT voltage, ac, luminosity, wind FROM sensors ORDER BY times DESC LIMIT 4")
	query2 := fmt.Sprintf("SELECT cpu, ram FROM cpuram ORDER BY timest DESC LIMIT 4")
	result, result1 := type1(query1, query2)

	if len(result) != 4 && len(result1) != 4 {
		t.Error("Number of returned variables is incorrect")
	} else {
		t.Error("CORRETO")
	}
}

func TestgetCpuMem(t *testing.T) {
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
