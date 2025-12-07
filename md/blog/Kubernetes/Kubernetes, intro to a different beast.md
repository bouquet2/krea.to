In this blog post, I will talk about my first hands-on experience with Kubernetes. This may or may not be factual information, and is purely based on my experience.
## What is Kubernetes?
Kubernetes (sometimes referred to as k8s) is a open source system that allows you to automate and scale software developments. I'd like to think of it as Docker in steroids.
## Why Kubernetes?
It is much more scalable than Docker, and much more managable (once you know how to use it). It also is more feature-rich than Docker, and things like GKE make it very accessible.

## Okay, but how do I set it up?
This is where the fun begins.

Kubernetes' core is an "engine", a daemon responsible for managing Kubernetes resources.

Unlike Docker, there isnt one standard Kubernetes engine called "Kubernetes"
 that you can just download.

There are multiple choices, each having cons and pros. Popular ones include [minikube](https://minikube.sigs.k8s.io/), [k3s](https://k3s.io), [OKD](https://www.okd.io) managed engines such as [GKE](https://cloud.google.com/kubernetes-engine), [EKS](https://aws.amazon.com/eks/), and much more.

Minikube is generally meant for testing, while others have varying usecases. K3s works well for most usecases, so it is what I went with.

Setting up k3s is very easy, and just takes a command;

```
curl -sfL https://get.k3s.io | sh -
```

After setting it up, we can continue with learning the jargon.

## Kubernetes jargon
Keep in mind that this isn't a complete list, but it will be enough to get you started.

A [namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) is a way to isolate certain resources from the rest. Think of it as a seperate sandbox. They are mainly intended for use in environments with many users working in multiple teams, or projects. You can still use them to seperate your instances though. 

A [pod](https://kubernetes.io/docs/concepts/workloads/pods/) is a container instance.

A [secret](https://kubernetes.io/docs/concepts/configuration/secret/) is a value store that can contain things like passwords, and sensitive information. Think of it like a `.env` file for your pods.

An [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) is a object that manages external access for a service.

For more information, please check out the Kubernetes documentation.

These are a couple resource types you will most encounter while working with Kubernetes. There are a bunch more, but they are mostly self-explanatory and well-documented.

## Kubernetes CLI - kubectl
`kubectl` is the CLI utility to manage Kubernetes resources. It allows you to see, modify, delete, move, create and do bunch more with Kubernetes resouruces. You will most definitely need to get used to it to get better at Kubernetes.

`kubectl` has a lot of commands, but major ones that you can focus on are `create`, `logs`, `port-forward`, `exec`, `apply` and `delete`.

## Kubernetes package manager - Helm
Helm is the Kubernetes package manager. It is the most common way to install and manage new software. Using Helm Charts, you can create pretty much whatever you want in Kubernetes with ease.

I won't talk about creating Helm Charts in this blog post, as I haven't looked too deep on it yet.

You can install Helm with a single command;
```
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

Helm is also available [on most popular package managers](https://helm.sh/docs/intro/install/).

Helm also has lots of commands, but main ones you should look at are `list`, `uninstall`, `upgrade`, `install`, `repo`.

You should also check out [ArtifactHUB](https://artifacthub.io) as it has most popular Helm charts.

## Things I recommend installing
First, you should set Traefik with cert-manager (Tutorial for that coming soon!).

It is the most convenient way to put services to the outside world.

If you want to monitor your installation, [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack) with [grafana-dashboards-kubernetes](https://github.com/dotdc/grafana-dashboards-kubernetes) is the easiest and one of the best ways to monitor. It comes with bunch of well-made Grafana dashboards, ready-to-use Prometheus instance and more.

I also recommend checking out [awesome-k8s-resources](https://github.com/tomhuang12/awesome-k8s-resources) for more cool Kubernetes stuff.

## Conclusion
Kubernetes is an awesome tool that might look scary at first, but with enough usage is not that hard. I recommend anyone to check it out.