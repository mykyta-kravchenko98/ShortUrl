package migration

import "embed"

// Files embeds every *.sql migration in this directory directly into
// whatever binary imports this package - so the migration Job image
// (cmd/migrate) always ships exactly the migrations that exist in this
// commit, with no separate copy to keep in sync (see
// shorturl-gitops/argocd - the old approach hand-copied these into a
// ConfigMap in the gitops repo, which drifted).
//
//go:embed *.sql
var Files embed.FS
