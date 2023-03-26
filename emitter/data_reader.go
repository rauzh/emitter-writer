package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"gitlab.com/Skinass/hakaton-2023-1-1/common/storage"
)

type DataReader struct {
	dataChan chan storage.Message
}

func (r *DataReader) parseAmazonData(dataStr string) (storage.Message, error) {
	var data storage.Message

	idRegex := `Id:\s*(?P<id>\d+)\s*`
	asinRegex := `ASIN:\s*(?P<asin>\w+)?\s*`
	titleRegex := `title:\s*(?P<title>.+)?\s*`
	groupRegex := `group:\s*(?P<group>.+)?\s*`
	salesrankRegex := `salesrank:\s*(?P<salesrank>\d+)?\s*`
	// similarRegex := `similar:\s*(?P<similar_num>\d+)\s+(?P<similars>(\w+\s+)+)`
	regex, err := regexp.Compile(idRegex + asinRegex + `((?P<discontinued>discontinued product)|` + titleRegex + groupRegex + salesrankRegex + ")")
	if err != nil {
		return data, fmt.Errorf("Emitter parser: regex compilation error")
	}

	if !regex.MatchString(dataStr) {
		return data, fmt.Errorf("Emitter parser: parse error")
	}

	groupNames := regex.SubexpNames()
	for _, match := range regex.FindAllStringSubmatch(dataStr, -1) {
		for groupIdx, groupValue := range match {
			name := groupNames[groupIdx]
			if name != "" {
				switch name {
				case "discontinued":
					if groupValue == "discontinued product" {
						data.IsDiscontinued = true
					}
				case "id":
					idInt, _ := strconv.ParseUint(groupValue, 10, 32)
					data.ID = uint(idInt)
				case "asin":
					data.ASIN = groupValue
				case "title":
					data.Title = groupValue
				case "group":
					data.Group = groupValue
				case "salesrank":
					data.Salesrank, _ = strconv.ParseUint(groupValue, 10, 32)
				}
			}
		}
	}

	return data, nil
}

func (r *DataReader) readAmazonFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var amazonDataLines string
	isMetainfoSkipped := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if !isMetainfoSkipped {
				isMetainfoSkipped = true
			} else {
				data, err := r.parseAmazonData(amazonDataLines)
				if err != nil {
					close(r.dataChan)
					return err
				}
				r.dataChan <- data
			}
			amazonDataLines = ""
		} else {
			amazonDataLines += line
		}
	}

	close(r.dataChan)

	return nil
}

func (r *DataReader) NextMessage() (storage.Message, bool) {
	message, ok := <-r.dataChan
	if ok {
		return message, true
	}
	return message, false
}

func NewDataReader(filename string) *DataReader {
	var reader DataReader
	reader.dataChan = make(chan storage.Message)

	go func() {
		err := reader.readAmazonFile("amazon-meta.txt")
		fmt.Println(err)
	}()

	return &reader
}

/*

// EXAMPLE:

func main() {
	reader := NewDataReader("amazon-meta.txt")
	message, ok := reader.NextMessage()
	for ok {
		fmt.Println(message)
		message, ok = reader.NextMessage()
	}
}

*/
