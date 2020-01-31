package versionbundle

type AppsLog map[string][]Log

type Log struct {
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Urls        []string `json:"urls"`
	Version     string   `json:"version"`
}

type ReleasesLog map[string]AppsLog
