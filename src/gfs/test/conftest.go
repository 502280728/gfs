package main

import (
	"fmt"
	"gfs/common"
	"gopkg.in/yaml.v2"
	"log"
)

var data = `
node:
  datadir: /data1
  blocksize: 24M
  infointerval: 10
  masters: ["localhost:9091","localhost:9092"]
  port: 8078

master: 
  defaultfs: localhost:9091
  defaultdir: /data2
`

func main() {
	fmt.Print("start")
	t := common.Conf{}
	err := yaml.Unmarshal([]byte(data), &t)
	fmt.Printf("--- t:\n%v\n\n", t)
	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))
}
