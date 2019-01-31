import os
from common.utils import *

class IndexesInvalid:
    check_name = "indexes_invalid"

    def __init__(self):
        self.current_dir = os.path.dirname(os.path.realpath(__file__))

    def run(self, db_conn):
        print("Hello from IndexesInvalid %s" % self.current_dir)
        res = get_resultset(
            db_conn,
            get_file_content(os.path.join(self.current_dir, "list_invalid_indexes.sql"))
        )
        print(str(res))

        res = get_resultset(
            db_conn,
            get_file_content(os.path.join(self.current_dir, "drop_create_ddls.sql"))
        )
        print(str(res))

    def to_md(self):
        pass
