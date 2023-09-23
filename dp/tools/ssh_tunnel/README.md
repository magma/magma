### To run this utility

1. Go to magma lab and login
2. Choose Actions/command dialog
3. Fill generic command (connect)
4. Pass following parameters

```
 {
  "command": "connect",
  "params": {
    "shell_params": [
      "21000"
    ]
  }
}
```

and note shell_params number (here 21000) which is ssh port where connecion to agw will be opened.

5. Run following command. Note that variable `REMOTE_PORT` value should match number set in previously.

```
make connect JUMPHOST="YOUR_JUMPHOST" REMOTE_PORT="YOUR_REMOTE_PORT" KEY="ABS_PATH_TO_PRIVATE_SSH_KEY"
```

6. Now your able to reach enodebd port 60055 on your local machine port 60055, as well as enodebd can reach your process on port //TODO

7. To stop running container run command:

```
docker stop ssh
```
