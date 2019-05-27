package a002

const VERSION_SOURCE_URL string = "https://git.postgresql.org/gitweb/?p=postgresql.git;a=tags"

const MSG_WRONG_VERSION_CONCLUSION string = "[P1] Unknown PostgreSQL version %s on %s."
const MSG_WRONG_VERSION_RECOMMENDATION string = "[P1] Check PostgreSQL version on %s."
const MSG_NOT_SUPPORTED_VERSION_CONCLUSION string = "[P1] Postgres major version being used is %s and it is " +
	"NOT supported by Postgres community and PGDG (supported ended %s). This is a major issue. New bugs and security " +
	"issues will not be fixed by community and PGDG. You are on your own! Read more: " +
	"[Versioning Policy](https://www.postgresql.org/support/versioning/)."
const MSG_NOT_SUPPORTED_VERSION_RECOMMENDATION string = "[P1] Please upgrade Postgres version %s to one of the " +
	"versions supported by the community and PGDG. To minimize downtime, consider using pg_upgrade or one " +
	"of solutions for logical replication."
const MSG_LAST_YEAR_SUPPORTED_VERSION_CONCLUSION string = "[P2] Postgres community and PGDG will stop supporting version %s" +
	" within the next 12 months (end of life is scheduled %s). After that, you will be on your own!"
const MSG_SUPPORTED_VERSION_CONCLUSION string = "Postgres major version being used is %s and it is " +
	"currently supported by Postgres community and PGDG (end of life is scheduled %s). It means that in case " +
	"of bugs and security issues, updates (new minor versions) with fixes will be released and available for use." +
	" Read more: [Versioning Policy](https://www.postgresql.org/support/versioning/)."
const MSG_LAST_MINOR_VERSION_CONCLUSION string = "%s is the most up-to-date Postgres minor version in the branch %s."
const MSG_NOT_LAST_MINOR_VERSION_CONCLUSION_1 string = "[P2] The minor version being used (%s) are not up-to-date (%s)."
const MSG_NOT_LAST_MINOR_VERSION_CONCLUSION_N string = "[P2] The minor versions being used (%s) are not up-to-date (%s)."
const MSG_NOT_ALL_VERSIONS_SAME_CONCLUSION_1 string = "[P2] Not all nodes have the same Postgres version. Node %s uses Postgres %s."
const MSG_NOT_ALL_VERSIONS_SAME_CONCLUSION_N string = "[P2] Not all nodes have the same Postgres version. Nodes %s uses Postgres %s respectively."
const MSG_NOT_ALL_VERSIONS_SAME_RECOMMENDATION string = "[P2] Please upgrade Postgres so its versions on all nodes match."
const MSG_ALL_VERSIONS_SAME_CONCLUSION string = "All nodes have the same Postgres version (%s)."

const MSG_NOT_LAST_MINOR_VERSION_RECOMMENDATION string = "[P2] Please upgrade Postgres to the most recent minor version: %s."
const MSG_NO_RECOMMENDATION string = "No recommendations."
const MSG_GENERAL_RECOMMENDATION string = "  \n" +
	"For more information about minor and major upgrades see:  \n" +
	" - Official documentation: https://www.postgresql.org/docs  \n" + ///XX.YY/upgrading.html
	" - [Major-version upgrading with minimal downtime](https://www.depesz.com/2016/11/08/major-version-upgrading-with-minimal-downtime/) (depesz.com)  \n" +
	" - [Upgrading PostgreSQL on AWS RDS with minimum or zero downtime](https://medium.com/preply-engineering/postgres-multimaster-34f2446d5e14)  \n" +
	" - [Near-Zero Downtime Automated Upgrades of PostgreSQL Clusters in Cloud](https://www.2ndquadrant.com/en/blog/near-zero-downtime-automated-upgrades-postgresql-clusters-cloud/) (2ndQuadrant.com)  \n" +
	" - [Updating a 50 terabyte PostgreSQL database](https://medium.com/adyen/updating-a-50-terabyte-postgresql-database-f64384b799e7)  \n"
