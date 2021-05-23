package client

import (
	"context"
	"log"

	"github.com/lwj5/ephemeral-launcher/pkg/launcher"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	batchv1type "k8s.io/client-go/kubernetes/typed/batch/v1"
)

type Status int

const (
	Running Status = iota - 1
	Failed
	Succeeded
)

func (i Status) String() string {
	return [...]string{"Running", "Failed", "Succeeded"}[i+1]
}

func GetJobClient(kubeClient *kubernetes.Clientset, namespace string) batchv1type.JobInterface {
	return kubeClient.BatchV1().Jobs(namespace)
}

func GetJobStatus(job *batchv1.Job) Status {
	if job.Status.Active > 0 {
		return Running
	} else {
		if job.Status.Succeeded > 0 {
			return Succeeded
		}
		return Failed
	}
}

func ListJobs(jobsClient batchv1type.JobInterface) (*batchv1.JobList, error) {
	// List Jobs
	list, err := jobsClient.List(context.TODO(), metav1.ListOptions{})
	return list, err
}

func CreateJob(jobsClient batchv1type.JobInterface, config *launcher.Configuration, job *batchv1.Job) error {
	// Create Job
	log.Printf("Creating job...")
	result, err := jobsClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Unable to create job %q.\n", job.ObjectMeta.GetName())
		return err
	}
	log.Printf("Created job %q.\n", result.GetObjectMeta().GetName())

	return nil
}

func CreateLauncherJob(jobsClient batchv1type.JobInterface, config *launcher.Configuration) error {

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "launcher",
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "launcher",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "launcher",
							Image: "lwj5/ephemeral-launcher",
							Env: []apiv1.EnvVar{
								{
									Name: "LAUNCHERNAMESPACE",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name:  "LAUNCHERREPO_URL",
									Value: config.RepoURL,
								},
								{
									Name:  "LAUNCHERCHART_RELEASE_NAME",
									Value: config.ChartReleaseName,
								},
								{
									Name:  "LAUNCHERCHART_NAME",
									Value: config.ChartName,
								},
								{
									Name:  "LAUNCHERCHART_VERSION",
									Value: config.ChartVersion,
								},
							},
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
				},
			},
		},
	}

	return CreateJob(jobsClient, config, job)

}
