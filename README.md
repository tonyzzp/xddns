
### 帮助

```bash
xddns --help
xddns --help set
```

### 配置文件

配置文件`ddns-config.yaml` (**搜索路径顺序:workingdir,exedir**)，内容如下
```yaml
ali:
  region: cn-shenzhen
  key_id: xxxxx
  key_secret: xxxxx
  domains:
    - a.com
    - b.com

cloudflare:
  token: "xxxxx"
  domains:
    "a.me": "xxxxxxx zoneid"
    "b.app": "xxxx zoneid"
```

### 通过命令行参数指定配置文件

--config参数一定要加在子命令前面

```bash
xddns --config /my/path/config.yaml set ...
xddns --config /my/path/config.yaml update ...
```


### 直接设置dns

```bash
xddns set --type A --domain myname.com --value 192.168.1.101
xddns set --type CNAME --domain www.myname.com --value  myname.com
```


### 直接更新为本机ip
```bash
xddns update --type ipv4 --domain a.myname.com
xddns update --type ipv6 --domain a.myname.com
```