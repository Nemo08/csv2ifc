package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/integrii/flaggy"
	"github.com/jszwec/csvutil"

	itl "csv2ifc/internal"
)

var (
	inputCsvFile         = ""
	outputIFCFile        = "out.ifc"
	version              = "version unknown"
	gitCommit, gitBranch string
)

type CsvRecord struct {
	X     string `csv:"x"`
	Y     string `csv:"y"`
	Z     string `csv:"z"`
	Name  string `csv:"name,omitempty"`
	IType string `csv:"type,omitempty"`
	Descr string `csv:"description,omitempty"`
	Tag   string `csv:"tag,omitempty"`
}

func main() {
	// Add a flag
	flaggy.String(&inputCsvFile, "c", "csv", "Input csv file")
	flaggy.String(&outputIFCFile, "o", "out", "Output ifc file")

	// Parse the flag
	flaggy.Parse()

	//check csv file is exist
	if _, err := os.Stat(inputCsvFile); err != nil {
		fmt.Printf("Csv file does not exist\n")
		os.Exit(1)
	}

	//check writeble ifc file
	file, err := os.OpenFile(outputIFCFile, os.O_WRONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Println("Unable to write to ", outputIFCFile)
			fmt.Println(err)
			os.Exit(1)
		}

	}
	file.Close()

	//open csv file
	csvFile, err := os.OpenFile(inputCsvFile, os.O_RDONLY, 0666)
	defer csvFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//open ifc file
	ifcFile, err := os.OpenFile(outputIFCFile, os.O_CREATE|os.O_TRUNC, 0666)
	defer ifcFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//write ifc file header
	_, err = ifcFile.Write([]byte(itl.IfcHeader))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	//read csv line and write to ifc file
	csvReader := csv.NewReader(csvFile)

	var csvHeader []string //empty csvHeader for csvutils -> csvutils use first line as header
	var count int32
	var b []byte
	var recordStruct CsvRecord
	var shapes []string
	count = 100

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	dec, err := csvutil.NewDecoder(csvReader, csvHeader...)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = dec.Decode(&recordStruct)
	if err != nil {
		log.Fatal("Csv decoder error: ", err)
		os.Exit(1)
	}

	for {
		if err := dec.Decode(&recordStruct); err == io.EOF {
			break
		}
		shapes = append(shapes, "#"+fmt.Sprint(count+6))
		b, count = itl.OneRecord(count, recordStruct.X, recordStruct.Y, recordStruct.Z, recordStruct.Name, recordStruct.IType, recordStruct.Descr, recordStruct.Tag)
		_, err = ifcFile.Write(b)
	}

	//write ifc relation
	_, err = ifcFile.Write([]byte("#" + fmt.Sprint(count) + "= IFCRELCONTAINEDINSPATIALSTRUCTURE('3krIADZCTB3AGKcs8$14SX',$,$,$,(" + strings.Join(shapes, ",") + "),#17);"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	//write ifc file bottom
	_, err = ifcFile.Write([]byte(itl.IfcBottom))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
