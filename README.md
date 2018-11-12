ToxiProxy k8s companion:
 - intercept REST proxy creation requests and open all the necessary ports
 - forward all the requests to toxiproxy (mainly inspired from https://hackernoon.com/writing-a-reverse-proxy-in-just-one-line-with-go-c1edfa78c84b)
 - does not proxy non REST "admin" requests
 - compatible with ToxiProxy 2.1.3