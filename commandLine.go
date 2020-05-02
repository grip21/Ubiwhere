package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func dbConnec() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/ubiwhere")
	if err != nil {
		panic(err.Error())
	}
	return db
}

////////****** Obtaining the cpu and ram usage*******////////////////
func getCpuMem() (float64, float64) {
	var sumCPU, sumMEM float64
	cmd := exec.Command("ps", "aux") // ps aux is the command used to get cpu and ram usage
	var out bytes.Buffer
	cmd.Stdout = &out //catching the command output
	err := cmd.Run()  //running the command
	if err != nil {
		log.Fatal(err)
	}
	for {
		line, err := out.ReadString('\n') //breaking the output in lines
		if err != nil {
			break
		}
		tokens := strings.Split(line, " ") //spliting each output line
		ft := make([]string, 0)
		for _, t := range tokens {
			if t != "" && t != "\t" {
				ft = append(ft, t) //for each line there is a buffer (ft) that keeps all the parameters
			}
		}
		if cpu, err := strconv.ParseFloat(ft[2], 32); err == nil { // parsing the cpu variable, as string, to float
			sumCPU += cpu //all the cpu's used capacity at the instant
		}
		if mem, err := strconv.ParseFloat(ft[3], 32); err == nil { // parsing the ram variable, as string, to float
			sumMEM += mem //all the ram's used capacity at the instant
		}
	}
	log.Println("Used CPU", sumCPU/8, "%", "  Used Memory RAM", sumMEM, "%")
	return sumCPU / 8, sumMEM //the cpu's total used capacity is splitted by 8 because its the total number of my PC's cores
	//otherwise, we would see outputs bigger than 100%
}

////////****** Obtaining 4 random variable *******////////////////
func randomFloat(min, max float32) (v float32, a float32, l float32, w float32) {
	v = min + rand.Float32()*(max-min) //generating the random variable 'voltage'
	a = min + rand.Float32()*(max-min) //generating the random variable 'current'
	l = min + rand.Float32()*(max-min) //generating the random variable 'luminosity'
	w = min + rand.Float32()*(max-min) //generating the random variable 'wind speed'
	return
}

// TYPE 1, get all columns
func type1(query1, query2 string) ([]string, []string) {
	db := dbConnec()
	var out, out1 []string
	rows, err := db.Query(query1) //FROM sensors
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var voltage, ac, luminosity, wind string
		err = rows.Scan(&voltage, &ac, &luminosity, &wind)
		if err != nil {
			panic(err.Error())
		}
		out = append(out, (fmt.Sprintf("Voltage: %s  Current: %s   Luminosity: %s   Wind Speed: %s\n", voltage, ac, luminosity, wind)))
		//log.Println("Voltage:", voltage, " Current:", ac, "  Luminosity:", luminosity, "  Wind Speed:", wind)
	}
	rows1, err := db.Query(query2) //FROM cpuram
	if err != nil {
		panic(err.Error())
	}
	defer rows1.Close()
	for rows1.Next() {
		var cpu, ram float32
		err = rows1.Scan(&cpu, &ram)
		if err != nil {
			panic(err.Error())
		}
		out1 = append(out1, fmt.Sprintf("CPU: %f  RAM: %f\n", cpu, ram))
		//log.Println("CPU:", cpu, " RAM:", ram)
	}
	return out, out1
}

// 1 Column from cpuram Table
func type2CpuRam1(query, split1 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query)
	defer rows.Close()
	var sumCPU, f float64
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		f++
		var x string
		err := rows.Scan(&x)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			out := fmt.Sprintf("%s: %s", split1, x)
			log.Println(out)
		}
		if inputType == 3 {
			if cpuram, err := strconv.ParseFloat(x, 32); err == nil {
				sumCPU += cpuram
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ": ", sumCPU/f)
	}
}

// Ram and Cpu columns from cpuram Table
func type2CpuRam2(query, split1, split2 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query)
	defer rows.Close()
	var sumCPU, sumRAM, f float64
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		f++
		var x, y string
		err := rows.Scan(&x, &y)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x, "  ", split2, ":", y)
		}
		if inputType == 3 {
			if cpuram, err := strconv.ParseFloat(x, 32); err == nil {
				sumCPU += cpuram
			}
			if cpuram2, err := strconv.ParseFloat(y, 32); err == nil {
				sumRAM += cpuram2
			}
		}
	}
	if inputType == 3 {
		log.Println("Average ", split1, ":", sumCPU/f, " Average", split2, ":", sumRAM/f)
	}
}

