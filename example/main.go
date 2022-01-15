package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	flyacl "github.com/fahmifan/flycasbin/acl"
	"github.com/fahmifan/flycasbin/acl/db"
	_ "github.com/lib/pq"
)

// define subjects
const (
	Reader flyacl.Role = "reader"
	Editor flyacl.Role = "editor"
)

// define resources
const (
	Story flyacl.Resource = "story"
)

// define actions
const (
	Read   flyacl.Action = "read"
	Write  flyacl.Action = "write"
	Delete flyacl.Action = "delete"
)

// define policies
var policies = []flyacl.Policy{
	{Reader, Story, Read},
	{Editor, Story, Read},
	{Editor, Story, Write},
}

// In this example we define the policies in the code and store it to db. And load it from db
func main() {
	conn, err := sql.Open("postgres", "user=root password=root dbname=flycasbin_acl sslmode=disable")
	panicErr(err)
	queries := db.New(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err = flyacl.StorePolicies(ctx, policies, queries)
	panicErr(err)

	policies, err = flyacl.LoadPolicies(ctx, queries)
	panicErr(err)

	acl, err := flyacl.NewACL(policies)
	panicErr(err)

	acl.Can(Editor, Delete, Story)
	err = acl.Can(Editor, Delete, Story)
	dump(Editor, Delete, Story, err)

	err = acl.Can(Editor, Read, Story)
	dump(Editor, Read, Story, err)

	err = acl.Can(Editor, Write, Story)
	dump(Editor, Write, Story, err)

	err = acl.Can(Reader, Write, Story)
	dump(Reader, Write, Story, err)
}

func dump(role flyacl.Role, act flyacl.Action, rsc flyacl.Resource, err error) {
	if err != nil {
		fmt.Printf("%s cannot %s a %s\n", role, act, rsc)
		return
	}
	fmt.Printf("%s can %s a %s\n", role, act, rsc)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
