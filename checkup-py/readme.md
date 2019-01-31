# Prepare

On windows

```bash
# install python
https://www.python.org/ftp/python/3.6.8/python-3.6.8-amd64.exe
python -m pip install --upgrade pip

cd C:\Program Files\Python36\Scripts
# or
cd C:\Users\USER\AppData\Local\Programs\Python\Python36\Scripts\

pip3.6.exe install py-postgresql
pip3.6.exe install paramiko
pip3.6.exe install sshtunnel
```

On RHEL

```bash
# install python
yum install -y python36
# if pip not installed
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
python3.6 get-pip.py
# else update pip
python3.6 -m pip install --upgrade pip

pip3.6 install py-postgresql
pip3.6 install paramiko
pip3.6 install sshtunnel
```

# Usage

```bash
python3.6 postgres-checkup.py \
	--ssh-host=194.67.206.177 \
	--ssh-port=22 \
	--ssh-user=root \
	--ssh-password=*** \
	--db-name=test \
	--check=indexes_invalid
```
