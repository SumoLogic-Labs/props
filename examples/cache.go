package main

import (
	"context"
	"log"
	"time"

	"github.com/SumoLogic-Labs/props"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func cacheExample() {
	initCache()
	valFromLocal := props.GetBool("prop_in_local", false)
	valFromS3 := props.GetString("prop_in_s3", "")
	valFromDb := props.GetInt("prop_key_1", 42)
	log.Println(valFromLocal, valFromS3, valFromDb)
}

func initCache() {
	cfg := *aws.NewConfig()

	source1 := props.NewFileSource("path_to_local_props")

	source2 := props.NewS3Source(cfg, "props_bucket", "props_file_path")

	source3 := props.NewDynamoDBGetterSource(cfg, props.DynamoDBGetterArgs{
		Table:     "props_table",
		KeyCol:    "key",
		ValCol:    "value",
		WatchKeys: []string{"prop_key_1", "prop_key_2"},
	})

	cache := props.Cache{
		Store:           props.GetProperties(),
		Source:          props.NewCompositeSource(source1, source2, source3),
		RefreshInterval: time.Minute,
		ExpireAfter:     5 * time.Minute,
	}

	if err := cache.Start(context.Background().Done()); err != nil {
		log.Fatal(err)
	}
}
