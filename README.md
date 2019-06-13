# duck
Expose local services to the Internet

## Usage

On a machine with a public IP, run:

```
./duck -l -addr :9990 -p test_password
```

on your local machine, run:

```
./duck -addr YOUR_SERVER_IP:9990 -p test_password <local service port> <local service port> ...
```

Now you can connect to the duck server and the tcp data will be forwarded to your local server. 