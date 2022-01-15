package acl

import (
	"context"
	"errors"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/fahmifan/flycasbin/acl/db"
)

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
m = r.role == p.role && r.resource == p.resource && r.action == p.action
`

// errors ..
var (
	ErrPermissionDenied = errors.New("permission denied")
)

// types ..
type (
	Role     string
	Resource string
	Action   string

	Policy struct {
		Role     Role
		Resource Resource
		Action   Action
	}
)

// NewACL create new ACL with the given policies
func NewACL(policies []Policy) (*ACL, error) {
	enf, err := newEnforcer()
	if err != nil {
		return nil, err
	}

	acl := &ACL{
		Policies: policies,
		enforcer: enf,
	}

	err = acl.initPolicies(policies)
	if err != nil {
		return nil, err
	}

	return acl, nil
}

// LoadPolicies load policies from db
func LoadPolicies(ctx context.Context, queries *db.Queries) ([]Policy, error) {
	dbPolices, err := queries.FindACLPolicies(ctx)
	if err != nil {
		return nil, err
	}

	var policies []Policy
	for _, pol := range dbPolices {
		policies = append(policies, Policy{
			Role(pol.Role), Resource(pol.Resource), Action(pol.Action),
		})
	}

	return policies, nil
}

// StorePolicies store policies to db
func StorePolicies(ctx context.Context, policies []Policy, queries *db.Queries) error {
	for _, pol := range policies {
		err := queries.CreateRole(ctx, string(pol.Role))
		if err != nil {
			return err
		}

		err = queries.CreateAction(ctx, string(pol.Action))
		if err != nil {
			return err
		}

		err = queries.CreateResource(ctx, string(pol.Resource))
		if err != nil {
			return err
		}

		err = queries.CreateACLPolicies(ctx, db.CreateACLPoliciesParams{
			Role:     string(pol.Role),
			Resource: string(pol.Resource),
			Action:   string(pol.Action),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type ACL struct {
	Policies []Policy
	enforcer *casbin.Enforcer
}

func (a *ACL) initPolicies(policies []Policy) error {
	for _, pol := range policies {
		_, err := a.enforcer.AddPolicy(
			string(pol.Role),
			string(pol.Resource),
			string(pol.Action),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Can check if the role can do action to a resource
func (a *ACL) Can(role Role, action Action, resource Resource) error {
	ok, err := a.enforcer.Enforce(string(role), string(resource), string(action))
	if err != nil {
		return err
	}
	if !ok {
		return ErrPermissionDenied
	}
	return nil
}

var aclCasbinModel model.Model
var casbinModel sync.Once

func newEnforcer() (*casbin.Enforcer, error) {
	var err error
	casbinModel.Do(func() {
		aclCasbinModel, err = model.NewModelFromString(aclModel)
	})
	if err != nil {
		return nil, err
	}

	enf, err := casbin.NewEnforcer(aclCasbinModel)
	if err != nil {
		return nil, err
	}
	return enf, nil
}
