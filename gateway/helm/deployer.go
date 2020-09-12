package helm

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/tierzer0/gateway/scheduling"
)

type CliDeployer struct {
	ChartFileName string
}

func (d *CliDeployer) Deploy(ctx context.Context, deployment scheduling.Deployment) error {
	f, err := ioutil.TempFile("", "values.yml")
	if err != nil {
		return errors.Wrap(err, "could not create values file")
	}

	defer f.Close()

	err = json.NewEncoder(f).Encode(deployment)
	if err != nil {
		return errors.Wrap(err, "could not write values file")
	}

	cmd := exec.Command("helm", "install", d.ChartFileName, "-f", f.Name())

	err = cmd.Run()
	return errors.Wrap(err, "could not run helm install")
}
