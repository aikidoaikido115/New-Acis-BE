//go:build ignore

package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&entities.Role{},
		&entities.User{},
		&entities.Staff{},
		&entities.StaffsFiles{},
		&entities.OTP{},
		&entities.TempToken{},
		&entities.AuditLogs{},
		&entities.Room{},
		&entities.Resident{},
		&entities.ResidentLabels{},
		&entities.IntakeLabels{},
		&entities.VitalSign{},
		&entities.LaboratoryValue{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
