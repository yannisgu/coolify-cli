package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate_SendsPortsMappings(t *testing.T) {
	path, body, err := runDatabaseCommand(t, "update", "db-uuid", "--ports-mappings", "10.8.0.5:5432:5432,[::1]:5433:5432")
	require.NoError(t, err)
	assert.Equal(t, "/api/v1/databases/db-uuid", path)
	assert.JSONEq(t, `{"ports_mappings":"10.8.0.5:5432:5432,[::1]:5433:5432"}`, body)
}

func TestUpdate_EmptyPortsMappingsIsSentToClear(t *testing.T) {
	path, body, err := runDatabaseCommand(t, "update", "db-uuid", "--ports-mappings", "")
	require.NoError(t, err)
	assert.Equal(t, "/api/v1/databases/db-uuid", path)
	assert.JSONEq(t, `{"ports_mappings":""}`, body)
}

func TestUpdate_PortsMappingsOmittedWhenFlagNotSet(t *testing.T) {
	_, body, err := runDatabaseCommand(t, "update", "db-uuid", "--name", "renamed")
	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"renamed"}`, body)
}

func TestCreate_SendsPortsMappings(t *testing.T) {
	path, body, err := runDatabaseCommand(t, "create", "postgresql",
		"--server-uuid", "srv-1", "--project-uuid", "proj-1", "--environment-name", "production",
		"--ports-mappings", "10.8.0.5:5432:5432")
	require.NoError(t, err)
	assert.Equal(t, "/api/v1/databases/postgresql", path)
	assert.JSONEq(t, `{
		"server_uuid":"srv-1",
		"project_uuid":"proj-1",
		"environment_name":"production",
		"ports_mappings":"10.8.0.5:5432:5432"
	}`, body)
}

func TestCreateSqlite_RejectsPortsMappings(t *testing.T) {
	path, _, err := runDatabaseCommand(t, "create", "sqlite",
		"--server-uuid", "srv-1", "--project-uuid", "proj-1", "--environment-name", "production",
		"--ports-mappings", "5432:5432")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not supported for sqlite")
	assert.Empty(t, path, "no request should be sent")
}
