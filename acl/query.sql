-- name: CreateRole :exec
INSERT INTO roles (id) VALUES (@role) ON CONFLICT DO NOTHING;

-- name: CreateResource :exec
INSERT INTO resources (id) VALUES (@resource) ON CONFLICT DO NOTHING;

-- name: CreateAction :exec
INSERT INTO actions (id) VALUES (@action) ON CONFLICT DO NOTHING;

-- name: CreateACLPolicies :exec
INSERT INTO acl_policies("role", "resource", "action") VALUES (@role, @resource, @action) ON CONFLICT DO NOTHING;

-- name: FindRoles :many
SELECT * FROM roles;

-- name: FindResources :many
SELECT * FROM resources;

-- name: FindActions :many
SELECT * FROM resources;

-- name: FindACLPolicies :many
SELECT * FROM acl_policies;
