package database

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"

	"github.com/stretchr/testify/assert"
)

type columnTestResult struct {
	Name      string
	isPrimary bool
	isForeign bool
}

func TestDatabaseIntegrations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	connectorFactory := NewConnectorFactory()
	testCases := []struct {
		dbType           DbType
		connectionString string
		schema           string
	}{
		{
			dbType:           Postgres,
			connectionString: "postgresql://user:password@localhost:5432/mermerd_test",
			schema:           "public",
		},
		{
			dbType:           MySql,
			connectionString: "mysql://user:password@tcp(127.0.0.1:3306)/mermerd_test",
			schema:           "mermerd_test",
		},
		{
			dbType:           MsSql,
			connectionString: "sqlserver://sa:securePassword1!@localhost:1433?database=mermerd_test",
			schema:           "dbo",
		},
	}

	for _, testCase := range testCases {
		connector, _ := connectorFactory.NewConnector(testCase.connectionString)

		getConnectionAndConnect := func(t *testing.T) Connector {
			err := connector.Connect()
			if err != nil {
				logrus.Error(err)
				t.FailNow()
			}
			return connector
		}

		t.Run(testCase.dbType.String(), func(t *testing.T) {
			t.Run("Connect", func(t *testing.T) {
				// Arrange
				connector := connector

				// Act
				err := connector.Connect()

				// Assert
				assert.Nil(t, err)
			})

			t.Run("GetSchemas", func(t *testing.T) {
				// Arrange
				connector := getConnectionAndConnect(t)

				// Act
				schemas, err := connector.GetSchemas()

				// Assert
				assert.Nil(t, err)
				assert.NotEmpty(t, schemas)
			})

			t.Run("GetTables", func(t *testing.T) {
				// Arrange
				connector := getConnectionAndConnect(t)
				schema := testCase.schema

				// Act
				tables, err := connector.GetTables(schema)

				// Assert
				expectedResult := []string{
					"article",
					"article_detail",
					"article_comment",
					"label",
					"article_label",
					"test_1_a",
					"test_1_b",
				}
				assert.Nil(t, err)
				assert.ElementsMatch(t, expectedResult, tables)
			})

			t.Run("GetColumns", func(t *testing.T) {
				connector := getConnectionAndConnect(t)
				testCases := []struct {
					tableName       string
					expectedColumns []columnTestResult
				}{
					{tableName: "article", expectedColumns: []columnTestResult{
						{Name: "id", isPrimary: true, isForeign: false},
						{Name: "title", isPrimary: false, isForeign: false},
					}},
					{tableName: "article_detail", expectedColumns: []columnTestResult{
						{Name: "id", isPrimary: true, isForeign: true},
						{Name: "created_at", isPrimary: false, isForeign: false},
					}},
					{tableName: "article_comment", expectedColumns: []columnTestResult{
						{Name: "id", isPrimary: true, isForeign: false},
						{Name: "article_id", isPrimary: false, isForeign: true},
						{Name: "comment", isPrimary: false, isForeign: false},
					}},
					{tableName: "label", expectedColumns: []columnTestResult{
						{Name: "id", isPrimary: true, isForeign: false},
						{Name: "label", isPrimary: false, isForeign: false},
					}},
					{tableName: "article_label", expectedColumns: []columnTestResult{
						{Name: "article_id", isPrimary: true, isForeign: true},
						{Name: "label_id", isPrimary: true, isForeign: true},
					}},
					{tableName: "test_1_a", expectedColumns: []columnTestResult{
						{Name: "id", isPrimary: true, isForeign: false},
						{Name: "xid", isPrimary: true, isForeign: false},
					}},
					{tableName: "test_1_b", expectedColumns: []columnTestResult{
						{Name: "aid", isPrimary: true, isForeign: true},
						{Name: "bid", isPrimary: true, isForeign: true},
					}},
				}

				for index, testCase := range testCases {
					t.Run(fmt.Sprintf("run #%d", index), func(t *testing.T) {
						// Arrange
						tableName := testCase.tableName
						var columnResult []columnTestResult

						// Act
						columns, err := connector.GetColumns(tableName)

						// Assert
						for _, column := range columns {
							columnResult = append(columnResult, columnTestResult{
								Name:      column.Name,
								isPrimary: column.IsPrimary,
								isForeign: column.IsForeign,
							})
						}

						assert.Nil(t, err)
						assert.ElementsMatch(t, testCase.expectedColumns, columnResult)
					})
				}
			})

			t.Run("GetConstraints", func(t *testing.T) {
				connector := getConnectionAndConnect(t)

				t.Run("One-to-one relation", func(t *testing.T) {
					// Arrange
					tableName := "article_detail"

					// Act
					constraintResults, err := connector.GetConstraints(tableName)

					// Assert
					assert.Nil(t, err)
					assert.Len(t, constraintResults, 1)
					constraint := constraintResults[0]
					assert.True(t, constraint.IsPrimary)
					assert.False(t, constraint.HasMultiplePK)
				})

				t.Run("Many-to-one relation #1", func(t *testing.T) {
					// Arrange
					tableName := "article_comment"

					// Act
					constraintResults, err := connector.GetConstraints(tableName)

					// Assert
					assert.Nil(t, err)
					assert.Len(t, constraintResults, 1)
					constraint := constraintResults[0]
					assert.False(t, constraint.IsPrimary)
					assert.False(t, constraint.HasMultiplePK)
				})

				t.Run("Many-to-one relation #2", func(t *testing.T) {
					// Arrange
					pkTableName := "article"
					fkTableName := "article_label"

					// Act
					constraintResults, err := connector.GetConstraints(pkTableName)

					// Assert
					assert.Nil(t, err)
					var constraint *ConstraintResult
					for _, item := range constraintResults {
						if item.FkTable == fkTableName {
							constraint = &item
							break
						}
					}
					assert.NotNil(t, constraint)
					assert.True(t, constraint.IsPrimary)
					assert.True(t, constraint.HasMultiplePK)
				})

				// Multiple primary keys (https://github.com/KarnerTh/mermerd/issues/8)
				t.Run("Test 1 (Issue #8)", func(t *testing.T) {
					// Arrange
					pkTableName := "test_1_b"

					// Act
					constraintResults, err := connector.GetConstraints(pkTableName)

					// Assert
					assert.Nil(t, err)
					assert.NotNil(t, constraintResults)
					assert.Len(t, constraintResults, 2)
					assert.True(t, constraintResults[0].IsPrimary)
					assert.True(t, constraintResults[0].HasMultiplePK)
					assert.Equal(t, constraintResults[0].ColumnName, "aid")
					assert.True(t, constraintResults[1].IsPrimary)
					assert.True(t, constraintResults[1].HasMultiplePK)
					assert.Equal(t, constraintResults[1].ColumnName, "bid")
				})
			})
		})
	}
}
