Kubernetes:

why: enabling lifestyle,
all code and dependencies in an isolated environment - container -> easy to create multiple containers, easy to tear down/replace, works the same way no matter where you create it
not meant to hang around for a long time
has to run on a virtual machines existing, as long as resource permit

kubernetes : containers places inside pods, kub8s manages pods not containers, by putting your application in a container, k8 puts a wrapper around and calls it a pod into a one unit.  pod lives in a node

pods -> nodes
combining all the nodes you get a cluster
worker nodes -> vast majority of pods in these worker nodes, different sizes,
controller nodes -> hub of all the tools required to keep your cluster running. manager the cluster

kubectl: CLI tool in your workstation, access the cluster by connecting to the controller node. all the communication from kubectl is done via APIs. kubeconfig -> allow you to authenticate and talk to your cluster
API server: the central management entity, only component that directly connects to etcd -> a key-value database commonly deployed with distributed systems, distributed reliable key-value store
scheduler: pods put on nodes, k8s decides, this makes sure that the pods and nodes are following the rules
controller manager: manages inside the master node as replication controller, endpoints controller, namespace controller, service accounts controller
kubelet: focused on running containers, eyes and ears in the cluster, runs all nodes, has an internal http server has endpoints /health, /pods, /spec, if it encounters problem sole connection to the node, means the node has a problem, responsible for managing the container

containers are build by docker, inside every pod you have containers running meaning inside each container there is some sort of container build engine

everything in the cluster are communicated to talk with APIs

image - your application images, created usually through docker, podman or other software

namespace - access control

resourceQuota : limit for your namespace

resource management : cluster an aggregate of all the nodes that are networked together, one whole, checking how many resources that are getting consumed and add guardrails

you need to limit the resource consumptions of the containers inside the pods so they don't consumer more or unnecessary resources
requests(min), limits (hard cap),

k8s:maintaining the container

probes : watch dogs, put on top pods and container to enforce certain behaviors,
liveness probe: fail = dead
    initialDelaySeconds: How soon after creation are we probing
    periodSeconds: How often thereafter are we probing
    timeoutSeconds: How long we give the container to respond
    failureThreshold: How many contiguous failures before killing the container
readiness probe: fail = not ready yet
    initialDelaySeconds: How soon after creation are we probing
    periodSeconds: How often thereafter are we probing
    timeoutSeconds: How long we give the container to respond
    failureThreshold: How many contiguous failures before killing the container, turn off traffic to the pod

configmap : pieces of data, env variables, that need to be changed or required for the future pods to be made, disengaged from the container
attaching using volumes
configmap is mounted as a pod, all the configuration maps can be reused by all the other pods, all the changes will be done at a single place
volumes: the configmap pod is mounted as a volume for the future pods to refer
subPath: instead of connecting the WHOLE configmap, let's write a subpath to just a file in the configmap, use it when you want to write to an existing folder and other files are present and you don't to delete them

secrets:is an object that contains a small amount of sensitive data such as a password, token, a key. meaning you don't need to include confidential data in your application code
on their own not encrypted, use your own encryption solution
env:
    - name: env variable for the secret to substitute
    valueFrom:
        secretKeyRef:
            name: <name of the secret created>
            key: <the key you put in the secret, not the value, the key>

logs:all the containers output their logs individually, these are the standard prints, err statement

labels : tags for objects, used from grouping, viewing and operating many objects at the same time, key-value pairs

all pods die, they are not meant to last,

pods are made as part of the deployment
deployments are able to create multiple pods of the same configuration, deployment make pods -> part of the configuration the manifest will be configuration for the pods, can also version control rollout and rollback with zero downtime, they expect change, they are design to handle change use replicaset to rollback and forth

replicaset - the one that creates the pods, deployment creates replicaset which inturn creates pods
deciding which pods goes to which node is decided by the scheduler

storage - chart a path between the container pod and the storage. K8s is agnostic, you need to create a storage class
persistent volumes is used to tell cluster how to handle different storage types, they represent the storage itself, access persistent volume key
1:1 relationship,
manual storage class - pick a node and use some chunk and use that chunk as storage space
storageClassName - create in advance, readWriteOnce: one node can write at once

kubectl describe svc -> IP storage, the pods in this deployment will have IPs in this range, will track the ip address of the pods that are linked to the deployment. The IP address of the service will never change. The service will be what points to the correct direction
deploy will service which will allow pods created in this deployment talk with pods outside

network policies- namespace you wish to control; ingress traffic coming in egress traffic coming out; podselector which pods are we controlling;
ingress traffic should come from either of these ipBlock, namespaceSelector, podSelector
egress traffic should go to similar constraints and use a particular protocol

targetPort: refers to the port on the target,
port: refer to the port on the service

pod can have multiple containers in them, to be able to access different containers we match the ports and services
service ports will be used

kube-proxy - round robin the traffic to different pods.
IP address of the nodes is different
node port services external access without using the kubectl
node port service mesh kind



cluster -> node -> pod -> container

### Commands
kubectl get pod - gets pods

- kubectl create
- kubectl apply : create the pod and run
- kubectl get pods
- kubectl describe RESOURCE RESOURCE_NAME
- kubectl delete RESOURCE RESOURCE_NAME
- kubectl get pods -n <namespace name>
- kubectl create ns <namespace name> // create a new namespace
- kubectl run // create the pod without using the manifest, the image will be in the registry
- kubectl port-forward <podname> LOCALPORT:CONTAINERPORT
- kubectl exec -it <podname> --sh
- kubectl logs <pod name> -c <container name> // --all-container all the logs at once, -f follow the logs updated in real time, --since <time> last time passed
- kubectl get pods --show-labels
- kubectl label <pod-name> key=value // --overwrite updates the existing values and keys, key- removes the label with this key, -L <key> gets all the pods with this label key, --selector=<key> get the pods with the key
- kubectl get deployment
- kubectl delete deploy --all
- kubectl rollout history deploy <deployment name>
- kubectl rollout undo deploy <deployment name>
- kubectl expose deploy <deployment name>


resource: https://www.youtube.com/watch?v=MTHGoGUFpvE
