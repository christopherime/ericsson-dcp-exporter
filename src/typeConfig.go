package main

type TypeConfig struct {
	// The name of the config file
	ConfigFile    string
	configPresent bool
	configValues  configValueObject
}

type configValueObject struct {
	SubToWatch []string `yaml:"subToWatch"`
}
