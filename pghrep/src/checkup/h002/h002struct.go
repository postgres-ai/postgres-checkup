package h002

import checkup ".."

type H002Index struct {
	Num                  int     `json:"num"`
	Reason               string  `json:"reason"`
	IndexId              string  `json:"index_id"`
	SchemaName           string  `json:"schema_name"`
	TableName            string  `json:"table_name"`
	IndexName            string  `json:"index_name"`
	IdxScan              int64   `json:"idx_scan"`
	AllScans             int64   `json:"all_scans"`
	IndexScanPct         float64 `json:"index_scan_pct"`
	Writes               int64   `json:"writes"`
	ScansPerWrite        float64 `json:"scans_per_write"`
	IndexSizeBytes       int64   `json:"index_size_bytes"`
	TableSizeBytes       int64   `json:"table_size_bytes"`
	Relpages             int64   `json:"relpages"`
	IdxIsBtree           bool    `json:"idx_is_btree"`
	IndexDef             string  `json:"index_def"`
	FormatedIndexName    string  `json:"formated_index_name"`
	FormatedSchemaName   string  `json:"formated_schema_name"`
	FormatedTableName    string  `json:"formated_table_name"`
	FormatedRelationName string  `json:"formated_relation_name"`
	Opclasses            string  `json:"opclasses"`
	SupportsFk           bool    `json:"supports_fk"`
	Grp                  int64   `json:"grp"`
}

type H002Indexes map[string]H002Index

type H002IndexesTotal struct {
	IndexSizeBytesSum int64 `json:"index_size_bytes_sum"`
	TableSizeBytesSum int64 `json:"table_size_bytes_sum"`
}

type DatabaseStat struct {
	StatsReset        string `json:"stats_reset"`
	StatsAge          string `json:"stats_age"`
	Days              int64  `json:"days"`
	DatabaseSizeBytes int64  `json:"database_size_bytes"`
}

type H002ReportHostResultData struct {
	NeverUsedIndexes       H002Indexes      `json:"never_used_indexes"`
	NeverUsedIndexesTotal  H002IndexesTotal `json:"never_used_indexes_total"`
	RarelyUsedIsndexes     H002Indexes      `json:"rarely_used_indexes"`
	RarelyUsedIndexesTotal H002IndexesTotal `json:"rarely_used_indexes_total"`
	Do                     []string         `json:"do"`
	UnDo                   []string         `json:"undo"`
	DatabaseStat           DatabaseStat     `json:"database_stat"`
	MinIndexSizeBytes      int64            `json:"min_index_size_bytes"`
}

type H002ReportHostResult struct {
	Data      H002ReportHostResultData `json:"data"`
	NodesJson checkup.ReportLastNodes  `json:"nodes.json"`
}

type H002ReportHostsResults map[string]H002ReportHostResult

type H002Report struct {
	Project       string                  `json:"project"`
	Name          string                  `json:"name"`
	CheckId       string                  `json:"checkId"`
	Timestamptz   string                  `json:"timestamptz"`
	Database      string                  `json:"database"`
	Dependencies  map[string]interface{}  `json:"dependencies"`
	LastNodesJson checkup.ReportLastNodes `json:"last_nodes_json"`
	Results       H002ReportHostsResults  `json:"results"`
}
