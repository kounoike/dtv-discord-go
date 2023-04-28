package dtv

import (
	"bytes"
	"context"
	"path"

	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/template"
	"golang.org/x/text/width"
)

const mkDirPerm = 0777

func (dtv *DTVUsecase) getContentPath(ctx context.Context, program db.Program, service db.Service) (string, error) {
	data := template.PathTemplateData{}

	_ = dtv.gpt.ParseTitle(ctx, program.Name, &data)
	data.Title = toSafePath(width.Fold.String(data.Title))
	data.Subtitle = toSafePath(width.Fold.String(data.Subtitle))

	data.Program = template.PathProgram{
		Name:      toSafePath(width.Fold.String(program.Name)),
		StartTime: program.StartTime(),
	}

	data.Service = template.PathService{
		Name: toSafePath(width.Fold.String(service.Name)),
	}

	var buffer bytes.Buffer
	err := dtv.contentPathTmpl.Execute(&buffer, data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (dtv *DTVUsecase) getEncodingOutputPath(contentPath string) string {
	ext := path.Ext(contentPath)
	return contentPath[:len(contentPath)-len(ext)] + dtv.encodedExt
}

func (dtv *DTVUsecase) getTranscriptionOutputPath(contentPath string) string {
	ext := path.Ext(contentPath)
	return contentPath[:len(contentPath)-len(ext)] + dtv.transcribedExt
}

func (dtv *DTVUsecase) getAribB24TextOutputPath(contentPath string) string {
	ext := path.Ext(contentPath)
	return contentPath[:len(contentPath)-len(ext)] + `.arib.txt`
}
