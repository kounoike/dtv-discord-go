package dag

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kounoike/dtv-discord-go/db"
	"gopkg.in/yaml.v2"
)

type dagTemplateData struct {
	Title    string
	Subtitle string
	Episode  int
	Program  *db.Program
	Service  *db.Service
}

type DagEngine struct {
	BasePath string
}

func NewDagEngine(basePath string) *DagEngine {
	return &DagEngine{BasePath: basePath}
}

func (d *DagEngine) GetDagFile(name string, contentPath string, title string, subtitle string, episode int, program *db.Program, service *db.Service) (*DagFile, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(d.BasePath, name, ".yml"))
	if err != nil {
		if os.IsNotExist(err) {
			bytes, err = ioutil.ReadFile(filepath.Join(d.BasePath, name, ".yaml"))
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	dagData := dagTemplateData{
		Title:    title,
		Subtitle: subtitle,
		Episode:  episode,
		Program:  program,
		Service:  service,
	}

	writer := new(strings.Builder)
	tmpl, err := template.New(fmt.Sprintf("dag_%s", name)).Parse(string(bytes))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(writer, dagData); err != nil {
		return nil, err
	}

	var dagFile DagFile
	if err := yaml.Unmarshal([]byte(writer.String()), &dagFile); err != nil {
		return nil, err
	}

	return nil, nil
}

func (d *DagEngine) ListDagFiles() ([]string, error) {
	return nil, nil
}
