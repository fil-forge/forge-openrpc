package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alanshaw/go-openrpc"
	"github.com/fil-forge/forge-openrpc/services/upload"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run ./main.go <service>")
		return
	}
	switch os.Args[1] {
	case "upload":
		schema, err := upload.BuildSchema()
		if err != nil {
			panic(fmt.Errorf("building upload schema: %w", err))
		}
		writeSchema(schema)
	default:
		fmt.Printf("Unknown service: %s\n", os.Args[1])
	}
}

func writeSchema(schema *openrpc.Schema) {
	b, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		panic(fmt.Errorf("marshaling openrpc schema: %w", err))
	}
	fmt.Println(string(b))
}
