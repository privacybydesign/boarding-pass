package main

import (
	log "boarding-pass/logging"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	ServerConfig ServerConfig `json:"server_config"`
}

func main() {

	configPath := flag.String("config", "", "Path for the config.json to use")
	flag.Parse()

	if *configPath == "" {
		log.Error.Fatal("please provide a config path using the --config flag")
	}

	log.Info.Printf("using config: %v", *configPath)

	config, err := readConfigFile(*configPath)
	if err != nil {
		log.Error.Fatalf("failed to read config file: %v", err)
	}

	server := NewServer(&ServerState{irmaServerURL: "http://localhost:8080"}, &config.ServerConfig)

	log.Info.Println("starting server on " + config.ServerConfig.Host + ":" + strconv.Itoa(config.ServerConfig.Port))
	err = server.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Error.Fatalf("server failed: %v", err)
	}

}

func readConfigFile(path string) (Config, error) {
	configBytes, err := os.ReadFile(path)

	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(configBytes, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
