package main

import (
	"fmt"
	"os"

	"server/common/models"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	// 1. INJECT REQUIRED EXTENSIONS
	// This tells Atlas that these extensions are part of your schema.
	// Atlas will execute these on the Dev DB before creating tables.
	fmt.Println(`CREATE EXTENSION IF NOT EXISTS "postgis";`)
	fmt.Println(`CREATE EXTENSION IF NOT EXISTS "pg_trgm";`)

	// 2. Load your GORM models
	stmts, err := gormschema.New("postgres").Load(
		&models.User{},
		&models.Team{},
		&models.Facility{},
		&models.Challenge{},
		&models.Sport{},
		&models.Invitation{},
		&models.Location{},
		&models.Notification{},
		&models.UserSettings{},
		&models.Message{},
		&models.Conversation{},
		&models.ConversationParticipant{},
		&models.EmergencyInfo{},
		&models.Report{},
		&models.EulaVersion{},
		&models.EulaAcceptance{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	// 3. Print the rest of the schema
	fmt.Print(stmts)
}
