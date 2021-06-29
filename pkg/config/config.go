package config

import (
	"log"
	"os"

	"github.com/go-openapi/swag"
	"github.com/jessevdk/go-flags"
)

// Options are configuration flags for Server
var Options struct {
	DBSource string `long:"db_source" description:"Database SourceName" default:"root:@tcp(localhost:3306)/product_category" env:"DB_SOURCE"`
	Host     string `long:"host" description:"the IP to listen on" default:"localhost" env:"HOST"`
	Port     int    `long:"port" description:"port to bind" default:"8000" env:"PORT"`

	// development logging mode
	DevMode bool `long:"dev_mode" description:"development mode" env:"DEV_MODE"`
}

// setConfigurations set all available server option list.
func setConfigurations() []swag.CommandLineOptionsGroup {
	return []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "Server Configuration",
			LongDescription:  "Server Configuration",
			Options:          &Options,
		},
		// more configuration options..
	}
}

// ParseConfigurations parse server configurations.
func ParseConfigurations() {
	// 1. set available server configurations.
	configurations := setConfigurations()
	// 2. Parse command line flags
	parser := flags.NewParser(nil, flags.Default)
	for _, optsGroup := range configurations {
		if _, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options); err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			code = 0
		}
		os.Exit(code)
	}
}
