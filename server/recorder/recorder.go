package recorder

import "fmt"

func Recorder(messages chan string) {

	for {

		fmt.Println(<-messages)
	}
}
