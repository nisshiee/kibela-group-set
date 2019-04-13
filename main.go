package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

var (
	endpoint, token, targetFolderID, targetGroupID string
)

func init() {
	team := os.Getenv("KIBELA_TEAM")
	if team == "" {
		panic("KIBELA_TEAM empty")
	}
	endpoint = fmt.Sprintf("https://%s.kibe.la/api/v1", team)

	token = os.Getenv("KIBELA_TOKEN")
	if token == "" {
		panic("KIBELA_TOKEN empty")
	}

	targetFolderID = os.Getenv("KIBELA_TARGET_FOLDER_ID")
	if targetFolderID == "" {
		panic("KIBELA_TARGET_FOLDER_ID empty")
	}

	targetGroupID = os.Getenv("KIBELA_TARGET_GROUP_ID")
	if targetGroupID == "" {
		panic("KIBELA_TARGET_GROUP_ID empty")
	}
}

func main() {
	var client *graphql.Client
	{
		cli := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
			TokenType:   "Bearer",
		}))
		client = graphql.NewClient(endpoint, cli)
	}

	count := 0
	if err := notesRecursiveEach(client, targetFolderID, func(summary *NoteSummary) error {
		if summary.hasGroup(targetGroupID) {
			return nil
		}

		note, err := FindNote(client, string(summary.ID))
		if err != nil {
			return err
		}
		if note.HasGroup(targetGroupID) {
			return nil
		}

		if err := note.AddGroup(client, targetGroupID); err != nil {
			return err
		}
		count++
		return nil
	}); err != nil {
		panic(err)
	}

	fmt.Printf("Updated Note Count: %d\n", count)
}

type NoteSummary struct {
	ID     graphql.String
	Groups []struct {
		ID graphql.String
	}
}

func (n *NoteSummary) hasGroup(groupID string) bool {
	for _, g := range n.Groups {
		if string(g.ID) == groupID {
			return true
		}
	}
	return false
}

type Handler func(*NoteSummary) error

func notesRecursiveEach(client *graphql.Client, rootFolderID string, handler Handler) error {
	if err := noteLoop(client, rootFolderID, "", handler); err != nil {
		return err
	}
	if err := folderLoop(client, rootFolderID, "", handler); err != nil {
		return err
	}
	return nil
}

func noteLoop(client *graphql.Client, folderID, after string, handler Handler) error {
	var res struct {
		Folder struct {
			Notes struct {
				Nodes    []NoteSummary
				PageInfo struct {
					EndCursor   graphql.String
					HasNextPage graphql.Boolean
				}
			} `graphql:"notes(first: 100, after: $after, onlyCoediting: true)"`
		} `graphql:"folder(id: $folderID)"`
	}
	variables := map[string]interface{}{
		"folderID": graphql.ID(folderID),
		"after":    graphql.String(after),
	}

	if err := client.Query(context.Background(), &res, variables); err != nil {
		return err
	}

	for _, note := range res.Folder.Notes.Nodes {
		if err := handler(&note); err != nil {
			return err
		}
	}

	if res.Folder.Notes.PageInfo.HasNextPage {
		if err := noteLoop(
			client,
			folderID,
			string(res.Folder.Notes.PageInfo.EndCursor),
			handler,
		); err != nil {
			return err
		}
	}

	return nil
}

func folderLoop(client *graphql.Client, folderID, after string, handler Handler) error {
	var res struct {
		Folder struct {
			Folders struct {
				Nodes []struct {
					ID graphql.String
				}
				PageInfo struct {
					EndCursor   graphql.String
					HasNextPage graphql.Boolean
				}
			} `graphql:"folders(first: 100, after: $after, active: true)"`
		} `graphql:"folder(id: $folderID)"`
	}
	variables := map[string]interface{}{
		"folderID": graphql.ID(folderID),
		"after":    graphql.String(after),
	}
	if err := client.Query(context.Background(), &res, variables); err != nil {
		return err
	}

	for _, folder := range res.Folder.Folders.Nodes {
		if err := notesRecursiveEach(client, string(folder.ID), handler); err != nil {
			return err
		}
	}

	if res.Folder.Folders.PageInfo.HasNextPage {
		if err := folderLoop(
			client,
			folderID,
			string(res.Folder.Folders.PageInfo.EndCursor),
			handler,
		); err != nil {
			return err
		}
	}

	return nil
}
