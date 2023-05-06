package discord_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/discord"
	"go.uber.org/zap"
)

type AutoSearchInfo struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Channel     string `json:"channel"`
	Genre       string `json:"genre"`
	FuzzySearch bool   `json:"fuzzy_search"`
	RegexSearch bool   `json:"regex_search"`
	KanaSearch  bool   `json:"kana_search"`
	Record      bool   `json:"record"`
}

func (h *DiscordHandler) CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	h.logger.Debug("CommandHandler", zap.Any("type", i.Type))
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		h.logger.Debug("InteractionApplicationCommand", zap.String("name", i.ApplicationCommandData().Name))
		switch i.ApplicationCommandData().Name {
		case "index":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "再インデックス処理を開始します",
				},
			})
			if err := h.dtv.Reindex(context.Background()); err != nil {
				h.logger.Error("Reindex error", zap.Error(err))
				return
			}
			s.ChannelMessageSend(i.ChannelID, "再インデックス処理の登録が完了しました(完了には時間がかかります)")
		case "delete":
			asCh, err := h.client.GetCachedChannel(discord.NotifyAndScheduleCategory, discord.AutoSearchChannelName)
			if err != nil {
				h.logger.Error("GetCachedChannel error", zap.Error(err))
				return
			}
			ch, err := h.client.GetChannel(i.ChannelID)
			if err != nil {
				h.logger.Error("GetChannel error", zap.Error(err))
				return
			}
			if asCh.ID != ch.ParentID {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "自動検索スレッドで実行してください",
					},
				})
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "自動検索スレッドを削除します",
				},
			})

			h.dtv.DeleteAutoSearch(context.Background(), i.ChannelID)
		case "create":
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "create",
					Title:    "自動検索の新規作成",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									Label:    "スレッドタイトル",
									CustomID: "name",
									Style:    discordgo.TextInputShort,
									Required: true,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									Label:    "タイトルの検索文字列（空で全番組）",
									CustomID: "title",
									Style:    discordgo.TextInputShort,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									Label:    "チャンネルの検索文字列（空で全チャンネル）",
									CustomID: "channel",
									Style:    discordgo.TextInputShort,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									Label:    "ジャンルの検索文字列（空で全番組）",
									CustomID: "genre",
									Style:    discordgo.TextInputShort,
								},
							},
						},
					},
				},
			}); err != nil {
				h.logger.Error("InteractionRespond error", zap.Error(err))
				return
			}
		}
	case discordgo.InteractionModalSubmit:
		h.logger.Debug("InteractionModalSubmit", zap.String("customID", i.ModalSubmitData().CustomID))
		switch i.ModalSubmitData().CustomID {
		case "create":
			name := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			title := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			channel := i.ModalSubmitData().Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			genre := i.ModalSubmitData().Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			data, err := json.Marshal(AutoSearchInfo{Name: name, Title: title, Channel: channel, Genre: genre})
			if err != nil {
				h.logger.Error("json.Marshal error", zap.Error(err))
				return
			}
			content := fmt.Sprintf("> %s\n検索方法は？", string(data))

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content:  content,
					CustomID: "search_method",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label:    "単純部分一致検索",
									Style:    discordgo.PrimaryButton,
									CustomID: "search_normal",
								},
								discordgo.Button{
									Label:    "かな部分一致検索",
									Style:    discordgo.PrimaryButton,
									CustomID: "search_kana",
								},
								discordgo.Button{
									Label:    "あいまい検索",
									Style:    discordgo.PrimaryButton,
									CustomID: "search_fuzzy",
								},
								discordgo.Button{
									Label:    "かなあいまい検索",
									Style:    discordgo.PrimaryButton,
									CustomID: "search_kana_fuzzy",
								},
								discordgo.Button{
									Label:    "正規表現検索",
									Style:    discordgo.PrimaryButton,
									CustomID: "search_regex",
								},
							},
						},
					},
				},
			})
		}
	case discordgo.InteractionMessageComponent:
		h.logger.Debug("InteractionMessageComponent", zap.String("customID", i.MessageComponentData().CustomID))
		if strings.HasPrefix(i.MessageComponentData().CustomID, "search") {
			var data AutoSearchInfo
			if err := json.Unmarshal([]byte(i.Message.Content[2:strings.Index(i.Message.Content, "\n")]), &data); err != nil {
				h.logger.Error("json.Unmarshal error", zap.Error(err))
				return
			}
			if i.MessageComponentData().CustomID == "search_kana" || i.MessageComponentData().CustomID == "search_kana_fuzzy" {
				data.KanaSearch = true
			}
			if i.MessageComponentData().CustomID == "search_fuzzy" || i.MessageComponentData().CustomID == "search_kana_fuzzy" {
				data.FuzzySearch = true
			} else if i.MessageComponentData().CustomID == "search_regex" {
				data.RegexSearch = true
			}
			jsonStr, err := json.Marshal(data)
			if err != nil {
				h.logger.Error("json.Marshal error", zap.Error(err))
				return
			}
			if err := s.ChannelMessageDelete(i.Message.ChannelID, i.Message.ID); err != nil {
				h.logger.Error("ChannelMessageDelete error", zap.Error(err))
			}
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content:  fmt.Sprintf("> %s\n録画する？", string(jsonStr)),
					CustomID: "search_method",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label:    "録画する",
									Style:    discordgo.PrimaryButton,
									CustomID: "record",
								},
								discordgo.Button{
									Label:    "録画しない",
									Style:    discordgo.SecondaryButton,
									CustomID: "record_not",
								},
							},
						},
					},
				},
			}); err != nil {
				h.logger.Error("InteractionRespond error", zap.Error(err))
				return
			}
		} else if strings.HasPrefix(i.MessageComponentData().CustomID, "record") {
			var data AutoSearchInfo
			if err := json.Unmarshal([]byte(i.Message.Content[2:strings.Index(i.Message.Content, "\n")]), &data); err != nil {
				h.logger.Error("json.Unmarshal error", zap.Error(err))
				return
			}
			if i.MessageComponentData().CustomID == "record" {
				data.Record = true
			}
			var recordStr string
			if data.Record {
				recordStr = "録画する"
			} else {
				recordStr = "録画しない"
			}
			jsonStr, err := json.Marshal(data)
			if err != nil {
				h.logger.Error("json.Marshal error", zap.Error(err))
				return
			}
			if err := s.ChannelMessageDelete(i.Message.ChannelID, i.Message.ID); err != nil {
				h.logger.Error("ChannelMessageDelete error", zap.Error(err))
			}
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("> %s\nスレッド名:%s\nタイトル:%s\nチャンネル:%s\nジャンル:%s\n録画:%s\nこれでよろしいですか？",
						jsonStr,
						data.Name,
						data.Title,
						data.Channel,
						data.Genre,
						recordStr,
					),
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label:    "OK",
									Style:    discordgo.PrimaryButton,
									CustomID: "confirm_ok",
								},
								discordgo.Button{
									Label:    "キャンセル",
									Style:    discordgo.SecondaryButton,
									CustomID: "confirm_cancel",
								},
							},
						},
					},
				},
			}); err != nil {
				h.logger.Error("InteractionRespond error", zap.Error(err))
				return
			}
		} else if strings.HasPrefix(i.MessageComponentData().CustomID, "confirm") {
			var data AutoSearchInfo
			if err := json.Unmarshal([]byte(i.Message.Content[2:strings.Index(i.Message.Content, "\n")]), &data); err != nil {
				h.logger.Error("json.Unmarshal error", zap.Error(err))
				return
			}
			if i.MessageComponentData().CustomID == "confirm_ok" {
				if err := s.ChannelMessageDelete(i.Message.ChannelID, i.Message.ID); err != nil {
					h.logger.Error("ChannelMessageDelete error", zap.Error(err))
				}
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "自動検索スレッドを作成します",
					},
				}); err != nil {
					h.logger.Error("InteractionRespond error", zap.Error(err))
					return
				}
				if err := h.dtv.CreateAutoSearch(i.Interaction.Member.User.ID, data.Name, data.Title, data.Channel, data.Genre, data.KanaSearch, data.FuzzySearch, data.RegexSearch, data.Record); err != nil {
					h.logger.Error("CreateAutoSearchThread error", zap.Error(err))
					return
				}
			} else if i.MessageComponentData().CustomID == "confirm_cancel" {
				// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
				// 	Data: &discordgo.InteractionResponseData{
				// 		Flags:   discordgo.MessageFlagsEphemeral,
				// 		Content: "キャンセルしました",
				// 	},
				// })
				if err := s.ChannelMessageDelete(i.Message.ChannelID, i.Message.ID); err != nil {
					h.logger.Error("ChannelMessageDelete error", zap.Error(err))
				}
			}
		}
	}
}

func (h *DiscordHandler) RegisterCommand() {
	h.session.AddHandler(h.CommandHandler)

	_, err := h.session.ApplicationCommandCreate(h.session.State.User.ID, h.session.State.Guilds[0].ID, &discordgo.ApplicationCommand{
		Name:        "index",
		Description: "視聴ちゃんの検索インデックスを再作成",
	})
	if err != nil {
		h.logger.Error("ApplicationCommandCreate error", zap.Error(err))
	}
	if _, err := h.session.ApplicationCommandCreate(h.session.State.User.ID, h.session.State.Guilds[0].ID, &discordgo.ApplicationCommand{
		Name:        "create",
		Description: "自動検索の新規作成",
	}); err != nil {
		h.logger.Error("ApplicationCommandCreate error", zap.Error(err))
	}
	if _, err := h.session.ApplicationCommandCreate(h.session.State.User.ID, h.session.State.Guilds[0].ID, &discordgo.ApplicationCommand{
		Name:        "delete",
		Description: "自動検索スレッドの削除",
	}); err != nil {
		h.logger.Error("ApplicationCommandCreate error", zap.Error(err))
	}
}
