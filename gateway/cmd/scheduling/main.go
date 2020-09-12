package main

import (
	"context"
	"time"

	"github.com/tierzer0/gateway/helm"
	"github.com/tierzer0/gateway/scheduling"
)

func main() {
	deployer := helm.CliDeployer{
		ChartFileName: "/charts/vscode.yaml",
	}

	scheduler := scheduling.NewScheduler(&deployer, 1*time.Millisecond, 5*time.Second)
	ctx, cancel := context.WithCancel(context.Background())

	go scheduler.Run(ctx)

	<-time.AfterFunc(30*time.Second, cancel).C
}
