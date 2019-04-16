Pghrep used by `checkup` during database analysis to create markdown reports from JSON reports.

To build report generator use:

`cd /pghrep`
`make all`

To generate report use:

`./bin/pghrep --checkdata=/path_to_json_repots_storage/A002_pgversion.json --outdir=/path_to_md_reports_storage/ `

also, to enable debug mode, use `--debug 1`
