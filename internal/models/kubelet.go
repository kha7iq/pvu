package models

type Summary struct {
	Pods []PodStats `json:"pods"`
}

type PodStats struct {
	PodRef PodRef        `json:"podRef"`
	Volume []VolumeStats `json:"volume"`
}

type PodRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type VolumeStats struct {
	Name           string  `json:"name"`
	AvailableBytes uint64  `json:"availableBytes"`
	CapacityBytes  uint64  `json:"capacityBytes"`
	UsedBytes      uint64  `json:"usedBytes"`
	PvcRef         *PvcRef `json:"pvcRef,omitempty"`
}

type PvcRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
