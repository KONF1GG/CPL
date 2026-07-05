package handler

type CreateVMRequest struct {
	Name   string `json:"name"`
	CPU    int    `json:"cpu"`
	RamMB  int    `json:"ram_mb"`
	DiskGB int    `json:"disk_gb"`
}

