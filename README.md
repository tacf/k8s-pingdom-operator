# K8s Pingdom Operator

K8s Pingdom Operator is a simple kubernates controller (based on the [operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/#operators-in-kubernetes)). It allows you to manage your _Pingdom_ service checks through the usage of [custom resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).

## Setup
Installing dependencies

```golang
go get -u
```

you may need to add `GO111MODULES=on`as a prefix for the command above.

## Building

In the root of the project just run

```golang
go build -o pingdom-controller .
```

## Running

Running it locally requires you to have you `kubeconfig` file properly setup so that the controller can easily communicate with the cluster api.

After setting the `kubeconfig` simply run

```shell
./pingdom-controller -kubeconfig=<path to kubeconfig>
```

## Limitations

As of now the definition and implementation of the controller only allows for simple _http_ checks to be created.