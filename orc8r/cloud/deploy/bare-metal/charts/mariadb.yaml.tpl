global:
  storageClass: nfs
image:
  tag: 10.3.22-debian-10-r27
master:
  extraFlags: --sql-mode=ANSI_QUOTES
rootUser:
  password: ${db_root_password}
