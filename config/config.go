package config

import "github.com/spf13/viper"

const (
	ShowAllConstraintsKey          = "showAllConstraints"
	UseAllTablesKey                = "useAllTables"
	SelectedTablesKey              = "selectedTables"
	SchemaKey                      = "schema"
	ConnectionStringKey            = "connectionString"
	ConnectionStringSuggestionsKey = "connectionStringSuggestions"
	OutputFileNameKey              = "outputFileName"
	EncloseWithMermaidBackticksKey = "encloseWithMermaidBackticks"
	DebugKey                       = "debug"
	OmitConstraintLabelsKey        = "omitConstraintLabels"
	OmitAttributeKeysKey           = "omitAttributeKeys"
)

type config struct{}

type MermerdConfig interface {
	ShowAllConstraints() bool
	UseAllTables() bool
	Schema() string
	ConnectionString() string
	OutputFileName() string
	ConnectionStringSuggestions() []string
	SelectedTables() []string
	EncloseWithMermaidBackticks() bool
	Debug() bool
	OmitConstraintLabels() bool
	OmitAttributeKeys() bool
}

func NewConfig() MermerdConfig {
	return config{}
}

func (c config) ShowAllConstraints() bool {
	return viper.GetBool(ShowAllConstraintsKey)
}

func (c config) UseAllTables() bool {
	return viper.GetBool(UseAllTablesKey)
}

func (c config) Schema() string {
	return viper.GetString(SchemaKey)
}

func (c config) ConnectionString() string {
	return viper.GetString(ConnectionStringKey)
}

func (c config) OutputFileName() string {
	return viper.GetString(OutputFileNameKey)
}

func (c config) ConnectionStringSuggestions() []string {
	return viper.GetStringSlice(ConnectionStringSuggestionsKey)
}

func (c config) SelectedTables() []string {
	return viper.GetStringSlice(SelectedTablesKey)
}

func (c config) EncloseWithMermaidBackticks() bool {
	return viper.GetBool(EncloseWithMermaidBackticksKey)
}

func (c config) Debug() bool {
	return viper.GetBool(DebugKey)
}

func (c config) OmitConstraintLabels() bool {
	return viper.GetBool(OmitConstraintLabelsKey)
}

func (c config) OmitAttributeKeys() bool {
	return viper.GetBool(OmitAttributeKeysKey)
}
