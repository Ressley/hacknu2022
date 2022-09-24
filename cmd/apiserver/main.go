package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/Ressley/hacknu/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	/*
		client, err := apiserver.GetMongoClient()
		if err != nil {
			log.Fatal(err)
		}
		accounts := client.Database(apiserver.DB).Collection(apiserver.ACCOUNTS)
		accounts.InsertOne(context.TODO(), bson.D{
			{"title", "The Polyglot Developer Podcast"},
			{"author", "Nic Raboy"},
			{"tags", bson.A{"development", "programming", "coding"}},
		})
	*/
	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
