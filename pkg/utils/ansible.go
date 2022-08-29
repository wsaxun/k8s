package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/apenella/go-ansible/pkg/adhoc"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/execute/measure"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"io"
	"log"
)

func Exec(host string, module string, args string) error {

	ansibleAdhocOptions := &adhoc.AnsibleAdhocOptions{
		ModuleName: module,
		Args:       args,
	}

	adhoc := &adhoc.AnsibleAdhocCmd{
		Pattern:        host,
		Options:        ansibleAdhocOptions,
		StdoutCallback: "json",
	}

	fmt.Println("Command: ", adhoc.String())

	err := adhoc.Run(context.TODO())
	return err
}

func Playbook(yml string) {
	var err error
	var res *results.AnsiblePlaybookJSONResults

	buff := new(bytes.Buffer)
	executorTimeMeasurement := measure.NewExecutorTimeMeasurement(
		execute.NewDefaultExecute(
			execute.WithWrite(io.Writer(buff)),
		),
	)

	playbooksList := []string{yml}
	playbooks := &playbook.AnsiblePlaybookCmd{
		Playbooks:      playbooksList,
		Exec:           executorTimeMeasurement,
		StdoutCallback: "json",
	}

	err = playbooks.Run(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	res, err = results.ParseJSONResultsStream(io.Reader(buff))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.String())
}


func Test(yml string) {
	var err error
	var res *results.AnsiblePlaybookJSONResults

	buff := new(bytes.Buffer)



	executorTimeMeasurement := measure.NewExecutorTimeMeasurement(
		execute.NewDefaultExecute(
			execute.WithWrite(io.Writer(buff)),
		),
	)

	playbooksList := []string{yml}
	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:         playbooksList,
		Exec:              executorTimeMeasurement,
		StdoutCallback:    "json",
	}

	err = playbook.Run(context.TODO())
	if err != nil {
		fmt.Println(err.Error())
	}

	res, err = results.ParseJSONResultsStream(io.Reader(buff))
	if err != nil {
		panic(err)
	}

	fmt.Println(res.String())
	fmt.Println("Duration: ", executorTimeMeasurement.Duration())


}