// Table sensors, 1 column
func type2sensors1(query, split1 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query)
	var sumSensors, f float64
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		f++
		var x string
		err := rows.Scan(&x)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x)
		}
		if inputType == 3 {
			if sumS1, err := strconv.ParseFloat(x, 32); err == nil {
				sumSensors += sumS1
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ": ", sumSensors/f)
	}
}

// Table sensors,2 columns
func type2sensors2(query, split1, split2 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query)
	defer rows.Close()
	var sumSensors1, sumSensors2, f float64
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		f++
		var x, y string
		err := rows.Scan(&x, &y)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x, "  ", split2, ":", y)
		}
		if inputType == 3 {
			if sumS1, err := strconv.ParseFloat(x, 32); err == nil {
				sumSensors1 += sumS1
			}
			if sumS2, err := strconv.ParseFloat(y, 32); err == nil {
				sumSensors2 += sumS2
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ":", sumSensors1/f, "  Average", split2, ":", sumSensors2/f)
	}
}

// Table sensors, 3 columns
func type2sensors3(query, split1, split2, split3 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	var sumSensors1, sumSensors2, sumSensors3, f float64
	for rows.Next() {
		f++
		var x, y, z string
		err := rows.Scan(&x, &y, &z)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x, "  ", split2, ":", y, split3, ":", z)
		}
		if inputType == 3 {
			if sumS1, err := strconv.ParseFloat(x, 32); err == nil {
				sumSensors1 += sumS1
			}
			if sumS2, err := strconv.ParseFloat(y, 32); err == nil {
				sumSensors2 += sumS2
			}
			if sumS3, err := strconv.ParseFloat(z, 32); err == nil {
				sumSensors3 += sumS3
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ":", sumSensors1/f, "  Average", split2, ":", sumSensors2/f, "Average", split3, ":", sumSensors3/f)
	}
}

// Table sensors, all columns
func type2sensors4(query, split1, split2, split3, split4 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	var sumSensors1, sumSensors2, sumSensors3, sumSensors4, f float64
	for rows.Next() {
		f++
		var x, y, z, w string
		err := rows.Scan(&x, &y, &z, &w)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x, "  ", split2, ":", y, split3, ":", z, split4, ":", w)
		}
		if inputType == 3 {
			if sumS1, err := strconv.ParseFloat(x, 32); err == nil {
				sumSensors1 += sumS1
			}
			if sumS2, err := strconv.ParseFloat(y, 32); err == nil {
				sumSensors2 += sumS2
			}
			if sumS3, err := strconv.ParseFloat(z, 32); err == nil {
				sumSensors3 += sumS3
			}
			if sumS4, err := strconv.ParseFloat(w, 32); err == nil {
				sumSensors4 += sumS4
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ":", sumSensors1/f, "  Average", split2, ":", sumSensors2/f, "Average", split3, ":", sumSensors3/f, "Average", split4, ":", sumSensors4/f)
	}
}

// 1 column from each table
func type2CR1S1(query1, query2, split1, split2 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query1)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}
	var sumCR1, sumS1, f, s float64
	for rows.Next() {
		f++
		var x string
		err := rows.Scan(&x)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x)
		}
		if inputType == 3 {
			if sumC1, err := strconv.ParseFloat(x, 32); err == nil {
				sumCR1 += sumC1
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ": ", sumCR1/f)
	}

	rows1, err := db.Query(query2)
	defer rows1.Close()
	if err != nil {
		panic(err.Error())
	}
	for rows1.Next() {
		s++
		var y string
		err := rows1.Scan(&y)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split2, ":", y)
		}
		if inputType == 3 {
			if sumSens1, err := strconv.ParseFloat(y, 32); err == nil {
				sumS1 += sumSens1
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split2, ": ", sumS1/s)
	}
}

//1 column from CPURAM and 2 columns from SENSORS
func type2CR1S2(query1, query2, split1, split2, split3 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query1)
	defer rows.Close()
	var sumCR1, sumS1, sumS2, f, s float64
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		f++
		var x string
		err := rows.Scan(&x)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x)
		}
		if inputType == 3 {
			if sumC1, err := strconv.ParseFloat(x, 32); err == nil {
				sumCR1 += sumC1
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ": ", sumCR1/f)
	}

	rows1, err := db.Query(query2)
	if err != nil {
		panic(err.Error())
	}
	for rows1.Next() {
		s++
		var y, z string
		err := rows1.Scan(&y, &z)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split2, ":", y, split3, ":", z)
		}
		if inputType == 3 {
			if sumSens1, err := strconv.ParseFloat(y, 32); err == nil {
				sumS1 += sumSens1
			}
			if sumSens2, err := strconv.ParseFloat(z, 32); err == nil {
				sumS2 += sumSens2
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split2, ": ", sumS1/s, "  Average", split3, ": ", sumS2/s)
	}
}

