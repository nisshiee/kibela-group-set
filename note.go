package main

import (
	"context"

	"github.com/shurcooL/graphql"
)

type Note struct {
	ID         graphql.String
	Title      graphql.String
	Content    graphql.String
	Coediting  graphql.Boolean
	FolderName graphql.String
	Groups     []struct {
		ID graphql.String
	}
}

func FindNote(client *graphql.Client, id string) (*Note, error) {
	var res struct {
		Note Note `graphql:"note(id: $id)"`
	}
	variables := map[string]interface{}{
		"id": graphql.ID(id),
	}
	if err := client.Query(context.Background(), &res, variables); err != nil {
		return nil, err
	}
	return &res.Note, nil
}

func (n *Note) HasGroup(groupID string) bool {
	for _, group := range n.Groups {
		if string(group.ID) == groupID {
			return true
		}
	}
	return false
}

func (n *Note) AddGroup(client *graphql.Client, groupID string) error {
	baseNote := n.toInput()

	n.Groups = append(n.Groups, struct{ ID graphql.String }{ID: graphql.String(groupID)})
	newNote := n.toInput()

	var res struct {
		UpdateNote struct {
			ClientMutationId graphql.String
		} `graphql:"updateNote(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": UpdateNoteInput{
			ID:       n.ID,
			BaseNote: baseNote,
			NewNote:  newNote,
			Draft:    graphql.Boolean(false),
		},
	}
	if err := client.Mutate(context.Background(), &res, variables); err != nil {
		return err
	}

	return nil
}

type UpdateNoteInput struct {
	ID       graphql.ID      `json:"id"`
	NewNote  NoteInput       `json:"newNote"`
	BaseNote NoteInput       `json:"baseNote"`
	Draft    graphql.Boolean `json:"draft"`
}

type NoteInput struct {
	Title      graphql.String  `json:"title"`
	Content    graphql.String  `json:"content"`
	Coediting  graphql.Boolean `json:"coediting"`
	FolderName graphql.String  `json:"folderName"`
	GroupIds   []graphql.ID    `json:"groupIds"`
}

func (n *Note) toInput() NoteInput {
	groupIds := make([]graphql.ID, len(n.Groups))
	for _, g := range n.Groups {
		groupIds = append(groupIds, g.ID)
	}

	return NoteInput{
		Title:      n.Title,
		Content:    n.Content,
		Coediting:  n.Coediting,
		FolderName: n.FolderName,
		GroupIds:   groupIds,
	}
}
