# localnet
用于获取用户localdns的出口ip，subnet 信息，用于判断用户dns是否设置正确，可能导致CDN调度异常

## 使用

### 部署
```
[root@xxx dns]# bash ./bulid.sh [local|linux]
[root@xxx dns]# ./localnet
2023/02/02 14:57:15 start dns server listen :53(udp)
[root@xxx dns]# dig @127.0.0.1 xx.example.com

```

### 自有域名dns设置（假设你持有域名： example.org）
#### 添加两条解析
* check  ns记录  ns1.example.org
* ns1    a记录   1.1.1.1 (这个IP是你localnet的公网IP)


### 测试
支持 a记录和txt记录
```
dig 1.check.example.org
;; ANSWER SECTION:
1.check.example.org.	0	IN	A	127.0.0.1 // 你localdns的出口IP
1.check.example.org.	0	IN	A	1.1.1.1 // 你localdns附加的subnet网段(可能没有)

dig 1.check.example.org txt 
;; ANSWER SECTION:
1.check.example.org.	0	IN	TXT	"localdns=127.0.0.1" // 你localdns的出口IP
1.check.example.org.	0	IN	TXT	"subnet=1.1.1.1/32/0" // 你localdns附加的subnet网段(可能没有)
1.check.example.org.	0	IN	TXT	"request_id=e60b5db7-f358-4f63-9402-8c2e70443d8a" // 查询的request_id，和日志结合排查问题，一般不需要关注
```