//1 column from CPURAM and 3 columns from SENSORS
func type2CR1S3(query1, query2, split1, split2, split3, split4 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query1)
	defer rows.Close()
	var sumCR1, sumS1, sumS2, sumS3, f, s float64
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		f++
		var x string
		err := rows.Scan(&x)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x)
		}
		if inputType == 3 {
			if sumC1, err := strconv.ParseFloat(x, 32); err == nil {
				sumCR1 += sumC1
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ": ", sumCR1/f)
	}
	rows1, err := db.Query(query2)
	defer rows1.Close()
	if err != nil {
		panic(err.Error())
	}
	for rows1.Next() {
		s++
		var y, z, w string
		err := rows1.Scan(&y, &z, &w)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split2, ":", y, " ", split3, ":", z, " ", split4, ":", w)
		}
		if inputType == 3 {
			if sumSens1, err := strconv.ParseFloat(y, 32); err == nil {
				sumS1 += sumSens1
			}
			if sumSens2, err := strconv.ParseFloat(z, 32); err == nil {
				sumS2 += sumSens2
			}
			if sumSens3, err := strconv.ParseFloat(w, 32); err == nil {
				sumS3 += sumSens3
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split2, ": ", sumS1/s, "  Average", split3, ": ", sumS2/s, "  Average", split4, ": ", sumS3/s)
	}
}

//1 column from CPURAM and 4 columns from SENSORS
func type2CR1S4(query1, query2, split1, split2, split3, split4, split5 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query1)
	defer rows.Close()
	var sumCR1, sumS1, sumS2, sumS3, sumS4, f, s float64
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		f++
		var x string
		err := rows.Scan(&x)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x)
		}
		if inputType == 3 {
			if sumC1, err := strconv.ParseFloat(x, 32); err == nil {
				sumCR1 += sumC1
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ": ", sumCR1/f)
	}

	rows1, err := db.Query(query2)
	if err != nil {
		panic(err.Error())
	}
	for rows1.Next() {
		s++
		var y, z, w, k string
		err := rows1.Scan(&y, &z, &w, &k)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split2, ":", y, " ", split3, ":", z, " ", split4, ":", w, " ", split5, ":", k)
		}
		if inputType == 3 {
			if sumSens1, err := strconv.ParseFloat(y, 32); err == nil {
				sumS1 += sumSens1
			}
			if sumSens2, err := strconv.ParseFloat(z, 32); err == nil {
				sumS2 += sumSens2
			}
			if sumSens3, err := strconv.ParseFloat(w, 32); err == nil {
				sumS3 += sumSens3
			}
			if sumSens4, err := strconv.ParseFloat(k, 32); err == nil {
				sumS4 += sumSens4
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split2, ": ", sumS1/s, "  Average", split3, ": ", sumS2/s, "  Average", split4, ": ", sumS3/s, "  Average", split5, ": ", sumS4/s)
	}
}

// 2 columns from CPURAM and 1 column from SENSORS
func type2CR2S1(query1, query2, split1, split2, split3 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query1)
	defer rows.Close()
	var sumCR1, sumCR2, sumS1, f, s float64
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		f++
		var x, y string
		err := rows.Scan(&x, &y)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x, split2, ":", y)
		}
		if inputType == 3 {
			if sumC1, err := strconv.ParseFloat(x, 32); err == nil {
				sumCR1 += sumC1
			}
			if sumC2, err := strconv.ParseFloat(y, 32); err == nil {
				sumCR2 += sumC2
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ": ", sumCR1/f, "Average", split2, ": ", sumCR2/f)
	}

	rows1, err := db.Query(query2)
	if err != nil {
		panic(err.Error())
	}
	for rows1.Next() {
		s++
		var z string
		err := rows1.Scan(&z)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split3, ":", z)
		}
		if inputType == 3 {
			if sumSens1, err := strconv.ParseFloat(z, 32); err == nil {
				sumS1 += sumSens1
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split3, ": ", sumS1/s)
	}
}

//2 columns from each table
func type2CR2S2(query1, query2, split1, split2, split3, split4 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query1)
	defer rows.Close()
	var sumCR1, sumCR2, sumS1, sumS2, f, s float64
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		f++
		var x, y string
		err := rows.Scan(&x, &y)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x, split2, ":", y)
		}
		if inputType == 3 {
			if sumC1, err := strconv.ParseFloat(x, 32); err == nil {
				sumCR1 += sumC1
			}
			if sumC2, err := strconv.ParseFloat(y, 32); err == nil {
				sumCR2 += sumC2
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ": ", sumCR1/f, "Average", split2, ": ", sumCR2/f)
	}

	rows1, err := db.Query(query2)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}
	for rows1.Next() {
		s++
		var z, w string
		err := rows1.Scan(&z, &w)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split2, ":", z, " ", split3, ":", w)
		}
		if inputType == 3 {
			if sumSens1, err := strconv.ParseFloat(z, 32); err == nil {
				sumS1 += sumSens1
			}
			if sumSens2, err := strconv.ParseFloat(w, 32); err == nil {
				sumS2 += sumSens2
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split3, ": ", sumS1/s, "Average", split4, ": ", sumS2/s)
	}
}

