package client

import (
	"bytes"
	"context"
	"io"

	corev1 "k8s.io/api/core/v1"
	corev1type "k8s.io/client-go/kubernetes/typed/core/v1"
)

func getPodLogs(pod corev1.Pod, podClient corev1type.PodInterface) string {
	podLogOpts := corev1.PodLogOptions{}

	req := podClient.GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()

	return str
}
