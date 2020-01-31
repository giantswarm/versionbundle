package versionbundle

type AppsChangelogs map[string][]AppChangelog

type AppChangelog struct {
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Urls        []string `json:"urls"`
	Version     string   `json:"version"`
}

type ReleasesChangelogs map[string]AppsChangelogs