//2 columns from  CPURAM and 3 columns from SENSORS
func type2CR2S3(query1, query2, split1, split2, split3, split4, split5 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query1)
	defer rows.Close()
	var sumCR1, sumCR2, sumS1, sumS2, sumS3, f, s float64
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		f++
		var x, y string
		err := rows.Scan(&x, &y)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x, split2, ":", y)
		}
		if inputType == 3 {
			if sumC1, err := strconv.ParseFloat(x, 32); err == nil {
				sumCR1 += sumC1
			}
			if sumC2, err := strconv.ParseFloat(y, 32); err == nil {
				sumCR2 += sumC2
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ": ", sumCR1/f, "Average", split2, ": ", sumCR2/f)
	}

	rows1, err := db.Query(query2)
	defer rows1.Close()
	if err != nil {
		panic(err.Error())
	}
	for rows1.Next() {
		s++
		var z, w, k string
		err := rows1.Scan(&z, &w, &k)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split2, ":", z, " ", split3, ":", w, split4, ":", k)
		}
		if inputType == 3 {
			if sumSens1, err := strconv.ParseFloat(z, 32); err == nil {
				sumS1 += sumSens1
			}
			if sumSens2, err := strconv.ParseFloat(w, 32); err == nil {
				sumS2 += sumSens2
			}
			if sumSens3, err := strconv.ParseFloat(k, 32); err == nil {
				sumS3 += sumSens3
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split3, ": ", sumS1/s, "Average", split4, ": ", sumS2/s, "Average", split5, ": ", sumS3/s)
	}
}

//2 columns from  CPURAM and 4 columns from SENSORS
func type2CR2S4(query1, query2, split1, split2, split3, split4, split5, split6 string, inputType int) {
	db := dbConnec()
	rows, err := db.Query(query1)
	defer rows.Close()
	var sumCR1, sumCR2, sumS1, sumS2, sumS3, sumS4, f, s float64
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		f++
		var x, y string
		err := rows.Scan(&x, &y)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split1, ":", x, split2, ":", y)
		}
		if inputType == 3 {
			if sumC1, err := strconv.ParseFloat(x, 32); err == nil {
				sumCR1 += sumC1
			}
			if sumC2, err := strconv.ParseFloat(y, 32); err == nil {
				sumCR2 += sumC2
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split1, ": ", sumCR1/f, "Average", split2, ": ", sumCR2/f)
	}

	rows1, err := db.Query(query2)
	if err != nil {
		panic(err.Error())
	}
	for rows1.Next() {
		s++
		var z, w, k, r string
		err := rows1.Scan(&z, &w, &k, &r)
		if err != nil {
			panic(err.Error())
		}
		if inputType == 2 {
			log.Println(split3, ":", z, " ", split4, ":", w, split5, ":", k, split6, ":", r)
		}
		if inputType == 3 {
			if sumSens1, err := strconv.ParseFloat(z, 32); err == nil {
				sumS1 += sumSens1
			}
			if sumSens2, err := strconv.ParseFloat(w, 32); err == nil {
				sumS2 += sumSens2
			}
			if sumSens3, err := strconv.ParseFloat(k, 32); err == nil {
				sumS3 += sumSens3
			}
			if sumSens4, err := strconv.ParseFloat(r, 32); err == nil {
				sumS4 += sumSens4
			}
		}
	}
	if inputType == 3 {
		log.Println("Average", split3, ": ", sumS1/s, "Average", split4, ": ", sumS2/s, "Average", split5, ": ", sumS3/s, "Average", split6, ": ", sumS4/s)
	}
}

