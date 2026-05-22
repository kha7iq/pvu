package k8s

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/client-go/kubernetes"

	"github.com/kha7iq/pvu/internal/models"
)

func GetPodVolumeView(ctx context.Context, clientset kubernetes.Interface, namespace, podName string) (models.PodVolumeView, error) {
	pod, err := GetPod(ctx, clientset, namespace, podName)
	if err != nil {
		return models.PodVolumeView{}, err
	}

	if pod.NodeName == "" {
		return models.PodVolumeView{}, fmt.Errorf("pod %q is not scheduled on a node yet", pod.Name)
	}

	req := clientset.CoreV1().RESTClient().Get().
		Resource("nodes").
		Name(pod.NodeName).
		SubResource("proxy").
		Suffix("stats/summary")

	res, err := req.DoRaw(ctx)
	if err != nil {
		return models.PodVolumeView{}, fmt.Errorf("query node stats for %q: %w", pod.NodeName, err)
	}

	var summary models.Summary
	if err := json.Unmarshal(res, &summary); err != nil {
		return models.PodVolumeView{}, fmt.Errorf("parse kubelet stats summary: %w", err)
	}

	view := models.PodVolumeView{
		PodName:   pod.Name,
		Namespace: pod.Namespace,
		NodeName:  pod.NodeName,
		Volumes:   []models.VolumeRow{},
	}

	for _, podStat := range summary.Pods {
		if podStat.PodRef.Name != pod.Name || podStat.PodRef.Namespace != pod.Namespace {
			continue
		}

		for _, vol := range podStat.Volume {
			if vol.PvcRef == nil {
				continue
			}

			var usagePercent float64
			if vol.CapacityBytes > 0 {
				usagePercent = float64(vol.UsedBytes) / float64(vol.CapacityBytes) * 100
			}

			view.Volumes = append(view.Volumes, models.VolumeRow{
				VolumeName:     vol.Name,
				PVCName:        vol.PvcRef.Name,
				CapacityBytes:  vol.CapacityBytes,
				UsedBytes:      vol.UsedBytes,
				AvailableBytes: vol.AvailableBytes,
				UsagePercent:   usagePercent,
			})
		}

		return view, nil
	}

	return models.PodVolumeView{}, fmt.Errorf("pod %q found, but kubelet stats are not available yet", pod.Name)
}
