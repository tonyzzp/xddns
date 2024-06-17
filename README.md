

### 直接设置dns

```bash
ali-ddns set --type A --domain myname.com --r @ --value 192.168.1.101
ali-ddns set --type CNAME --domain myname.com --r www --value  www.myname.com
```


### 直接更新为本机ip
```bash
ali-ddns update --type ipv4 --domains a.myname.com,b.myname.com
ali-ddns update --type ipv6 --domains a.myname.com
```