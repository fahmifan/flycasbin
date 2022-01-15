CREATE TABLE IF NOT EXISTS "resources" (
    id TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS "actions" (
    id TEXT PRIMARY KEY

);

CREATE TABLE IF NOT EXISTS "roles" (
    id TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS "acl_policies" (
    "role" TEXT NOT NULL REFERENCES "roles"(id),
    "resource" TEXT NOT NULL REFERENCES "resources"(id),
    "action" TEXT NOT NULL REFERENCES "actions"(id)
);