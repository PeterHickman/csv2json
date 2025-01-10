# csv2json

In a perfect world this would not exist as this functionality is available in the `column` command. On modern Unix clones that is but in the world of BSD varients nothing so useful exists. So here we are

```bash
$ csv2json --input /etc/passwd --delimit : -names username,password,uid,gid,gecos,home_dir,shell
```