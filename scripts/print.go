package main

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(any interface{}) {
	b, err := json.MarshalIndent(any, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
