package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/slack-go/slack"
)

func main() {
	err := doMassRenameChannel(context.Background(), map[string]string{
		"times_motemen_z": "times_motemen",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func doMassRenameChannel(ctx context.Context, nameMapping map[string]string) error {
	api := slack.New(os.Getenv("SLACK_TOKEN"))

	idToNewName := map[string]string{}

	params := &slack.GetConversationsParameters{
		Limit: 1000,
	}

	for {
		log.Println(params)

		chs, cursor, err := api.GetConversationsContext(ctx, params)
		if err != nil {
			return fmt.Errorf("GetConversationsContext: %w", err)
		}

		log.Println(len(chs))
		if len(chs) == 0 {
			break
		}
		if cursor == "" {
			break
		}

		for _, ch := range chs {
			if newName, ok := nameMapping[ch.Name]; ok {
				idToNewName[ch.ID] = newName
			}
		}

		params.Cursor = cursor

		time.Sleep(5 * time.Millisecond)
	}

	log.Println(idToNewName)

	for fromId, newName := range idToNewName {
		_, err := api.RenameConversationContext(
			ctx,
			fromId,
			newName,
		)

		return err
	}

	return nil
}

func doRenameChannel(ctx context.Context, from, to string) error {
	api := slack.New(os.Getenv("SLACK_TOKEN"))

	var channelID string

	params := &slack.GetConversationsParameters{}

getConversasions:
	for {
		chs, cursor, err := api.GetConversationsContext(ctx, params)
		if err != nil {
			return fmt.Errorf("GetConversationsContext: %w", err)
		}

		if len(chs) == 0 {
			break
		}

		for _, ch := range chs {
			if ch.Name == from {
				channelID = ch.ID
				break getConversasions
			}
		}

		params.Cursor = cursor
	}

	if channelID == "" {
		return fmt.Errorf("channel %q not found", from)
	}

	log.Println(channelID)

	_, err := api.RenameConversationContext(
		ctx,
		channelID,
		to,
	)

	return err
}
