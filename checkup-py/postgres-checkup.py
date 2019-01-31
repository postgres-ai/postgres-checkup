from sshtunnel import SSHTunnelForwarder
from threading import Thread
import time
import argparse
import postgresql
from common.utils import *
from checks.indexes_invalid.impl import IndexesInvalid     # each check must be imported here


class SysConf:
    def __init__(self):
        self.current_dir = os.path.dirname(os.path.realpath(__file__))
        self.application_name = 'postgres-checkup'


class CheckupGlobal:
    sys_conf = None
    args = None
    common_report = [       # checks order defined
        IndexesInvalid()
    ]
    check_names = {v.__class__.check_name:v for v in common_report}

    def __init__(self):
        try:
            parser = argparse.ArgumentParser()
            parser.add_argument(
                "--db-name",
                type=str,
                default="test"
            )
            parser.add_argument(
                "--db-remote-port",
                type=int,
                default=5432
            )
            parser.add_argument(
                "--db-local-port",
                type=int,
                default=5400
            )
            parser.add_argument(
                "--db-user-name",
                type=str,
                default="postgres"
            )
            parser.add_argument(
                "--db-user-password",
                type=str,
                default="postgres"
            )

            parser.add_argument("--ssh-host", type=str)
            parser.add_argument("--ssh-port", type=str, default="22")
            parser.add_argument("--ssh-user", type=str)
            parser.add_argument("--ssh-password", type=str)

            parser.add_argument("--check", type=str)

            try:
                self.args = parser.parse_args()
            except:
                print(exception_helper())
                sys.exit(0)

            if not len(sys.argv) > 1:
                print("No arguments. Type -h for help.")
                sys.exit(0)

            self.sys_conf = SysConf()
        except SystemExit as e:
            print("Exiting...")
            sys.exit(0)
        except:
            print(exception_helper())
            print("Can't initialize application. Exiting...")
            sys.exit(0)


def run_checks():
    str_conn = 'pq://%s:%s@%s:%s/%s' % (
        Checkup.args.db_user_name,
        Checkup.args.db_user_password,
        '127.0.0.1',
        Checkup.args.db_local_port,
        Checkup.args.db_name
    )
    db_conn = postgresql.open(str_conn)

    get_resultset(db_conn, "select 1, 2, 3")
    if Checkup.args.check == "ALL":
        # TODO
        pass
    else:
        Checkup.check_names[Checkup.args.check].run(db_conn)

    db_conn.close()


if __name__ == "__main__":
    print('===========> checkup started')
    Checkup = CheckupGlobal()           # global object with configuration and logger

    try:
        server = SSHTunnelForwarder(
            Checkup.args.ssh_host,
            ssh_port=Checkup.args.ssh_port,
            ssh_username=Checkup.args.ssh_user,
            ssh_password=Checkup.args.ssh_password,
            remote_bind_address=('127.0.0.1', Checkup.args.db_remote_port),
            local_bind_address=('127.0.0.1', Checkup.args.db_local_port)
        )

        server.start()
    except:
        print(exception_helper())
        print("Can't start SSHTunnelForwarder. Exiting...")
        sys.exit(0)

    str_conn = 'pq://%s:%s@%s:%s/%s' % (
        Checkup.args.db_user_name,
        Checkup.args.db_user_password,
        '127.0.0.1',
        Checkup.args.db_local_port,
        Checkup.args.db_name
    )

    worker_threads = []
    print('Start threads initialization for host "%s"' % Checkup.args.ssh_host)
    worker_threads.append(Thread(target=run_checks, args=[]))

    for thread in worker_threads: thread.start()
    alive_count = 1
    print('Threads successfully initialized for host %s' % Checkup.args.ssh_host)
    while alive_count > 0:
        alive_count = len([thread for thread in worker_threads if thread.is_alive()])
        if alive_count == 0: break
        print('Live %s threads' % alive_count)
        time.sleep(1)

    server.stop()
    print('<========== checkup finished')