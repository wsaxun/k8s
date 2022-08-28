package utils

import (
	"context"
	"fmt"
	"github.com/apenella/go-ansible/pkg/adhoc"
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
