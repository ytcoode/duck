# duck
Expose local servers behind NAT or firewall to the Internet

## Usage

On a publicly accessible machine, run:

```
./duck -l -addr :9990 -p test_password
```

On the machine where your local servers are running, run:

```
./duck -addr DUCK_SERVER_IP:9990 -p test_password <local server port> <local server port> ...
```

Now you can connect to the public duck server and the tcp traffic will be forwarded to your specific local server according to port.