## Load testing

I used [Hey](https://github.com/rakyll/hey)

```bash
url_api="http://localhost:8080"
data='{"name":"Bob"}'

hey -n 10000 -c 1000 -m POST -T "Content-Type: application/json" -H "rwaapi_token:tada" -H "rwaapi_data:hello world" -d ${data} "${url_api}"

url_api="http://localhost:8080"
data='Bob'
hey -n 10000 -c 1000 -m POST -d ${data} "${url_api}"
```
