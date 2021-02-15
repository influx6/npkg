package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/influx6/npkg/nerror"
	"github.com/influx6/npkg/nnet"
)

func generateFromCSV(targetFile string, output io.Writer) error {
	var recordFile, recordErr = os.Open(targetFile)
	if recordErr != nil {
		return nerror.WrapOnly(recordErr)
	}

	defer recordFile.Close()

	var reader = csv.NewReader(recordFile)

	var _, printErr = fmt.Fprint(output, "var Records = []nnet.Location{\n")
	if printErr != nil {
		return nerror.WrapOnly(printErr)
	}

	var line = -1
	for {
		var currentLine, curErr = reader.Read()
		if curErr == io.EOF {
			break
		}
		if curErr != nil {
			return nerror.WrapOnly(curErr)
		}

		line++
		fmt.Printf("Handling: %+q on line %d\n", currentLine, line)

		var loc nnet.Location
		loc.FromIPNumeric = currentLine[0]
		loc.ToIPNumeric = currentLine[1]
		loc.CountryCode = currentLine[2]
		loc.CountryName = currentLine[3]
		loc.RegionName = currentLine[4]
		loc.City = currentLine[5]
		loc.Lat = currentLine[6]
		loc.Long = currentLine[7]
		loc.Zipcode = currentLine[8]
		loc.Timezone = currentLine[9]

		if len(loc.FromIPNumeric) != 0 {
			var fromIP, fromIPErr = nnet.IPLongNotation2IPFromString(loc.FromIPNumeric)
			if fromIPErr != nil {
				return nerror.WrapOnly(fromIPErr)
			}
			loc.FromIP = fromIP.String()
		}

		if len(loc.ToIPNumeric) != 0 {
			var toIP, toIPErr = nnet.IPLongNotation2IPFromString(loc.ToIPNumeric)
			if toIPErr != nil {
				return nerror.WrapOnly(toIPErr)
			}
			loc.ToIP = toIP.String()
		}

		if writeErr := writeFile(output, loc); writeErr != nil {
			return nerror.WrapOnly(writeErr)
		}
	}

	var _, err = fmt.Fprint(output, "\n}")
	if err != nil {
		return nerror.WrapOnly(err)
	}
	return nil
}

func writeFile(output io.Writer, value nnet.Location) error {
	var _, err = fmt.Fprintf(output, `{IP:%q,Street:%q,City:%q,State:%q,Postal:%q,CountryCode:q,CountryName:%q,RegionCode:%q,RegionName:%q,Zipcode:%q,Lat:%q,Long:%q,MetroCode:%q,Timezone:%q,AreaCode:%q,FromIP:%q,ToIP:%q,FromIPNumeric:%q,ToIPNumeric:%q},`,
		value.IP,
		value.Street, value.City, value.State,
		value.Postal, value.CountryCode, value.CountryName,
		value.RegionCode, value.RegionName, value.Zipcode,
		value.Lat, value.Long, value.MetroCode, value.Timezone,
		value.AreaCode, value.FromIP, value.ToIP,
		value.FromIPNumeric, value.ToIPNumeric,
	)
	if err != nil {
		return nerror.WrapOnly(err)
	}
	return nil
}

func main() {
	var targetCSVFile = flag.String("ip_file", "", "csv file from IP2Location.com")
	var targetGoFile = flag.String("out", "", "go file")
	var targetPackageName = flag.String("package", "", "name of package to use")

	flag.Parse()

	if len(*targetCSVFile) == 0 || len(*targetGoFile) == 0 || len(*targetPackageName) == 0 {
		fmt.Println(`Generates a go file containing a list of Location objects which are mapped from a IPLocation.com
csv files. This allows us embed this as usable list to find a suited location if any for a target
IP.

Please always provide:

- ip_file: csv file from IP2Location.com
- out: the path to the golang file to generate
- package: the name of the package to use.
`)
		return
	}

	var targetGoFileWriter, targetGoFileErr = os.Create(*targetGoFile)
	if targetGoFileErr != nil {
		log.Fatalf("Error occurred: %+s\n", targetGoFileErr)
		return
	}

	defer targetGoFileWriter.Close()

	var _, printErr = fmt.Fprintf(targetGoFileWriter, "package %s\n", *targetPackageName)
	if printErr != nil {
		log.Fatalf("Error occurred: %+s\n", printErr)
		return
	}

	if genErr := generateFromCSV(*targetCSVFile, targetGoFileWriter); genErr != nil {
		log.Fatalf("Error occurred: %+s\n", genErr)
		return
	}

	log.Println("Finished generating ip files")
}
