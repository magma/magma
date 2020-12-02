global:
  storageClass: nfs
rootUser:
  password: ${db_admin_password}
extraFlags: "--sql-mode=ANSI_QUOTES"
image:
  tag: 10.5.6-debian-10-r7