func main() {
	optionsCpuRam := [2]string{}
	optionsCpuRam[0] = "cpu"
	optionsCpuRam[1] = "ram"
	optionsSensors := [4]string{}
	optionsSensors[0] = "voltage"
	optionsSensors[1] = "ac"
	optionsSensors[2] = "luminosity"
	optionsSensors[3] = "wind"
	var boolean1, boolean2, type2, type3, errorInput bool
	var cpuramCount, sensorsCount int
	var Nrows string
	db := dbConnec()
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to This Challenge")
	fmt.Println("---------------------")
	fmt.Println("Choose one of four input types!")
	log.Println("If you want obtain from your PC, the CPU and RAM usage percentage. Input -> '0' ")
	fmt.Println("If you want the last 'X' (integer number) samples of the variables press 1 followed by the number of samples (Input example)-> 1 'X'")
	fmt.Println("If you want the last 'X' (integer number) samples of the variables press 1 followed by the number " +
		"of samples and the types of available variables (Input example)-> 2 'X' 'cpu' 'ram' 'voltage' 'ac' 'luminosity' 'wind'")
	fmt.Println("If you want to get an average of the value of one or more variables (Input example)-> 3 'cpu' 'ram' 'voltage' 'ac' 'luminosity' 'wind'")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		splitBuf := strings.Split(text, " ")
		if strings.Compare("0", splitBuf[0]) != 0 {
			Nrows = splitBuf[1]
		}
		/******* Type 0 */ ///
		if strings.Compare("0", splitBuf[0]) == 0 && len(splitBuf) == 1 {
			db := dbConnec()
			for tick := range time.Tick(1 * time.Second) { // get CPU and RAM memory usage is executed every second
				cpu, ram := getCpuMem() //executing function to obtain cpu and ram
				timestam := time.Now()
				insert, err := db.Prepare("INSERT INTO cpuram (cpu, ram, timest)VALUES (?,?,?)") //preparing query to do the insert
				if err != nil {
					panic(err.Error())
				}
				insert.Exec(cpu, ram, timestam) //inserting both variables in the cpuram table
				defer insert.Close()

				//Generate 4 random variables to sensors table every second
				db1 := dbConnec()
				rand.Seed(time.Now().UnixNano())
				voltage, ac, luminosity, wind := randomFloat(0.0, 99.9)
				log.Println("Voltage", voltage, "AC", ac, "Luminosity", luminosity, "Wind", wind, "timestamp", timestam)
				insert1, err := db1.Prepare("INSERT INTO sensors (voltage,ac,luminosity,wind,times) VALUES (?,?,?,?,?)")
				if err != nil {
					panic(err.Error())
				}
				insert1.Exec(voltage, ac, luminosity, wind, timestam) //inserting variables in the sensors table
				defer insert1.Close()
				log.Println(tick)
			}
		}
		/***************TYPE1************************/
		if strings.Compare("1", splitBuf[0]) == 0 { // TYPE 1->SELECT last N rows of all variables. ! Query for each table
			query1 := fmt.Sprintf("SELECT voltage, ac, luminosity, wind FROM sensors ORDER BY times DESC LIMIT %s", Nrows)
			query2 := fmt.Sprintf("SELECT cpu, ram FROM cpuram ORDER BY timest DESC LIMIT %s", Nrows)
			log.Println(type1(query1, query2))

		}
		/////////**********TYPE2*******************////////
		if strings.Compare("2", splitBuf[0]) == 0 { // TYPE 2-> SELECT last N rows for X variables
			type2 = true
			for i := 2; i < len(splitBuf); i++ {
				for j := 0; j <= 1; j++ { //check how many vars are required from the cpuram table
					if strings.Compare(splitBuf[i], optionsCpuRam[j]) == 0 {
						cpuramCount++
						boolean1 = true
					}
				}
				for k := 0; k < len(optionsSensors); k++ { //check how many vars are required from the sensors table
					if strings.Compare(splitBuf[i], optionsSensors[k]) == 0 {
						sensorsCount++
						boolean2 = true
					}
				}
			}
			total := cpuramCount + sensorsCount
			if total < len(splitBuf)-2 {
				log.Println("Your input parameters are incorrect!! Try again")
				type2 = false
			}
		}
		//////////////***********TYPE3*********************////////////////////////////////
		if strings.Compare("3", splitBuf[0]) == 0 && !errorInput { // TYPE 3-> SELECT last N rows for X variables
			type3 = true
			for i := 1; i < len(splitBuf); i++ {
				for j := 0; j <= 1; j++ { //check how many vars are required from the cpuram table
					if strings.Compare(splitBuf[i], optionsCpuRam[j]) == 0 {
						cpuramCount++
						boolean1 = true
					}
				}
				for k := 0; k < len(optionsSensors); k++ { //check how many vars are required from the sensors table
					if strings.Compare(splitBuf[i], optionsSensors[k]) == 0 {
						sensorsCount++
						boolean2 = true
					}
				}
			}
			total := cpuramCount + sensorsCount
			if total < len(splitBuf)-1 {
				log.Println("Your input parameters are incorrect!! Try again")
				type3 = false
			}
		}
		errorInput = false
		//FROM now on, all the function calls have the following logic -> (query, queryN...,column1,columnN..., inputType)
		// if boolean1 && !boolean2 -> Search only on CPURAM table
		if boolean1 && !boolean2 {
			if len(splitBuf) == 3 && type2 { //-> only CPU or RAM
				query := fmt.Sprintf("SELECT %s FROM cpuram ORDER BY timest DESC LIMIT %s", splitBuf[2], Nrows)
				type2CpuRam1(query, splitBuf[2], 2)
			}
			if len(splitBuf) == 2 && type3 { //-> only CPU or RAM
				query := fmt.Sprintf("SELECT %s FROM cpuram ", splitBuf[1])
				type2CpuRam1(query, splitBuf[1], 3)
			}

			if len(splitBuf) == 4 && type2 { //-> CPU and RAM
				query := fmt.Sprintf("SELECT %s,%s FROM cpuram ORDER BY timest DESC LIMIT %s", splitBuf[2], splitBuf[3], Nrows)
				type2CpuRam2(query, splitBuf[2], splitBuf[3], 2)
			}
			if len(splitBuf) == 3 && type3 { //-> CPU and RAM
				query := fmt.Sprintf("SELECT %s,%s FROM cpuram", splitBuf[1], splitBuf[2])
				type2CpuRam2(query, splitBuf[1], splitBuf[2], 3)
			}
			type2 = false
			type3 = false
			cpuramCount = 0
			boolean1 = false
		}
		//Search on SENSORS table
		if !boolean1 && boolean2 {
			if len(splitBuf) == 3 && type2 { //-> only VOLTAGE or AC or LUMINOSITY or WIND
				query := fmt.Sprintf("SELECT %s FROM sensors ORDER BY times DESC LIMIT %s", splitBuf[2], Nrows)
				type2sensors1(query, splitBuf[2], 2)
			}
			if len(splitBuf) == 2 && type3 { //-> only VOLTAGE or AC or LUMINOSITY or WIND
				query := fmt.Sprintf("SELECT %s FROM sensors", splitBuf[1])
				type2sensors1(query, splitBuf[1], 3)
			}

			if len(splitBuf) == 4 && type2 { //-> Voltage-AC or Voltage-LUM or Voltage-Wind or AC-LUM or AC-Wind or LUM-Wind
				query := fmt.Sprintf("SELECT %s,%s FROM sensors ORDER BY times DESC LIMIT %s", splitBuf[2], splitBuf[3], Nrows)
				type2sensors2(query, splitBuf[2], splitBuf[3], 2)
			}
			if len(splitBuf) == 3 && type3 { //-> Voltage-AC or Voltage-LUM or Voltage-Wind or AC-LUM or AC-Wind or LUM-Wind
				query := fmt.Sprintf("SELECT %s,%s FROM sensors", splitBuf[1], splitBuf[2])
				type2sensors2(query, splitBuf[1], splitBuf[2], 3)
			}

			if len(splitBuf) == 5 && type2 { //-> Voltage-AC-LUM or Voltage-AC-WIND or Voltage-Wind-LUM or AC-LUM-WIND
				query := fmt.Sprintf("SELECT %s,%s,%s FROM sensors ORDER BY times DESC LIMIT %s", splitBuf[2], splitBuf[3], splitBuf[4], Nrows)
				type2sensors3(query, splitBuf[2], splitBuf[3], splitBuf[4], 2)
			}
			if len(splitBuf) == 4 && type3 { //-> Voltage-AC-LUM or Voltage-AC-WIND or Voltage-Wind-LUM or AC-LUM-WIND
				query := fmt.Sprintf("SELECT %s,%s,%s FROM sensors", splitBuf[1], splitBuf[2], splitBuf[3])
				type2sensors3(query, splitBuf[1], splitBuf[2], splitBuf[3], 3)
			}

			if len(splitBuf) == 6 && type2 { //-> Voltage-AC-LUM-WIND
				query := fmt.Sprintf("SELECT %s,%s,%s,%s FROM sensors ORDER BY times DESC LIMIT %s", splitBuf[2], splitBuf[3], splitBuf[4], splitBuf[5], Nrows)
				type2sensors4(query, splitBuf[2], splitBuf[3], splitBuf[4], splitBuf[5], 2)
			}
			if len(splitBuf) == 5 && type3 { //-> Voltage-AC-LUM-WIND
				query := fmt.Sprintf("SELECT %s,%s,%s,%s FROM sensors", splitBuf[1], splitBuf[2], splitBuf[3], splitBuf[4])
				type2sensors4(query, splitBuf[1], splitBuf[2], splitBuf[3], splitBuf[4], 3)
			}
			type2 = false
			type3 = false
			sensorsCount = 0
			boolean2 = false
		}

		if boolean1 && boolean2 { //search on both tables
			if len(splitBuf) == 4 && cpuramCount == 1 && sensorsCount == 1 && type2 { //1 var from SENSORS 1 var from CPURAM
				query1 := fmt.Sprintf("SELECT %s FROM cpuram ORDER BY timest DESC LIMIT %s", splitBuf[2], Nrows)
				query2 := fmt.Sprintf("SELECT %s FROM sensors ORDER BY times DESC LIMIT %s", splitBuf[3], Nrows)
				type2CR1S1(query1, query2, splitBuf[2], splitBuf[3], 2)
			}
			if len(splitBuf) == 3 && cpuramCount == 1 && sensorsCount == 1 && type3 { //1 var from SENSORS 1 var from CPURAM
				query1 := fmt.Sprintf("SELECT %s FROM cpuram", splitBuf[1])
				query2 := fmt.Sprintf("SELECT %s FROM sensors", splitBuf[2])
				type2CR1S1(query1, query2, splitBuf[1], splitBuf[2], 3)
			}

			if len(splitBuf) == 5 && cpuramCount == 1 && sensorsCount == 2 && type2 { // 1 var from CPURAM 2 vars from SENSORS
				query1 := fmt.Sprintf("SELECT %s FROM cpuram ORDER BY timest DESC LIMIT %s", splitBuf[2], Nrows)
				query2 := fmt.Sprintf("SELECT %s,%s FROM sensors ORDER BY times DESC LIMIT %s", splitBuf[3], splitBuf[4], Nrows)
				type2CR1S2(query1, query2, splitBuf[2], splitBuf[3], splitBuf[4], 2)
			}
			if len(splitBuf) == 4 && cpuramCount == 1 && sensorsCount == 2 && type3 { // 1 var from CPURAM 2 vars from SENSORS
				query1 := fmt.Sprintf("SELECT %s FROM cpuram", splitBuf[1])
				query2 := fmt.Sprintf("SELECT %s,%s FROM sensors", splitBuf[2], splitBuf[3])
				type2CR1S2(query1, query2, splitBuf[1], splitBuf[2], splitBuf[3], 3)
			}

			if len(splitBuf) == 6 && cpuramCount == 1 && sensorsCount == 3 && type2 { // 1 var from CPURAM 3 vars from SENSORS
				query1 := fmt.Sprintf("SELECT %s FROM cpuram ORDER BY timest DESC LIMIT %s", splitBuf[2], Nrows)
				query2 := fmt.Sprintf("SELECT %s,%s,%s FROM sensors ORDER BY times DESC LIMIT %s", splitBuf[3], splitBuf[4], splitBuf[5], Nrows)
				type2CR1S3(query1, query2, splitBuf[2], splitBuf[3], splitBuf[4], splitBuf[5], 2)
			}
			if len(splitBuf) == 5 && cpuramCount == 1 && sensorsCount == 3 && type3 { // 1 var from CPURAM 3 vars from SENSORS
				query1 := fmt.Sprintf("SELECT %s FROM cpuram", splitBuf[1])
				query2 := fmt.Sprintf("SELECT %s,%s,%s FROM sensors", splitBuf[2], splitBuf[3], splitBuf[4])
				type2CR1S3(query1, query2, splitBuf[1], splitBuf[2], splitBuf[3], splitBuf[4], 3)
			}

			if len(splitBuf) == 7 && cpuramCount == 1 && sensorsCount == 4 && type2 { // 1 var from CPURAM 4 vars from SENSORS
				query1 := fmt.Sprintf("SELECT %s FROM cpuram ORDER BY timest DESC LIMIT %s", splitBuf[2], Nrows)
				query2 := fmt.Sprintf("SELECT %s,%s,%s,%s FROM sensors ORDER BY times DESC LIMIT %s", splitBuf[3], splitBuf[4], splitBuf[5], splitBuf[6], Nrows)
				type2CR1S4(query1, query2, splitBuf[2], splitBuf[3], splitBuf[4], splitBuf[5], splitBuf[6], 2)
			}
			if len(splitBuf) == 6 && cpuramCount == 1 && sensorsCount == 4 && type3 { // 1 var from CPURAM 4 vars from SENSORS
				query1 := fmt.Sprintf("SELECT %s FROM cpuram", splitBuf[1])
				query2 := fmt.Sprintf("SELECT %s,%s,%s,%s FROM sensors", splitBuf[2], splitBuf[3], splitBuf[4], splitBuf[5])
				type2CR1S4(query1, query2, splitBuf[1], splitBuf[2], splitBuf[3], splitBuf[4], splitBuf[5], 3)
			}

			if len(splitBuf) == 5 && cpuramCount == 2 && sensorsCount == 1 && type2 { // 2 vars from CPURAM 1 var from SENSORS
				query1 := fmt.Sprintf("SELECT %s,%s FROM cpuram ORDER BY timest DESC LIMIT %s", splitBuf[2], splitBuf[3], Nrows)
				query2 := fmt.Sprintf("SELECT %s FROM sensors ORDER BY times DESC LIMIT %s", splitBuf[4], Nrows)
				type2CR2S1(query1, query2, splitBuf[2], splitBuf[3], splitBuf[4], 2)
			}
			if len(splitBuf) == 4 && cpuramCount == 2 && sensorsCount == 1 && type3 { // 2 vars from CPURAM 1 var from SENSORS
				query1 := fmt.Sprintf("SELECT %s,%s FROM cpuram", splitBuf[1], splitBuf[2])
				query2 := fmt.Sprintf("SELECT %s FROM sensors", splitBuf[3])
				type2CR2S1(query1, query2, splitBuf[1], splitBuf[2], splitBuf[3], 3)
			}

			if len(splitBuf) == 6 && cpuramCount == 2 && sensorsCount == 2 && type2 { //2 vars from each table
				query1 := fmt.Sprintf("SELECT %s,%s FROM cpuram ORDER BY timest DESC LIMIT %s", splitBuf[2], splitBuf[3], Nrows)
				query2 := fmt.Sprintf("SELECT %s,%s FROM sensors ORDER BY times DESC LIMIT %s", splitBuf[4], splitBuf[5], Nrows)
				type2CR2S2(query1, query2, splitBuf[2], splitBuf[3], splitBuf[4], splitBuf[5], 2)
			}
			if len(splitBuf) == 5 && cpuramCount == 2 && sensorsCount == 2 && type3 { //2 vars from each table
				query1 := fmt.Sprintf("SELECT %s,%s FROM cpuram", splitBuf[1], splitBuf[2])
				query2 := fmt.Sprintf("SELECT %s,%s FROM sensors", splitBuf[3], splitBuf[4])
				type2CR2S2(query1, query2, splitBuf[1], splitBuf[2], splitBuf[3], splitBuf[4], 3)
			}

			if len(splitBuf) == 7 && cpuramCount == 2 && sensorsCount == 3 && type2 { //2 vars from CPURAM 3 vars from SENSORS
				query1 := fmt.Sprintf("SELECT %s,%s FROM cpuram ORDER BY timest DESC LIMIT %s", splitBuf[2], splitBuf[3], Nrows)
				query2 := fmt.Sprintf("SELECT %s,%s,%s FROM sensors ORDER BY times DESC LIMIT %s", splitBuf[4], splitBuf[5], splitBuf[6], Nrows)
				type2CR2S3(query1, query2, splitBuf[2], splitBuf[3], splitBuf[4], splitBuf[5], splitBuf[6], 2)
			}
			if len(splitBuf) == 6 && cpuramCount == 2 && sensorsCount == 3 && type3 { //2 vars from CPURAM 3 vars from SENSORS
				query1 := fmt.Sprintf("SELECT %s,%s FROM cpuram", splitBuf[1], splitBuf[2])
				query2 := fmt.Sprintf("SELECT %s,%s,%s FROM sensors", splitBuf[3], splitBuf[4], splitBuf[5])
				type2CR2S3(query1, query2, splitBuf[1], splitBuf[2], splitBuf[3], splitBuf[4], splitBuf[5], 3)
			}

			if len(splitBuf) == 8 && cpuramCount == 2 && sensorsCount == 4 && type2 {
				query1 := fmt.Sprintf("SELECT %s,%s FROM cpuram ORDER BY timest DESC LIMIT %s", splitBuf[2], splitBuf[3], Nrows)
				query2 := fmt.Sprintf("SELECT %s,%s,%s,%s FROM sensors ORDER BY times DESC LIMIT %s", splitBuf[4], splitBuf[5], splitBuf[6], splitBuf[7], Nrows)
				type2CR2S4(query1, query2, splitBuf[2], splitBuf[3], splitBuf[4], splitBuf[5], splitBuf[6], splitBuf[7], 2)
			}
			if len(splitBuf) == 7 && cpuramCount == 2 && sensorsCount == 4 && type3 {
				query1 := fmt.Sprintf("SELECT %s,%s FROM cpuram", splitBuf[1], splitBuf[2])
				query2 := fmt.Sprintf("SELECT %s,%s,%s,%s FROM sensors", splitBuf[3], splitBuf[4], splitBuf[5], splitBuf[6])
				type2CR2S4(query1, query2, splitBuf[1], splitBuf[2], splitBuf[3], splitBuf[4], splitBuf[5], splitBuf[6], 3)
			}

			sensorsCount = 0
			cpuramCount = 0
			boolean1 = false
			boolean2 = false
		}
	}
}
