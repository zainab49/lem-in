package main
import "functions"
import (
	"log"
	"os"
)

func ReadData(filename string) ()

func main() {
	if len(os.Args) !=2 {
		log.Fatal("Invalid Arguments")
		return
	}
	graph, err := ReadData(os.Args[1])
}