# Multipass Setup

Start a new multipass instance with the following command:
```bash
multipass launch focal \
  --name orc8r \
  --disk 100G \
  --mem 8G \
  --cpus 4
```

Check if instance has started:
```bash
multipass ls
```

Get access to the shell:
```bash
multipass shell orc8r
```

Add your public SSH key to the instance:
```bash
vim .ssh/authorized_keys
```

### Uninstall

Delete Instance
```bash
multipass delete orc8r
multipass purge
```
