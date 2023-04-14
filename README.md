# Overview

`props` is a thread-safe, multi-source cache for Java properties in Go.

It wraps around the [properties](https://github.com/magiconair/properties) library and provides a read-only concurrent cache that is continusously kept in sync with configured properties sources.

## Features

1. Supports multiple sources:
    1. S3
    1. DynamoDB
    1. Local File
1. Supports a combination of any of the above with custom precedence order.
1. Caching and syncing with remote sources.
1. Go-routine/concurrent/thread-safe global props.

You can also bring your own source by implementing the `Poller` interface and use it during cache initialization:
```go
type Poller interface {
	Poll(context.Context) (*properties.Properties, error)
}
```

Feel free to send in a PR for some useful sources!

## Example Usage

Configure Sources

```go
source1 := props.NewFileSource("path_to_local_props")

source2 := props.NewS3Source(*cfg, "props_bucket", "props_file_path")

source3 := props.NewDynamoDBGetterSource(*cfg, props.DynamoDBGetterArgs{
  Table:     "PropsTable",
  KeyCol:    "Key",
  ValCol:    "Value",
  WatchKeys: []string{"prop_key_1", "prop_key_2"},
})
```

Initialize and Start Cache

```go
cache := props.Cache{
  Store:           props.GetProperties(),
  Source:          props.NewCompositeSource(source1, source2, source3),
  RefreshInterval: time.Minute,
  ExpireAfter:     5 * time.Minute,
}

cache.Start(context.Background().Done())
```

Read props from anywhere, anytime

```
valFromLocal := props.GetBool("prop_in_local", false)
valFromS3 := props.GetString("prop_in_s3", "")
valFromDb := props.GetInt("prop_key_1", 42)

log.Println(valFromLocal, valFromS3, valFromDb)
```
