package migrations

import "embed"

var(

	//go:embed *sql
	MigrationsFs embed.FS
)