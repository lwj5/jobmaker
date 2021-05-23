package main

import (
	"context"
	"log"

	"github.com/lwj5/ephemeral-launcher/pkg/launcher"
	"github.com/lwj5/jobmaker/pkg/client"
	"github.com/lwj5/jobmaker/pkg/jobmaker"
)

// CreateLauncherJob implements jobmaker.jobmakerServer
func (s *server) CreateLauncherJob(ctx context.Context, in *jobmaker.CreateLauncherJobRequest) (*jobmaker.JobResponse, error) {
	log.Printf("Received: %v", in.GetRepoURL())

	kubeClient := client.GetClient()
	jobsClient := client.GetJobClient(kubeClient, in.GetNamespace())
	if err := client.CreateLauncherJob(jobsClient, &launcher.Configuration{
		Namespace:        in.GetNamespace(),
		RepoURL:          in.GetRepoURL(),
		ChartReleaseName: in.GetChartReleaseName(),
		ChartName:        in.GetChartName(),
		ChartVersion:     in.GetChartVersion(),
	}); err != nil {
		return &jobmaker.JobResponse{}, err
	}

	return &jobmaker.JobResponse{Code: 200, Message: "Hello " + in.GetChartName()}, nil
}
