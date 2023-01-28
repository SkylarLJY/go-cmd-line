package main

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type pipelineStep struct {
	Name, Cmd, Message        string
	Timeout, Exception, Print bool
}

func initPipeline(file, proj, branch string) ([]executer, error) {
	configs, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	configList := []pipelineStep{}
	err = json.Unmarshal(configs, &configList)
	if err != nil {
		return nil, err
	}
	res := []executer{}
	for _, line := range configList {
		var thisStep executer
		exe := strings.Split(line.Cmd, " ")[0]
		args := strings.Split(line.Cmd, " ")[1:]
		if line.Timeout {
			thisStep = newTimeoutStep(line.Name, exe, line.Message, proj, append(args, branch), 10*time.Second)
		} else if line.Exception {
			thisStep = newExceptionStep(line.Name, exe, line.Message, proj, args)
		} else if line.Print {
			thisStep = newPrintStep(line.Name, exe, line.Message, proj, args)
		} else {
			thisStep = newStep(line.Name, exe, line.Message, proj, args)
		}

		res = append(res, thisStep)
	}
	return res, nil
}
