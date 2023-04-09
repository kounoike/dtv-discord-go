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

func (dtv *DTVUsecase) getEncodingOutputPath(ctx context.Context, program db.Program, service db.Service, pathData *template.PathTemplateData) (string, error) {
	var b bytes.Buffer

	err := dtv.encodingOutputPathTmpl.Execute(&b, pathData)
	if err != nil {
		return "", err
	}
	outputPath := b.String()

	return outputPath, nil
}

func (dtv *DTVUsecase) getTranscriptionOutputPath(ctx context.Context, program db.Program, service db.Service, pathData *template.PathTemplateData) (string, error) {
	var b bytes.Buffer
	err := dtv.transcriptionOutputPathTmpl.Execute(&b, pathData)
	if err != nil {
		return "", err
	}
	outputPath := b.String()

	return outputPath, nil
}
