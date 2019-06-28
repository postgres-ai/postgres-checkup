package l003

const MSG_HIGH_RISKS_TABLE string = "    - `%s`: reached value %d, %.0f%% of `%s` capacity"
const MSG_HIGH_RISKS_CONCLUSION_1 string = "[P1] High risks of out-of-range errors for integer column. The column listed below, being part of " +
	"primary keys, have high risks to reach 100%% of the integer capacity (2^31-1 = 2147483647 for `int4` columns, and 2^15-1 = 32767 for `int2` ones; see " +
	"[the documentation](https://www.postgresql.org/docs/current/datatype-numeric.html). " +
	"Once it happens, INSERTs of new rows will not be possible (unless they use some non-incremental " +
	"values, such as negative values) and the fix will require long downtime. Here is that %d column:  \n%s."

const MSG_HIGH_RISKS_CONCLUSION_N string = "[P1] High risks of out-of-range errors for integer columns. The columns listed below, being part of " +
	"primary keys, have high risks to reach 100%% of the integer capacity (2^31-1 = 2147483647 for `int4` columns, and 2^15-1 = 32767 for `int2` ones; see " +
	"[the documentation](https://www.postgresql.org/docs/current/datatype-numeric.html). " +
	"Once it happens, INSERTs of new rows will not be possible (unless they use some non-incremental " +
	"values, such as negative values) and the fix will require long downtime. %d such column are found:  \n%s"

const MSG_HIGH_RISKS_RECOMMENDATION string = "[P1] High risks of out-of-range errors for integer columns. Consider using `int8` in all PK columns, " +
	"always. To convert existing columns to `int8`, consider the following approaches:  \n" +
	"    1. Blocking `UPDATE`: a straightforward solution requiring downtime (maintenance window).  \n" +
	"    1. \"New column\": create a new column, update it in batches (lasting up not more than a few seconds, " +
	"not to block other queries), and then switch to using it, redefining all the constraints. Notice, that " +
	"to re-define a primary key constraint, `ALTER TABLE .. ALTER COLUMN .. SET NOT NULL` will be needed. " +
	"It is a blocking operation, up to Postgres 12 (where it might be lightweight if a proper `CHECK` " +
	"constraint is defined). Since Postgres 11, it is possible to use a trick: when adding a column, use " +
	"`DEFAULT` with `NOT NULL`, it will be not a blocking operation. For Postgres versions, prior 11, a " +
	"specific downtime (maintenance window) will be needed anyway.  \n" +
	"    1. \"New table\": create a new table with the same schema as the existing one, capture all ongoing " +
	"changes to an additional \"log\" table, copy existing data from the old table to the new one, and switch. " +
	"This method, as the previous one, is non-trivial and requires careful development and testing under " +
	"load (consider using [Nancy](https://gitlab.com/postgres-ai/nancy) for database experiments developing " +
	"this solution)."
