package model

type AuthenticationResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type Metadata struct {
	Guid string
}

type Entity struct {
	Name string
}

type Resource struct {
	Metadata Metadata
	Entity   Entity
}

type ApiResponse struct {
	Resources []Resource
}

type AppStat struct {
	Instances []string
	Data      map[string]Instance
}

type Instance struct {
	State string `json:"state"`
	Stats struct {
		Name                string   `json:"name"`
		URIs                []string `json:"uris"`
		Host                string   `json:"host"`
		Port                int      `json:"port"`
		Uptime              int64    `json:"uptime"`
		MemoryQuota         int64    `json:"mem_quota"`
		DiskQuota           int64    `json:"disk_quota"`
		FiledescriptorQuota int      `json:"fds_quota"`
		Usage               struct {
			Time   string  `json:"time"`
			CPU    float64 `json:"cpu"`
			Memory int64   `json:"mem"`
			Disk   int64   `json:"disk"`
		} `json:"usage"`
	} `json:"stats"`
}