# Kubernetes capstan

capstan is a benchmark which contains series of workloads and testing tools for Kubernetes. You can obtain the performance data of each workload and each component in the specific configuration of Kubernetes cluster offered by different cloud offering.

## What is the scope of this project?

capstan aims to provide a series of workloads and testing tools for Kubernetes cluster, which contains two component:

- capstan: manage cluster lifecycle, start propeller service, analysis and display the testing results.

- propeller: run every workload's testing cases, collect the testing data and the performance data of Kubernetes component and Kubernetes cluster.

## What is not in scope for this project?

- Building a new cluster lifecycle management tool(e.g. [kubeadm](https://github.com/kubernetes/kubeadm),[kops](https://github.com/kubernetes/kops),[kubernetes-anywhere](https://github.com/kubernetes/kubernetes-anywhere)).

- Building a new data collection and analysis tool(e.g. [cadvisor](https://github.com/google/cadvisor),[heapster](https://github.com/kubernetes/heapster)).

## Roadmap

- Basic cluster lifecycle management including creating, destroying and other operations (Optional).
- Design the testing indicators（P0).
- Design and Implement the framework of capstan（P0).
- Implement multiple workloads and testing tools（P1).
- Implement the online ranking system（P2).
