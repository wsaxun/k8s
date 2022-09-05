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

func Playbook(yml string, hosts string) {
	var err error
	var res *results.AnsiblePlaybookJSONResults

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: hosts,
	}

	buff := new(bytes.Buffer)
	executorTimeMeasurement := measure.NewExecutorTimeMeasurement(
		execute.NewDefaultExecute(
			execute.WithWrite(io.Writer(buff)),
		),
	)

	playbooksList := []string{yml}
	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:      playbooksList,
		Exec:           executorTimeMeasurement,
		StdoutCallback: "json",
		Options:        ansiblePlaybookOptions,
	}

	err = playbook.Run(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	res, err = results.ParseJSONResultsStream(io.Reader(buff))
	fmt.Println(res.String())
	if err != nil {
		log.Fatalln(err)
	}

}
