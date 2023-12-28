package main

import (
	"log"
	"flag"
	"fmt"

)

//enter port numbers
const frontend = "3000"
const backend = "2999"

//command to seed database: ./bin/backend --seed true
func main() {

	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()

	data, err := NewPostgressStore()
	if err != nil {
		log.Fatal(err)
	}
	
	if err := data.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("seeding the database")
		SeedData(data)
	}

	server := NewAPIServer(":" + backend, data)
	server.Run()

}