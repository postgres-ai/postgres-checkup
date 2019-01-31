import traceback
import sys
import os


def exception_helper():
    exc_type, exc_value, exc_traceback = sys.exc_info()
    return "\n".join(
        [
            v for v in traceback.format_exception(exc_type, exc_value, exc_traceback)
        ]
    )


def get_scalar(conn, query):
    p_query = conn.prepare(query)
    res = p_query()
    return None if len(res) == 0 else next(row[0] for row in res)


def get_resultset(conn, query):
    p_query = conn.prepare(query)
    return p_query()


def get_file_content(file_name):
    current_file = open(file_name, 'r')
    file_content = current_file.read()
    current_file.close()
    return file_content