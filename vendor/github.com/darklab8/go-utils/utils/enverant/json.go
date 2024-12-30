package enverant

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

var EnverantDebug = os.Getenv("ENVERANT_PRINT_JSON") == "true"

func ReadJson(path string) map[string]interface{} {
	env_map := make(map[string]interface{})

	if EnverantDebug {
		fmt.Println("enverant path=", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Println(err, "not found env file at path=", path)
		return env_map
	}

	err = yaml.Unmarshal(data, &env_map)
	if err != nil {
		panic(fmt.Sprintln(err, "failed to parse yaml/json with env vars"))
	}

	for key, value := range env_map {

		if _, ok := os.LookupEnv(key); ok {
			continue
		}

		switch v := value.(type) {
		case bool:
			os.Setenv(key, strconv.FormatBool(v))
		case string:
			os.Setenv(key, v)
		case int:
			os.Setenv(key, fmt.Sprintf("%d", v))
		case float64:
			os.Setenv(key, fmt.Sprintf("%.0f", v))
		default:
			panic(fmt.Sprintln("enverant value in file has not supported type", key, fmt.Sprintf("%T", value)))
		}

	}

	return env_map
}
