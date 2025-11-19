package environment

import (
	"errors"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadVars() (map[string]string, []string, error) {
	nytKey := os.Getenv("NYT_API_KEY")
	if nytKey == "" {
		return nil, nil, errors.New("new york times api key not found")
	}

	polygonKeys := []string{}
	polygonEnvNames := []string{
		"POLYGON_API_KEY",
		"POLYGON_API_KEY_1",
		"POLYGON_API_KEY_2",
		"POLYGON_API_KEY_3",
		"POLYGON_API_KEY_4",
		"POLYGON_API_KEY_5",
	}
	for _, name := range polygonEnvNames {
		v := os.Getenv(name)
		if v == "" {
			return nil, nil, errors.New(name + " not found")
		}
		polygonKeys = append(polygonKeys, v)
	}

	geminiKey := os.Getenv("GOOGLE_GEMINI_API_KEY")
	if geminiKey == "" {
		return nil, nil, errors.New("gemini api key not found")
	}

	mongoUsername := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	if mongoUsername == "" {
		return nil, nil, errors.New("mongo username not found")
	}

	mongoPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	if mongoPassword == "" {
		return nil, nil, errors.New("mongo password not found")
	}

	mongoPortString := os.Getenv("MONGO_PORT")
	if mongoPortString == "" {
		return nil, nil, errors.New("mongo port not found")
	}

	// validate port is an integer but keep it as string in the returned map
	if _, err := strconv.Atoi(mongoPortString); err != nil {
		return nil, nil, errors.Join(errors.New("failed to convert mongo port to int"), err)
	}

	throttleTimeString := os.Getenv("THROTTLE_TIME")
	if throttleTimeString == "" {
		return nil, nil, errors.New("throttle time not found")
	}

	// validate time is an integer but keep it as string in the returned map
	if _, err := strconv.Atoi(throttleTimeString); err != nil {
		return nil, nil, errors.Join(errors.New("failed to convert throttle time to int"), err)
	}

	vars := map[string]string{
		"NYT_API_KEY":                nytKey,
		"GOOGLE_GEMINI_API_KEY":      geminiKey,
		"MONGO_INITDB_ROOT_USERNAME": mongoUsername,
		"MONGO_INITDB_ROOT_PASSWORD": mongoPassword,
		"MONGO_PORT":                 mongoPortString,
		"THROTTLE_TIME":              throttleTimeString,
	}

	return vars, polygonKeys, nil
}

// Loads env based on whether we are in a Docker container or not
func LoadEnvironment() error {
	// Load env vars, if not in Docker container
	log.Println("Checking if Docker is running...")
	if dockerRunning := os.Getenv("DOCKER_RUNNING"); dockerRunning != "true" {
		if err := LoadLocalEnvironment(); err != nil {
			return errors.Join(errors.New("couldn't load environment"), err)
		}
	}
	return nil
}

// LoadLocalEnvironment gets environment variables from a .env file in the project directory and allows us to use them
// Falls through if in Docker environment, to protect against infinite loops
func LoadLocalEnvironment() error {
	log.Println("Double-checking if Docker is running...")
	// Safeguard to prevent infinite loops
	if dockerRunning := os.Getenv("DOCKER_RUNNING"); dockerRunning == "true" {
		return nil
	}

	log.Println("Loading environment variables...")
	// Find root directory
	for {
		workingDirectoryPath, err := os.Getwd()
		if err != nil {
			return errors.Join(errors.New("could not get current working directory for environment file loading"), err)
		}
		_, workingDirectoryName := path.Split(workingDirectoryPath)
		log.Println("Current directory:", workingDirectoryName)
		if workingDirectoryName == "backend" {
			break
		}
		os.Chdir("..")
	}
	err := godotenv.Load(".env", "mongo.env")
	if err != nil {
		return errors.Join(errors.New("could not load environment variables"), err)
	}
	log.Println("Environment variables loaded.")

	return nil
}
