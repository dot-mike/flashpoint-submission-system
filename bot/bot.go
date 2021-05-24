package bot

import (
	"fmt"
	"github.com/Dri0m/flashpoint-submission-system/types"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Bot struct {
	Session            *discordgo.Session
	FlashpointServerID string
	L                  *logrus.Logger
}

// ConnectBot connects bot or panics
func ConnectBot(l *logrus.Logger, token string) *discordgo.Session {
	l.Infoln("connecting discord bot...")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		l.Fatal(err)
	}
	return dg
}

// GetFlashpointRoleIDsForUser returns user role IDs
func (b *Bot) GetFlashpointRoleIDsForUser(uid int64) ([]string, error) {
	b.L.WithField("uid", uid).Info("getting flashpoint role ID for user")
	member, err := b.Session.GuildMember(b.FlashpointServerID, fmt.Sprint(uid))
	if err != nil {
		return nil, err
	}

	return member.Roles, nil
}

// GetFlashpointRoles returns list of flashpoint server roles
func (b *Bot) GetFlashpointRoles() ([]types.DiscordRole, error) {
	b.L.Info("getting flashpoint roles")
	roles, err := b.Session.GuildRoles(b.FlashpointServerID)
	if err != nil {
		return nil, err
	}

	result, err := formatDiscordgoRoles(roles)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func formatDiscordgoRoles(roles []*discordgo.Role) ([]types.DiscordRole, error) {
	formattedRoles := make([]types.DiscordRole, 0, len(roles))
	for _, role := range roles {
		id, err := strconv.ParseInt(role.ID, 10, 64)
		if err != nil {
			return nil, err
		}
		formattedRoles = append(formattedRoles, types.DiscordRole{ID: id, Name: role.Name, Color: fmt.Sprintf("#%06x", role.Color)})
	}
	return formattedRoles, nil
}

// IsUserAuthorized contacts discord api to check if user has sufficient roles to use this site
func (b *Bot) IsUserAuthorized(uid int64) (bool, error) {
	userRoles := make([]types.DiscordRole, 0)

	roles, err := b.GetFlashpointRoles()
	if err != nil {
		return false, err
	}

	roleIDs, err := b.GetFlashpointRoleIDsForUser(uid)
	if err != nil {
		return false, err
	}

	for _, roleID := range roleIDs {
		for _, role := range roles {
			id, err := strconv.ParseInt(roleID, 10, 64)
			if err != nil {
				return false, err
			}
			if role.ID == id {
				userRoles = append(userRoles, role)
			}
		}
	}

	authorizedRoles := []string{"Administrator", "Moderator", "Curator", "Tester", "Mechanic", "Hunter", "Hacker"}

	isAuthorized := false
	for _, role := range userRoles {
		for _, authorizedRole := range authorizedRoles {
			if role.Name == authorizedRole {
				isAuthorized = true
				break
			}
		}
	}

	return isAuthorized, nil
}