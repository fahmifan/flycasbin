package flycasbin

import (
	"errors"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

var once sync.Once
var enforcer *casbin.Enforcer
var aclModel = `
# Request definition
[request_definition]
r = role, resource, action

# Policy definition
[policy_definition]
p = role, resource, action

# Policy effect
[policy_effect]
e = some(where (p.eft == allow))

# Matchers
[matchers]
m = r.role == p.role && r.resource == p.resource && r.act == p.action
`

func init() {
	once.Do(func() {
		acl, err := model.NewModelFromString(aclModel)
		panicErr(err)

		enforcer, err = casbin.NewEnforcer(acl)
		panicErr(err)
	})
}

// errors ..
var (
	ErrPermissionDenied = errors.New("permission denied")
)

// Role ..
type Role string

func (s Role) Can(act Action, rsc Resource) error {
	ok, err := enforcer.Enforce(string(s), string(rsc), string(act))
	if err != nil {
		return err
	}
	if !ok {
		return ErrPermissionDenied
	}

	return nil
}

// types ..
type (
	Resource string
	Action   string

	ACL struct {
		Subject  Role
		Resource Resource
		Action   Action
	}
)

func InitPolicies(acls []ACL) {
	for _, acl := range acls {
		enforcer.AddPolicy(
			string(acl.Subject),
			string(acl.Resource),
			string(acl.Action),
		)
	}
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
