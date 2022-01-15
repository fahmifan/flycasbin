package flycasbin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestACL(t *testing.T) {
	// subjects ..
	const (
		Member Role = "MEMBER"
		Editor Role = "EDITOR"
	)

	// resources ..
	const (
		Story Resource = "story"
	)

	// actions ..
	const (
		Read   Action = "read"
		Write  Action = "write"
		Delete Action = "delete"
	)

	policies := []ACL{
		{Member, Story, Read},
		{Editor, Story, Read},
		{Editor, Story, Write},
	}

	InitPolicies(policies)
	err := Editor.Can(Read, Story)
	require.NoError(t, err)
}
