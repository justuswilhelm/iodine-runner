package lib

import (
	"bufio"
	"log"
)

func readStderr(reader *bufio.Reader, name string, matcher func(string)) {
	log.Printf("Started logger for %s", name)
	for {
		output, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from stderr: %+v", err)
			log.Printf("Closing logger for %s", name)
			return
		}
		log.Printf("%s: %s", name, output)
		matcher(output)
	}
}
