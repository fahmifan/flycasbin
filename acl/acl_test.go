package acl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestACL(t *testing.T) {
	// define subjects
	const (
		Reader Role = "reader"
		Editor Role = "editor"
	)

	// define resources
	const (
		Story Resource = "story"
	)

	// define actions
	const (
		Read   Action = "read"
		Write  Action = "write"
		Delete Action = "delete"
	)

	// define policies
	policies := []Policy{
		{Reader, Story, Read},
		{Editor, Story, Read},
		{Editor, Story, Write},
	}

	acl, err := NewACL(policies)
	require.NoError(t, err)

	err = acl.Can(Editor, Read, Story)
	require.NoError(t, err)
}
