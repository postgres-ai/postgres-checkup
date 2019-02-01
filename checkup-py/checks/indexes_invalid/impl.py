import os
from common.utils import *

class IndexesInvalid:
    check_name = "indexes_invalid"

    def __init__(self, ctx):
        self.current_dir = os.path.dirname(os.path.realpath(__file__))
        self.ctx = ctx

    def run(self, db_conn):
        # print("Hello from IndexesInvalid %s" % self.current_dir)
        tbl_1 = get_resultset(db_conn, get_file_content(os.path.join(self.current_dir, "list_invalid_indexes.sql")))
        tbl_2 = get_resultset(db_conn, get_file_content(os.path.join(self.current_dir, "drop_create_ddls.sql")))
        self.to_md(tbl_1, tbl_2)

    def to_md(self, tbl_1, tbl_2):
        md_content = ""
        md_content += "# Indexes -> invalid #\n\n"
        md_content += "## Observations ##\n\n"
        md_content += "### Master %s ###\n" % self.ctx.args.db_name

        if len(tbl_1) > 0:
            md_content += "Schema name | Table name | Index name | Index size\n"
            md_content += "------------|------------|------------|------------\n"
            for row in tbl_1:
                md_content += "%s | %s | %s | %s\n" % (row[0], row[1], row[2], row[3])
        else:
            md_content += "Invalid indexes not found\n"

        md_content += "\n\n## Recommendations ##\n\n"

        if len(tbl_2) > 0:
            md_content += """\n\n#### "DO" database migration code ####\n\n"""
            md_content += "```\n"
            md_content += """-- Call each line separately. "CONCURRENTLY" queries cannot be\n"""
            md_content += """-- combined in multi-statement requests.\n"""
            for row in tbl_2:
                md_content += "%s\n" % (row[0])
            md_content += "```\n"
            md_content += """\n\n#### "UNDO" database migration code ####\n\n"""
            md_content += "```\n"
            for row in tbl_2:
                md_content += "%s\n" % (row[1])
            md_content += "```\n"

        save_to(os.path.join(self.ctx.sys_conf.current_dir, "output", self.check_name + ".md"), md_content)