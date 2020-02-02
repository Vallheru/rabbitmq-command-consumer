package main

import (
	"fmt"
	"io/ioutil"
)

/**


// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Fatalf("%s: %s", msg, err)
// 	}
// }

 */


func main() {
	data, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	content := string(data)
	
	fmt.Println(parseConfigString(&content))
}
