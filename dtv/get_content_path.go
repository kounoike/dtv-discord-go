package dtv

import (
	"bytes"
	"context"

	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/template"
)

func (dtv *DTVUsecase) getContentPath(ctx context.Context, program db.Program, service db.Service) (string, error) {
	data := template.PathTemplateData{}

	_ = dtv.gpt.ParseTitle(ctx, program.Name, &data)

	data.Program = template.PathProgram{
		Name:      program.Name,
		StartTime: program.StartTime(),
	}

	data.Service = template.PathService{
		Name: service.Name,
	}

	var buffer bytes.Buffer
	err := dtv.contentPathTmpl.Execute(&buffer, data)
	if err != nil {
		return "", err
	}
	contentPath := toSafePath(buffer.String())
	return contentPath, nil
}

func (dtv *DTVUsecase) getOutputPath(ctx context.Context, program db.Program, service db.Service) (string, error) {
	var b bytes.Buffer
	data := template.PathTemplateData{}

	_ = dtv.gpt.ParseTitle(ctx, program.Name, &data)

	data.Program = template.PathProgram{
		Name:      program.Name,
		StartTime: program.StartTime(),
	}
	data.Service = template.PathService{
		Name: service.Name,
	}

	err := dtv.outputPathTmpl.Execute(&b, data)
	if err != nil {
		return "", err
	}
	outputPath := b.String()

	return outputPath, nil
}
