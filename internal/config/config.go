package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	GrpcServerPort int
	HttpServerPort int
	ClientsCfgFile string
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
	cfgFs.fs.IntVar(&cfgFs.cfg.GrpcServerPort, flagNameGrpcPort, 0, "port for HTTP server (admin web site)")
	cfgFs.fs.IntVar(&cfgFs.cfg.HttpServerPort, flagNameHttpPort, 0, "port for gGRPC server (clients api)")
	cfgFs.fs.StringVar(&cfgFs.cfg.ClientsCfgFile, flagNameClientsCfg, "", "file with clients settings")
	return cfgFs
}

// cleanLastCfg clean parsed values
func (c *FlagSet) cleanLastCfg() {
	c.cfg.GrpcServerPort = 0
	c.cfg.HttpServerPort = 0
	c.cfg.ClientsCfgFile = ""
}

// Flag names:

const (
	flagNameGrpcPort   = "grpc-port"
	flagNameHttpPort   = "http-port"
	flagNameClientsCfg = "clients-cfg"
)

// Usage printing "how usage flags" message
func (c *FlagSet) Usage() { c.fs.Usage() }

// Parse is Config factory method
func (c *FlagSet) Parse(args []string) (*Config, error) {
	missingArgExit := func(argName string) error {
		return fmt.Errorf("missing `%s` argument: %w", argName, WrongArgumentError)
	}

	err := c.fs.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse given config arguments: %w", err)
	}
	newCfg := c.cfg // copy parsed values
	c.cleanLastCfg()

	if newCfg.GrpcServerPort == 0 {
		return nil, missingArgExit(flagNameGrpcPort)
	}
	if newCfg.HttpServerPort == 0 {
		return nil, missingArgExit(flagNameHttpPort)
	}
	if newCfg.ClientsCfgFile == "" {
		return nil, missingArgExit(flagNameClientsCfg)
	}

	return &newCfg, nil
}

var WrongArgumentError = errors.New("wrong argument")
