package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
)

var optSize = flag.Int("size", 1, "Number of objects")

type Data struct {
	Number int    `json:"number"`
	Text   string `json:"text"`
}

type DataList []*Data

func (this *DataList) String() string {
	out := "["
	for i, v := range *this {
		if i != 0 {
			out += ","
		}
		out += fmt.Sprintf("%+v", *v)
	}
	out += "]"
	return out
}

func raw(n int) []byte {
	raw := "["
	for i := 0; i < n; i++ {
		if i != 0 {
			raw += ","
		}
		raw += fmt.Sprintf(`{"number":%d,"text":"%s"}`, i+1, "xyz")
	}
	raw += "]"
	return []byte(raw)
}

func data(n int) DataList {
	data := make(DataList, n)

	for i := 0; i < n; i++ {
		data[i] = &Data{Number: i + 1, Text: "xyz"}
	}

	return data
}

func jsonParse(raw *[]byte) DataList {
	var result DataList
	err := json.Unmarshal(*raw, &result)
	if err != nil {
		log.Println("Error parsing JSON", err)
		return nil
	}
	return result
}

func jsonRaw(data *DataList) []byte {
	raw, err := json.Marshal(*data)
	if err != nil {
		log.Println("Error generating JSON", err)
		return nil
	}
	return raw
}

func main() {
	fmt.Println("JSON encode/decode")
	flag.Parse()

	fmt.Println("Size:", *optSize)

	fmt.Println("(Raw to Data)")
	raw := raw(*optSize)
	result1 := jsonParse(&raw)
	fmt.Println(&result1)

	fmt.Println("(Data to Raw)")
	data := data(*optSize)
	result2 := jsonRaw(&data)
	fmt.Println(string(result2))
}
