package l003

const MSG_HIGH_RISKS_TABLE string = "    - `%s`: reached value %d, or %.0f%% of `%s` capacity\n"

const MSG_HIGH_RISKS_CONCLUSION_1 string = "[P1] High risks of out-of-range errors for an integer column. " +
    "The column listed below, being part of a primary key, has high risks to reach 100%% " +
    "of the integer capacity (`2^31-1`, or `2147483647` for `int4` columns, and `2^15-1`, or `32767` for `int2` columns; " +
	"see [the documentation](https://www.postgresql.org/docs/current/datatype-numeric.html). " +
	"If it happens, INSERTs of new rows are not be possible (unless they use some non-incremental " +
	"values, such as some negative values) and fixing it will require a long downtime. Here is that %d column:\n%s"

const MSG_HIGH_RISKS_CONCLUSION_N string = "[P1] High risks of out-of-range errors for an integer column. " +
    "The columns listed below, being part of a primary key, have high risks to reach 100%% " +
    "of the integer capacity (`2^31-1`, or `2147483647` for `int4` columns, and `2^15-1`, or `32767` for `int2` columns; " +
	"see [the documentation](https://www.postgresql.org/docs/current/datatype-numeric.html). " +
	"If it happens, INSERTs of new rows are not be possible (unless they use some non-incremental " +
	"values, such as some negative values) and fixing it will require a long downtime. %d such column are found:\n%s"

const MSG_HIGH_RISKS_RECOMMENDATION string = "[P1] High risks of out-of-range errors for an integer column. " +
    "Consider using `int8` in all PK columns,  always. To convert existing columns to `int8`, consider the " +
    " following approaches:\n" +
	"    1. Blocking `ALTER TABLE .. ALTER COLUMN`: a straightforward solution requiring significant downtime (a maintenance window).\n" +
	"    1. \"New column\": create a new column, update it in batches (runnong not longer than a few seconds, " +
	"to avoid blocking issues), and then switch to using it, redefining all the constraints. Notice, that " +
	"to redefine a primary key constraint, `ALTER TABLE .. ALTER COLUMN .. SET NOT NULL` will be needed. " +
	"It is a blocking operation in all Postgres versions up to 12 (where it might be lightweight if a proper `CHECK` " +
	"constraint is defined first; such constraint can be defined in a non-blocking way). " +
	"Since Postgres 11, it is possible to use a trick: when adding a column, use " +
	"`DEFAULT` with `NOT NULL`, it will be a non-blocking operation. For all Postgres versions prior to 11, a " +
	"specific downtime (maintenance window) will be needed anyway.\n" +
	"    1. \"New table\": create a new table with the same schema as the existing one, capture all ongoing " +
	"changes to an additional \"log\" table, copy existing data from the old table to the new one, and switch. " +
	"This method, as the previous one, is non-trivial and requires careful development and testing under " +
	"load (consider using [Nancy](https://gitlab.com/postgres-ai/nancy) for database experiments developing " +
	"this solution). This approach is non-blocking regardless of Postgres version, but it requires significantly " +
	"more efforts to implement."
