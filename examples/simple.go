package main

import (
	"context"
	"log"
	"time"

	"github.com/SumoLogic-Labs/props"
)

func simpleExample() {
	cache := props.Cache{
		Store:           props.GetProperties(),
		Source:          props.NewFileSource("example.props"),
		RefreshInterval: time.Minute,
		ExpireAfter:     2 * time.Minute,
	}

	if err := cache.Start(context.Background().Done()); err != nil {
		log.Fatal(err)
	}

	log.Println(
		props.GetString("my_string", "default"),
		props.GetInt("my_number", 0),
		props.GetBool("prop_in_local", false),
	)
}
