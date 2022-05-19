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

const (
	defaultNameForEmptyPset = "Common"
	defaultNameForEmptyType = "IFCLABEL"
)

var (
	inputCsvFile                  = ""
	outputIFCFile                 = "out.ifc"
	emptyPset2Common              bool
	psetPresent                   bool
	psetUseTypes                  bool
	version, gitCommit, gitBranch string
	psetHeader                    map[string]map[string]itl.PsetValue
	onePset                       map[string]string
)

type CsvRecord struct {
	X         string                       `csv:"x"`
	Y         string                       `csv:"y"`
	Z         string                       `csv:"z"`
	Name      string                       `csv:"name,omitempty"`
	IType     string                       `csv:"type,omitempty"`
	Descr     string                       `csv:"description,omitempty"`
	Tag       string                       `csv:"tag,omitempty"`
	OtherData map[string]map[string]string `csv:"-,omitempty"`
}

func main() {
	// Add a flag
	flaggy.String(&inputCsvFile, "c", "csv", "Input csv file")
	flaggy.String(&outputIFCFile, "o", "out", "Output ifc file")
	flaggy.Bool(&psetPresent, "p", "pset", `If flag setted then second line of CSV file interpret as Pset name, property in first line
	           except required and optional fields (x,y,z,name,type,description,tag)`)
	flaggy.Bool(&emptyPset2Common, "e", "empty", `Work with setted pset flag, create pset 'Pset_Common' for all not empty fields in header,
	           except required and optional fields (x,y,z,name,type,description,tag)`)
	flaggy.Bool(&psetUseTypes, "t", "types", `Work with setted pset flag, get type of property from third line`)
	flaggy.SetVersion(version)

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
	var csvHeader []string //empty csvHeader for csvutils -> csvutils use first line as header
	var count int32
	var b, bPset []byte
	var recordStruct CsvRecord
	var shapes []string

	csvReader := csv.NewReader(csvFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	dec, err := csvutil.NewDecoder(csvReader, csvHeader...)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	//read second string if present -p flag
	if psetPresent {
		psetHeader = make(map[string]map[string]itl.PsetValue)
		err = dec.Decode(&recordStruct)
		if err != nil {
			log.Fatal("Csv decoder error: ", err)
			os.Exit(1)
		}
		var psetNamesLineRecords, typesLineRecords []string
		header := dec.Header()
		psetNamesLineRecords = dec.Record()

		//read third line with types if -t flag present for ValueTypeDefinition
		if psetUseTypes {
			err = dec.Decode(&recordStruct)
			if err != nil {
				log.Fatal("Csv decoder error: ", err)
				os.Exit(1)
			}
			typesLineRecords = dec.Record()
		}

		for _, i := range dec.Unused() {
			if header[i] != "" {
				if (psetNamesLineRecords[i] == "") && (emptyPset2Common) {
					psetNamesLineRecords[i] = defaultNameForEmptyPset
				}
				if psetNamesLineRecords[i] == "" {
					break
				}
				if psetHeader[psetNamesLineRecords[i]] == nil {
					psetHeader[psetNamesLineRecords[i]] = make(map[string]itl.PsetValue)
				}
				var pv itl.PsetValue
				pv.Column = i

				if psetUseTypes {
					pv.ValueTypeDefinition = typesLineRecords[i]
				}

				psetHeader[psetNamesLineRecords[i]][header[i]] = pv
			}
		}
	}

	//read csv records until EOF
	count = 100
	for {
		//write point data
		if err := dec.Decode(&recordStruct); err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("Csv decoder error: ", err)
			os.Exit(1)
		}

		shapes = append(shapes, "#"+fmt.Sprint(count+6))
		b, count = itl.OneRecord(count, recordStruct.X, recordStruct.Y, recordStruct.Z, recordStruct.Name, recordStruct.IType, recordStruct.Descr, recordStruct.Tag)
		_, err = ifcFile.Write(b)

		//write Pset data
		bproxy := count - 1
		var pt itl.PsetValue

		for k, v := range psetHeader {
			onePsetData := make(map[string]itl.PsetValue)
			for pk, pv := range v {
				pt.Value = dec.Record()[pv.Column]
				if pv.ValueTypeDefinition != "" {
					pt.ValueTypeDefinition = pv.ValueTypeDefinition
				} else {
					pt.ValueTypeDefinition = defaultNameForEmptyType
				}
				onePsetData[pk] = pt
			}

			bPset, count = itl.OnePset(count, bproxy, k, onePsetData)
			_, err = ifcFile.Write(bPset)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		}

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
