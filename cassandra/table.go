package cassandra

import "github.com/scylladb/gocqlx/v2/table"

var sampleTableMetadata = table.Metadata{
	Name: "sample_table",
	Columns: []string{
		"column1",
		"column2",
		"column3",
	},
	PartKey: []string{"column1"},
	SortKey: []string{"column2"},
}

var sampleTable = table.New(sampleTableMetadata)

var sampleTableDdl = `CREATE TABLE IF NOT EXISTS sample_table (
    column1 text,
    column2 text,
    column3 text,
    PRIMARY KEY (column1, column2)
    ) WITH CLUSTERING ORDER BY (column2 ASC);
`

type SampleTable struct {
	Column1 string
	Column2 string
	Column3 string
}

var sampleTableSelectStmt, sampleTableSelectNames = sampleTable.Select()

func (c *Client) Select(column1, column2 string) (*SampleTable, error) {
	sess, err := c.ConnectX()
	if err != nil {
		return nil, err
	}

	var result SampleTable
	if err := sess.Query(sampleTableSelectStmt, sampleTableSelectNames).BindMap(map[string]any{
		"column1": column1,
		"column2": column2,
	}).SelectRelease(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

var sampleTableInsertStmt, sampleTableInsertNames = sampleTable.Insert()

func (c *Client) Insert(data *SampleTable) error {
	sess, err := c.ConnectX()
	if err != nil {
		return err
	}

	if err := sess.Query(sampleTableInsertStmt, sampleTableInsertNames).BindStruct(data).ExecRelease(); err != nil {
		return err
	}

	return nil
}

var sampleTableUpdateStmt, sampleTableUpdateNames = sampleTable.Update()

func (c *Client) Update(data *SampleTable) error {
	sess, err := c.ConnectX()
	if err != nil {
		return err
	}

	if err := sess.Query(sampleTableUpdateStmt, sampleTableUpdateNames).BindStruct(data).ExecRelease(); err != nil {
		return err
	}

	return nil
}

var sampleTableDeleteStmt, sampleTableDeleteNames = sampleTable.Delete()

func (c *Client) Delete(column1, column2 string) error {
	sess, err := c.ConnectX()
	if err != nil {
		return err
	}

	if err := sess.Query(sampleTableDeleteStmt, sampleTableDeleteNames).BindMap(map[string]any{
		"column1": column1,
		"column2": column2,
	}).ExecRelease(); err != nil {
		return err
	}

	return nil
}
