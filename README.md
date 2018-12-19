[![CircleCI](https://circleci.com/gh/marqub/resiproxy/tree/master.svg?style=svg)](https://circleci.com/gh/marqub/broadcast/tree/master)
[![codecov](https://codecov.io/gh/marqub/resiproxy/branch/master/graph/badge.svg)](https://codecov.io/gh/solcates/gobwa) 
# ResiProxy
### For more info about how to use ResiProxy, have a look at [this story on Medium](http://bit.ly/2EAXRYk)

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

The helm chart creates several K8s resources
```
$ helm install chart/ --namespace toxy --name resiproxy
NAME:   resiproxy
LAST DEPLOYED: Tue Nov 13 17:24:50 2018
NAMESPACE: toxy
STATUS: DEPLOYED

RESOURCES:
==> v1/ServiceAccount
NAME       SECRETS  AGE
resiproxy  1        0s

==> v1beta1/ClusterRole
NAME            AGE
resiproxy-read  0s

==> v1beta1/ClusterRoleBinding
NAME            AGE
resiproxy-read  0s

==> v1/Service
NAME                 TYPE       CLUSTER-IP      EXTERNAL-IP  PORT(S)   AGE
resiproxy-resiproxy  ClusterIP  xxx.xxx.xxx.xxx  <none>       8080/TCP  0s
resiproxy-toxiproxy  ClusterIP  xxx.xxx.xxx.xxx   <none>      8474/TCP  0s

==> v1beta1/Deployment
NAME                 DESIRED  CURRENT  UP-TO-DATE  AVAILABLE  AGE
resiproxy-resiproxy  1        1        1           0          0s

==> v1beta1/Ingress
NAME                 HOSTS                 ADDRESS  PORTS  AGE
resiproxy-resiproxy  resiproxy.marqub.com  80, 443  0s

==> v1/Pod(related)
NAME                                 READY  STATUS       RESTARTS  AGE
resiproxy-resiproxy-cfc4ccbff-6v9w5  0/3    Init:0/1     0         0s

==> v1/ConfigMap
NAME                 DATA  AGE
configmap-resiproxy  4     0s
```
Use the ingress to access ResiProxy from outside the cluster `https://resiproxy.marqub.com`
A sample request to create a proxy would be
```
curl -X POST http://resiproxy.resiliency-testing.com/proxies \  
  -H 'Content-Type: application/json' \
  -d '{
    "name": "proxy_service2",
    "listen": "[::]:8081",
    "upstream": "service2-go-service.resiliency-testing:8080",
    "enabled": true
}'
```
This sample request opens the port 8081 at the `toxiproxy` service level to redirect the incoming traffic to the port `8080` of `service2` in the namespace `resiliency-testing`: `service2-go-service.resiliency-testing:8080`
```
$ kubectl describe service resiproxy-toxiproxy -n resiliency-testing
Name:                     resiproxy-toxiproxy
Namespace:                resiliency-testing
Labels:                   app=resiproxy
                          chart=resiproxy-0.0.1
                          heritage=Tiller
                          release=resiproxy
Annotations:              <none>
Selector:                 app=resiproxy,release=resiproxy
Type:                     NodePort
IP:                       xxx.xxx.xxx.xxx
Port:                     http-toxiproxy  8474/TCP
TargetPort:               8474/TCP
NodePort:                 http-toxiproxy  31974/TCP
Endpoints:                xxx.xxx.xxx.xxx:8474
Port:                     8081  8081/TCP
TargetPort:               8081/TCP
NodePort:                 8081  30795/TCP
Endpoints:                xxx.xxx.xxx.xxx:8081
Session Affinity:         None
External Traffic Policy:  Cluster
Events:                   <none>
```
Use the toxiproxy service to redirect your traffic to your dependencies from any services running inside your cluster; for example to send a request to service2 use `resiproxy-toxiproxy.toxy:8081`
