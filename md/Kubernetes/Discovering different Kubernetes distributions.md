In this blog post, I am going to look at the different Kubernetes distributions available today, and how they differ from each other. 

## What is a Kubernetes distribution?
Kubernetes has multiple parts that make it work. A Kubernetes distribution takes Kubernetes, and makes it easier to use while adding their own spin to it. The differences include (but not limited to);

* The CNI used for networking
* The default Ingress and LoadBalancer
* The interface that you can control resources in
* Support
* Purpose (i.e Bare-metal, IOT etc.)

## Kubernetes distributions that will be covered in this post
I've tried;

* Red Hat OpenShift
* K3s
* K0s
* Vanilla Kubernetes (through kubeadm)

Please keep in mind that I mostly looked at these distributions on a surface-level, and haven't had enough time to go through each and every feature.

We will start with K3s which is what I am most familiar with.

## K3s
K3s is a lightweight Kubernetes distribution that is originally developed by SUSE for IoT and Edge Computing scenarios. It is a single-binary Kubernetes distribution that is very easy to install. It is usually recommended for new users of Kubernetes, as it is very light and stable. It is also the first Kubernetes distribution I have used.

It comes with;

* Traefik as IngressController (Optional)
* Containerd/cri-dockerd as container runtime
* Flannel as CNI
* Local-path-provisioner for PVs
* No etcd

While it being really light is nice, it is usually not enough for multi-node clusters which is why I moved away from it. I still recommend it for learning and single-node clusters.

## K0s
K0s is an another single-binary Kubernetes distribution. While it is similar to K3s, there are a few differences;

* Works decent for Multi-node setups
* Kube-Router/Calico as CNI
* Etcd for multi-node clusters
* Management using `k0sctl` including backing up, setting up and upgrading clusters

k0sctl is really nice to use, you can bootstrap a cluster by just having a config and a SSH connection. It also allows you to completely back up the cluster, or easily upgrade it. While it is not as advanced as stuff such as openshift-client, it is still a nice to have. 

I didn't have a good experience with Kube-Router so I chose Calico when I set up the cluster. It managed to work well through connecting my two nodes through Tailscale. Highly recommend if you are starting out, wish I knew about it when I was setting up my first cluster.

## Vanilla Kubernetes
Vanilla Kubernetes (sometimes referred as Plain Kubernetes) is a Kubernetes setup that gets as basic as it gets, with no CNI or StorageClass by default set up. Vanilla Kubernetes nowadays is set up by using `kubeadm`, a very nice tool that allows you to setup/join clusters. It is similar to `k0sctl`, but requires more tinkering as not many things are set up by default.

Setting it up was quite the hurdle, and required extensive Kubernetes knowledge but it is still an incredible learning experience, and after setting it up right it can be rock stable.

For learners or people that want to tinker, Vanilla Kubernetes will teach you about how a cluster works very well but if you want a cluster that "just works", I would stay away from this one.

(PS: There might be a blog coming soon about setting up Vanilla K8s on Tailscale, stay tuned!)

## OpenShift
This is the most different Kubernetes distribution that is in this post, as others mostly try to stick to vanilla. OpenShift is the Enterprise version of OKD, which is "a very opinionated deployment of Kubernetes". Unlike others mentioned in this post, OpenShift works in its own OS called RHCOS (Red Hat Enterprise Linux CoreOS), while OKD runs on SCOS (CentOS Streams CoreOS). OpenShift has its own CNI, IngressController and CLI to manage clusters.

It is mainly for Multi-node setups, although single-node is also available. It has really high requirements compared to other distributions in this list. You can either use `openshift-install` or the Red Hat Hybrid Cloud Console to install it. I chose the Hybrid Console since I installed it for the first time.

It has a pretty seamless installation, I was in the UI in 15 minutes. Its UI is the most comprehensive Kubernetes UI I've ever seen. It includes pretty much everything you would ask for about Kubernetes.

It has also a CLI (openshift-client or oc) that is very nice to use. While in general my opinions are positive, it just requires too much resources to make sense for most users.

If you want to learn Red Hat stuff like I do, or if you have a massive homelab setup, I would recommend checking it out. Otherwise, the other distributions are enough to learn and be productive in Kubernetes.


## Conclusion (and what i picked)

All of these distributions have their own pros and cons. I recommend trying all of them before sticking with one (if you have the time, of course). For me, I picked OpenShift for now, as I want to learn OpenShift specific technology since I am focused at Red Hat stuff at the moment.

If you don't have time and just want something that works, I was very impressed with K0s. It is light and good enough for most usecases. K3s also works if you need something that is a bit more lighter.
