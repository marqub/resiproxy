# ResiProxy
### ResiProxy is a [ToxiProxy](https://github.com/Shopify/toxiproxy) k8s companion:
 - intercepts REST proxy creation requests and open all the necessary ports
 - forwards all the requests to toxiproxy (mainly inspired by https://hackernoon.com/writing-a-reverse-proxy-in-just-one-line-with-go-c1edfa78c84b)
 - does not proxy non REST "admin" requests
 - compatible with ToxiProxy 2.1.3

### [ToxiProxy](https://github.com/Shopify/toxiproxy) Supported endpoints
All endpoints are JSON.

- GET /proxies - List existing proxies and their toxics
- POST /proxies - Create a new proxy
- ~~POST /populate - Create or replace a list of proxies~~
- GET /proxies/{proxy} - Show the proxy with all its active toxics
- POST /proxies/{proxy} - Update a proxy's fields
- DELETE /proxies/{proxy} - Delete an existing proxy
- GET /proxies/{proxy}/toxics - List active toxics
- POST /proxies/{proxy}/toxics - Create a new toxic
- GET /proxies/{proxy}/toxics/{toxic} - Get an active toxic's fields
- POST /proxies/{proxy}/toxics/{toxic} - Update an active toxic
- DELETE /proxies/{proxy}/toxics/{toxic} - Remove an active toxic
- POST /reset - Enable all proxies and remove all active toxics
- GET /version - Returns the server version number

Use the ingress to access ResiProxy from outside the cluster