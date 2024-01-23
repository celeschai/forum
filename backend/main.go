package main

import (
	"log"
	"flag"
	"fmt"
	"os"
)

//command to seed database: ./bin/backend --seed true
func main() {
	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()

	data, err := NewPostgressStore()
	if err != nil {
		log.Fatal(err)
	}
	defer data.db.Close()
	
	if err := data.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("seeding the database")
		SeedData(data)
	}

	server := NewAPIServer(":" + os.Getenv("BACK_PORT"), data)
	server.Run()

}