On first terminal:
```
sudo bash -x setup-net-http.sh setup
sudo ip netns exec envoy_ns bash
```

In the NS: `bash setup-envoy.sh`

Open second terminal: There is different header added according to UE and target HTTP server

UE1:
```
sudo ip netns exec ue1_ns curl   3.3.3.3:80/index
sudo ip netns exec ue1_ns curl   4.4.4.4:80/index
```
UE2:
```
sudo ip netns exec ue2_ns curl   3.3.3.3:80/index
sudo ip netns exec ue2_ns curl   4.4.4.4:80/index
```
Destroy setup:
```
sudo bash -x setup-net-http.sh destroy
```
