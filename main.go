package main

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type NFT struct {
	Format           string      `json:"format"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	MintingTool      string      `json:"minting_tool"`
	SensitiveContent bool        `json:"sensitive_content"`
	SeriesNumber     string      `json:"series_number"`
	SeriesTotal      int         `json:"series_total"`
	Collection       Collections `json:"collection"`
	Gender           string      `json:"gender"`
	Uuid             string      `json:"uuid"`
	Hash             string      `json:"hash"`
}

type Collections struct {
	Name       string       `json:"name"`
	ID         string       `json:"id"`
	Attributes []Attributes `json:"attributes"`
}
type Attributes struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

var collection []string

func main() {
	data := ReadCsv()
	nft := GenerateJSON(data)
	GenerateCSV(nft)
}

func ReadCsv() [][]string {
	var sun string
	flag.StringVar(&sun, "file", "", "enter path to the file")
	flag.Parse()

	file, err := os.Open(strings.TrimSpace(sun))
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return data
}

func GenerateCSV(team []NFT) {
	os.Chdir("../HNGi9 csv")
	file, err := os.OpenFile("nft.output.csv", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Println(err.Error())
		return
	}
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()
	var data [][]string
	data = append(data, []string{"Series Number", "Filename", "Description", "Gender", "Attribute", "UUID", "Hash"})
	for i, v := range team {
		row := []string{v.SeriesNumber, v.Collection.Name, v.Description, v.Gender, fmt.Sprintf("%v", collection[i]), v.Uuid, v.Hash}
		data = append(data, row)

	}
	err = csvWriter.WriteAll(data)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func GenerateJSON(data [][]string) []NFT {
	var NFTs []NFT
	os.Mkdir("HNGi9 csv", 0777)
	os.Mkdir("teams json data", 0777)
	os.Chdir("teams json data")
	for i, t := range data {
		if i == 0 {
			continue
		}
		var attribute Attributes
		var attributes []Attributes
		newAttribute := strings.Split(t[6], ";")

		for _, a := range newAttribute {
			attributeSlice := strings.Split(a, ":")
			attribute.Type = attributeSlice[0]
			attribute.Value = strings.Join(attributeSlice[1:], " ")
			attributes = append(attributes, attribute)
		}
		collection = append(collection, t[6])
		nft := NFT{
			Format:           "CHIP-0007",
			Name:             t[3],
			Description:      t[4],
			MintingTool:      "SuperMinter/2.5.2",
			SensitiveContent: false,
			SeriesNumber:     t[1],
			SeriesTotal:      len(data),
			Collection: Collections{
				Name:       t[2],
				ID:         strconv.Itoa(i + 1),
				Attributes: attributes,
			},
			Gender: t[5],
			Uuid:   t[7],
		}
		h := sha256.New()
		jsonData, _ := json.Marshal(nft)
		h.Write(jsonData)
		nft.Hash = hex.EncodeToString(h.Sum(nil))

		jsonData, err := json.Marshal(&nft)
		if err != nil {
			log.Println(err)
			return nil
		}
		file, err := os.OpenFile(t[0]+".output.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		_, err = file.Write(jsonData)
		if err != nil {
			log.Println(err)
			return nil
		}
		NFTs = append(NFTs, nft)
	}
	return NFTs
}
