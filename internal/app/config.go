package app

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	GrpcServerPort int
	HttpServerPort int
	DbFile         string
}

// Flag-feature part:

// FlagSet is Config factory
type FlagSet struct {
	cfg Config
	fs  *flag.FlagSet
}

func InitFlagSet() *FlagSet {
	cfgFs := new(FlagSet)
	cfgFs.fs = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	cfgFs.fs.IntVar(&cfgFs.cfg.GrpcServerPort, flagNameGrpcPort, 0, "port for gGRPC server (clients api)")
	cfgFs.fs.IntVar(&cfgFs.cfg.HttpServerPort, flagNameHttpPort, 0, "port for HTTP server (admin web site)")
	cfgFs.fs.StringVar(&cfgFs.cfg.DbFile, flagDbFile, "", "path to sqlite database source")
	return cfgFs
}

// cleanLastCfg clean parsed values
func (c *FlagSet) cleanLastCfg() {
	c.cfg.GrpcServerPort = 0
	c.cfg.HttpServerPort = 0
	c.cfg.DbFile = ""
}

// Flag names:

const (
	flagNameGrpcPort = "grpc-port"
	flagNameHttpPort = "http-port"
	flagDbFile       = "db"
)

// Usage printing "how usage flags" message
func (c *FlagSet) Usage() { c.fs.Usage() }

// Parse is Config factory method
func (c *FlagSet) Parse(args []string) (*Config, error) {
	missingArgExit := func(argName string) error {
		return fmt.Errorf("missing `%s` argument: %w", argName, ErrorWrongArgument)
	}

	err := c.fs.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse given config arguments: %w", err)
	}
	newCfg := c.cfg // copy parsed values
	c.cleanLastCfg()

	// check what all fields defined
	switch {
	case newCfg.GrpcServerPort == 0:
		return nil, missingArgExit(flagNameGrpcPort)
	case newCfg.HttpServerPort == 0:
		return nil, missingArgExit(flagNameHttpPort)
	case newCfg.DbFile == "":
		return nil, missingArgExit(flagDbFile)
	}

	return &newCfg, nil
}

var ErrorWrongArgument = errors.New("wrong argument")
