package models

type Options struct {
	PodName    string
	Namespace  string
	Kubeconfig string
}

type PodListItem struct {
	Name      string
	Namespace string
	NodeName  string
	Phase     string
}

type VolumeRow struct {
	VolumeName     string
	PVCName        string
	CapacityBytes  uint64
	UsedBytes      uint64
	AvailableBytes uint64
	UsagePercent   float64
}

type PodVolumeView struct {
	PodName   string
	Namespace string
	NodeName  string
	Volumes   []VolumeRow
}
