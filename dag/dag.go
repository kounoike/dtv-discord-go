package dag

type OutputType string

const (
	Video OutputType = "video"
	Text  OutputType = "text"
	Other OutputType = "other"
)

type Output struct {
	Type OutputType `yaml:"type"`
	Name string     `yaml:"name"`
	Path string     `yaml:"path"`
}

type Step struct {
	Name    string   `yaml:"name"`
	Command string   `yaml:"command"`
	Queue   string   `yaml:"queue"`
	Output  Output   `yaml:"output"`
	Depends []string `yaml:"depends"`
}

type DagFile struct {
	Steps []*Step `yaml:"steps"`
}
