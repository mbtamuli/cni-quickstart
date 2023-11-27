# CNI Quickstart

:warning: This is a extremely experimental :construction: WIP repository. Read any further at your own risk! :warning:

This CNI plugin does nothing but attempt to modify the pod's network related details! :laughing:

Actually, this plugin doesn't modify any of the pod networking at all, just the information shown in the pod description.

```
kubectl describe pods test
Name:             test
Namespace:        default
Priority:         0
Service Account:  default
Node:             cni-test-control-plane/172.18.0.2
Start Time:       Mon, 27 Nov 2023 20:53:06 +0530
Labels:           run=test
Annotations:      <none>
Status:           Running
IP:               1.2.3.0
IPs:
  IP:  1.2.3.0
```

But if you exec into the pod and try to reach the internet it still does! :astonished: _<insert-astonished-Pickachu-GIF>_

```
+ kubectl exec -it test -- sh
/ # ping 1.1.1.1
PING 1.1.1.1 (1.1.1.1): 56 data bytes
64 bytes from 1.1.1.1: seq=0 ttl=62 time=62.596 ms
64 bytes from 1.1.1.1: seq=1 ttl=62 time=55.099 ms
^C
--- 1.1.1.1 ping statistics ---
2 packets transmitted, 2 packets received, 0% packet loss
round-trip min/avg/max = 55.099/58.847/62.596 ms
/ # ping kubernetes.io
PING kubernetes.io (147.75.40.148): 56 data bytes
64 bytes from 147.75.40.148: seq=0 ttl=62 time=100.265 ms
64 bytes from 147.75.40.148: seq=1 ttl=62 time=105.223 ms
^C
--- kubernetes.io ping statistics ---
2 packets transmitted, 2 packets received, 0% packet loss
round-trip min/avg/max = 100.265/102.744/105.223 ms
/ # apk add curl
(1/7) Installing ca-certificates (20230506-r0)
(2/7) Installing brotli-libs (1.0.9-r14)
(3/7) Installing libunistring (1.1-r1)
(4/7) Installing libidn2 (2.3.4-r1)
(5/7) Installing nghttp2-libs (1.57.0-r0)
(6/7) Installing libcurl (8.4.0-r0)
(7/7) Installing curl (8.4.0-r0)
Executing busybox-1.36.1-r2.trigger
Executing ca-certificates-20230506-r0.trigger
OK: 16 MiB in 26 packages
/ # curl -kO https://raw.githubusercontent.com/containernetworking/cni/main/logo.png
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 11604  100 11604    0     0  24207      0 --:--:-- --:--:-- --:--:-- 24175
```

If you dig deeper, you find the reality. :smiling_imp:
```
/ # ip a s
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: tunl0@NONE: <NOARP> mtu 1480 qdisc noop state DOWN qlen 1000
    link/ipip 0.0.0.0 brd 0.0.0.0
3: ip6tnl0@NONE: <NOARP> mtu 1452 qdisc noop state DOWN qlen 1000
    link/tunnel6 00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00 brd 00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00
4: eth0@if151: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 65535 qdisc noqueue state UP
    link/ether ae:89:21:5f:e3:5f brd ff:ff:ff:ff:ff:ff
    inet 10.244.0.148/24 brd 10.244.0.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::ac89:21ff:fe5f:e35f/64 scope link
       valid_lft forever preferred_lft forever
/ # ip route
default via 10.244.0.1 dev eth0
10.244.0.0/24 via 10.244.0.1 dev eth0  src 10.244.0.148
10.244.0.1 dev eth0 scope link  src 10.244.0.148
```

## Steps

1. Get a running Kubernetes cluster.
  ```
  kind create cluster --name cni-test --kubeconfig ~/.kube/cni-test.yml
  ```
2. Build the binary
  ```
  # I have a M2 MacBook so the Kind Cluster node I have also has arm64 architecture
  env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o cni-quickstart main.go
  ```
3. Copy the binary and the network config to enable the plugin
  ```
  docker cp cni-quickstart cni-test-control-plane:/opt/cni/bin/
  docker cp 10-kindnet.conflist cni-test-control-plane:/etc/cni/net.d/
  ```
4. Run a pod and see the magic! :sunglasses:
  ```
  kubectl run test --image alpine -- sleep infinity
  ```

## Gotchas

Can't run go build if you try to import "github.com/containernetworking/plugins/pkg/ns" on macOS as the "github.com/containernetworking/plugins/pkg/ns" contains files named `ns_linux.go` and `ns_windows.go` triggering implicit build constraint! - https://pkg.go.dev/cmd/go#hdr-Build_constraints

```
go build -o cni-quickstart main.go
github.com/containernetworking/plugins/pkg/ns: build constraints exclude all Go files in /Users/mriyam.tamuli/workspace/go/pkg/mod/github.com/containernetworking/plugins@v1.3.0/pkg/ns
```

Can build using
```
env GOOS="linux" GOARCH="amd64" CGO_ENABLED="0" go build -o cni-quickstart main.go
```
