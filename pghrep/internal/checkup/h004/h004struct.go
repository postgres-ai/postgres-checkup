package h004

import (
	"gitlab.com/postgres-ai/postgres-checkup/pghrep/internal/checkup"
)

type H004Index struct {
	Num                  int    `json:"num"`
	IndexId              string `json:"index_id"`
	SchemaName           string `json:"schema_name"`
	TableName            string `json:"table_name"`
	TableSizeBytes       int64  `json:"table_size_bytes"`
	IndexName            string `json:"index_name"`
	AccessMethod         string `json:"access_method"`
	Reason               string `json:"reason"`
	MainIndexDef         string `json:"main_index_def"`
	MainIndexSize        string `json:"main_index_size"`
	IndexDef             string `json:"index_def"`
	IndexSizeBytes       int64  `json:"index_size_bytes"`
	IndexUsage           int64  `json:"index_usage"`
	FormatedIndexName    string `json:"formated_index_name"`
	FormatedSchemaName   string `json:"formated_schema_name"`
	FormatedTableName    string `json:"formated_table_name"`
	FormatedRelationName string `json:"formated_relation_name"`
	SupportsFk           bool   `json:"supports_fk"`
}

type H004Indexes map[string]H004Index

type H004IndexesTotal struct {
	IndexSizeBytesSum int64 `json:"index_size_bytes_sum"`
	TableSizeBytesSum int64 `json:"table_size_bytes_sum"`
}

type DatabaseStat struct {
	StatsReset        string `json:"stats_reset"`
	StatsAge          string `json:"stats_age"`
	Days              int64  `json:"days"`
	DatabaseSizeBytes int64  `json:"database_size_bytes"`
}

type H004ReportHostResultData struct {
	RedundantIndexes      H004Indexes      `json:"redundant_indexes"`
	RedundantIndexesTotal H004IndexesTotal `json:"redundant_indexes_total"`
	Do                    []string         `json:"do"`
	UnDo                  []string         `json:"undo"`
	DatabaseStat          DatabaseStat     `json:"database_stat"`
	MinIndexSizeBytes     int64            `json:"min_index_size_bytes"`
}

type H004ReportHostResult struct {
	Data      H004ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type H004ReportHostsResults map[string]H004ReportHostResult

type H004Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       H004ReportHostsResults  `json:"results"`
}
