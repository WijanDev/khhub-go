package db

import "embed"

// Migrations holds golang-migrate SQL files (*.up.sql / *.down.sql).
//
//go:embed migrations/*.sql
var Migrations embed.FS
