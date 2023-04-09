package dtv

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (dtv *DTVUsecase) checkMirakcUpdateTask(ctx context.Context) error {
	mirakcVersion, err := dtv.mirakc.GetVersion()
	if err != nil {
		return err
	}
	dtv.logger.Debug("checkMirakcUpdateTask", zap.String("Current", mirakcVersion.Current), zap.String("Latest", mirakcVersion.Latest))
	if mirakcVersion.Current != mirakcVersion.Latest {
		mirakcTableVersion, err := dtv.queries.GetComponentVersion(ctx, "mirakc")
		if errors.Cause(err) == sql.ErrNoRows {
			err = dtv.queries.InsertComponentVersion(ctx, db.InsertComponentVersionParams{Component: "mirakc", Version: mirakcVersion.Latest})
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			if mirakcTableVersion.Version == mirakcVersion.Latest {
				return nil
			}
		}
		content := fmt.Sprintf("mirakcの新しいバージョン(%s)が出ています。現在(%s)", mirakcVersion.Latest, mirakcVersion.Current)
		_, err = dtv.discord.SendMessage(discord.InformationCategory, discord.UpdateChannel, content)
		if err != nil {
			return err
		}
	}
	return nil
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func (dtv *DTVUsecase) checkDtvDiscordGoVersion(ctx context.Context, version string) error {
	url := "https://api.github.com/repos/kounoike/dtv-discord-go/releases/latest"
	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("X-GitHub-Api-Version", "2022-11-28").
		Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode()/100 != 2 {
		dtv.logger.Error("StatusCode error")
		return fmt.Errorf("StatusCode:%d", resp.StatusCode())
	}
	var release GitHubRelease
	json.Unmarshal(resp.Body(), &release)
	ghVersion := release.TagName
	currentVersion := "v" + version
	dtv.logger.Debug("checkDtvDiscordGoVersion", zap.String("Current", currentVersion), zap.String("Latest", ghVersion))
	if ghVersion != currentVersion {
		tVersion, err := dtv.queries.GetComponentVersion(ctx, "dtv-discord-go")
		if errors.Cause(err) == sql.ErrNoRows {
			err := dtv.queries.InsertComponentVersion(ctx, db.InsertComponentVersionParams{Component: "dtv-discord-go", Version: ghVersion})
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if tVersion.Version == ghVersion {
			return nil
		}
	}
	content := fmt.Sprintf("dtv-discord-goの新しいバージョン(%s)が出ています。現在(%s)", ghVersion, currentVersion)
	_, err = dtv.discord.SendMessage(discord.InformationCategory, discord.UpdateChannel, content)
	if err != nil {
		return err
	}
	err = dtv.queries.UpdateComponentVersion(ctx, db.UpdateComponentVersionParams{Component: "dtv-discord-go", Version: ghVersion})
	if err != nil {
		return err
	}
	return nil
}

func (dtv *DTVUsecase) CheckUpdateTask(ctx context.Context, version string) error {
	err := dtv.checkMirakcUpdateTask(ctx)
	if err != nil {
		return err
	}
	err = dtv.checkDtvDiscordGoVersion(ctx, version)
	if err != nil {
		return err
	}
	return nil
}
