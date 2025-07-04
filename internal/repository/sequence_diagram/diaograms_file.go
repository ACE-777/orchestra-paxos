package sequence_diagram

import (
	"fmt"
	"os"
	"sync"
)

var mu sync.Mutex

func CreateNewFile(paxosType string) {
	err := os.WriteFile(fmt.Sprintf("./artifacts/%s-paxos-output.txt", paxosType), []byte("sequenceDiagram\n"), 0755)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}
}

func WriteToFile(text string) {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile("./artifacts/multi-paxos-output.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		//fmt.Printf("Can't write error: %v", err)
		return
	}
	defer file.Close()
	//mu.Lock()
	_, err = file.WriteString(fmt.Sprintf("%s\n", text))
	//mu.Unlock()
}
