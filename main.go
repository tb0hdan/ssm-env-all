package main

import (
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
	"os/exec"
	"fmt"
	"strings"
	"flag"
	"syscall"
)

var param_path = flag.String("path", "/", "Parameter path, separated by slashes: /env/role/param")

func main() {
	flag.Parse()
	if val, status := os.LookupEnv("SSM_PATH"); status && *param_path == "/" {
		*param_path = val
	}
	s := ssm.New(session.New())
	input := &ssm.GetParametersByPathInput{}
	input.SetPath(*param_path)
	input.SetRecursive(true)
	input.SetWithDecryption(true)
	out, err := s.GetParametersByPath(input)
	if err != nil {
		log.Panicf("Could not run GetParametersByPath: %+v\n", err)
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

	for _, param := range out.Parameters {
		key := strings.Split(*param.Name, "/")[len(strings.Split(*param.Name, "/")) - 1]
		key = strings.ToUpper(key)
		val := *param.Value
		env = append(env, fmt.Sprintf("%s=%s", key, val))

	}
	err = syscall.Exec(binary, flag.Args(), env)
	if err != nil {
		log.Panicf("Command exec error: %+v\n", err)
	}
}
