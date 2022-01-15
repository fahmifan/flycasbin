# Fly Casbin

## Feature
- ACL with casbin
- Add type to Role, Resource and Action
- Define Policies in code, then store & load them from db

## Usage
```go
var policies = []flyacl.Policy{
	{Reader, Story, Read},
	{Editor, Story, Read},
	{Editor, Story, Write},
}

acl, err := flyacl.NewACL(policies)
checkErr(err)

// error
err := acl.Can(Editor, Delete, Story)
checkErr(err)

// ok
err := acl.Can(Reader, Read, Story)
checkErr(err)
```