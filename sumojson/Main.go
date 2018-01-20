package sumojson

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/SumoLogic/sumoshell/util"
)

func BuildAndConnect(args []string) {
	if len(args[1:]) > 0 {
		read(args[1])
	} else {
		read("")
	}
}

func read(filterString string) {
	r, w := io.Pipe()
	handler := util.NewRawInputHandler(w)
	go util.ConnectToReader(JsonOperator{util.NewJsonWriter()}, r)
	bio := bufio.NewReader(os.Stdin)
	var line, hasMoreInLine, err = bio.ReadLine()
	for err != io.EOF || hasMoreInLine {
		handler.Process(line)
		line, hasMoreInLine, err = bio.ReadLine()
	}
	handler.Flush()
}

type JsonOperator struct {
	output *util.JsonWriter
}

func (jsonOperator JsonOperator) Process(inp map[string]interface{}) {
	raw := util.ExtractRaw(inp)
	if raw == "" {
		return
	}
	rawBytes := []byte(raw)
	var objmap map[string]interface{}
	err := json.Unmarshal(rawBytes, &objmap)
	if err != nil {
		inp["parse_failure"] = raw

	}
	for k, v := range objmap {
		inp[k] = v
	}
	jsonOperator.output.Write(inp)
}
