package helm

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/tierzer0/gateway/deployment"
	"gopkg.in/square/go-jose.v2/json"
)

type CLIDeployer struct {
	ChartPath string
}

func (d CLIDeployer) Deploy(ctx context.Context, deployment deployment.Deployment) error {
	if !d.helmInstalled(ctx) {
		return fmt.Errorf("helm is not installed on this system")
	}

	f, err := ioutil.TempFile("", "values.yaml")
	if err != nil {
		return fmt.Errorf("could not create values file: %w", err)
	}

	defer f.Close()

	err = json.NewEncoder(f).Encode(deployment)
	if err != nil {
		return fmt.Errorf("could not write to values file: %w", err)
	}

	releaseName := fmt.Sprintf("vscode-%s", deployment.ID)

	log.Println(releaseName, d.ChartPath, f.Name())
	cmd := exec.CommandContext(ctx, "helm", "install", releaseName, d.ChartPath, "--debug")
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not execute helm command, check log for full details: %w", err)
	}

	return nil
}

func (d CLIDeployer) helmInstalled(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "helm", "version")
	err := cmd.Run()
	return err == nil
}
