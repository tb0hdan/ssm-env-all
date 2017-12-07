package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	ignore_errors = flag.Bool("ignore-errors", false, "Ignore authentication/authorization errors")
	param_path    = flag.String("path", "/", "Parameter path, separated by slashes: /env/role/param")
)

func populateEnvFromPath(path, token string, ignore_errors bool) (result []string) {
	s := ssm.New(session.New())
	input := &ssm.GetParametersByPathInput{}
	if token != "" {
		input.SetNextToken(token)
	}
	input.SetPath(path)
	input.SetRecursive(true)
	input.SetWithDecryption(true)
	out, err := s.GetParametersByPath(input)

	if out.NextToken != nil {
		result = populateEnvFromPath(path, *out.NextToken, ignore_errors)
	}

	if err != nil {
		message := fmt.Sprintf("Could not run GetParametersByPath: %+v\n", err)
		if ignore_errors {
			log.Print(message)
			return
		} else {
			log.Panic(message)
		}
	}
	for _, param := range out.Parameters {
		key := strings.Split(*param.Name, "/")[len(strings.Split(*param.Name, "/"))-1]
		key = strings.ToUpper(key)
		val := *param.Value
		result = append(result, fmt.Sprintf("%s=%s", key, val))
	}
	return
}

func main() {
	flag.Parse()
	if val, status := os.LookupEnv("SSM_PATH"); status && *param_path == "/" {
		*param_path = val
	}

	if val, status := os.LookupEnv("SSM_IGNORE_ERRORS"); status {
		*ignore_errors = (val == "1") || (strings.ToLower(val) == "yes")
	}

	lookup := "env"
	if len(flag.Args()) > 0 {
		lookup = flag.Args()[0]
	}
	binary, err := exec.LookPath(lookup)
	if err != nil {
		log.Panicf("Command error: %+v\n", err)
	}

	env := os.Environ()

	for _, path := range strings.Split(*param_path, ",") {
		for _, envvar := range populateEnvFromPath(path, "", *ignore_errors) {
			env = append(env, envvar)
		}
	}

	err = syscall.Exec(binary, flag.Args(), env)
	if err != nil {
		log.Panicf("Command exec error: %+v\n", err)
	}
}
