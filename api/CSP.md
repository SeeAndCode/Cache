# CSP

Cache Serialization Protocol，基于TCP实现的应用层协议，Cache-Server默认基于本协议向外提供服务。

## ABNF描述

```ABNF
request = data-command-req;

key = length SP content
value = length SP content
key-value = length SP length SP content content
length = 1 * DIGIT
content = *OCTET

data-command-req = set-req
set-req = "00001" SP key-value;

data-command-req =/ get-req
get-req = "00002" SP key

data-command-req =/ getrange-req
getrange-req = "00003" SP start SP end SP key
start = 1 * DIGIT
end = 1 * DIGIT

data-command-req =/ getset-req
getset = "00004" SP key-value

data-command-req =/ getbit-req
getbit-req = "00005" SP offset SP key
offset = 1 * DIGIT

data-command-req =/ mget-req
mget = "00006" SP num SP 1 * key

data-command-req =/ setbit-req
setbit = "00007" SP offset SP key-value



response = data-command-resp

code = 1 * DIGIT
message = bytes-array

data-command-resp = set-resp
set-resp = "10001" SP code SP message

data-command-resp =/ get-resp
get-resp = "10002" SP code SP message value



```