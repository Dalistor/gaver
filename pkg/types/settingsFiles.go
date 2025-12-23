package types

type GaverModuleFile struct {
	Type                string
	ProjectName         string
	ProjectVersion      string
	ProjectModules      []string
	ProjectDatabaseType string
	MigrationTag        int
}
