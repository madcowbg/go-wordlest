package game

import (
	"encoding/csv"
	"log"
	"os"
)

type WordList []Word

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Can't close file "+filePath, err)
		}
	}(f)

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func ReadWordleList(filePath string) WordList {
	loaded_csv := readCsvFile(filePath)
	wordlist := make([]Word, len(loaded_csv))

	for i := range loaded_csv {
		wordlist[i] = toWord(loaded_csv[i][0])
	}
	return wordlist
}
