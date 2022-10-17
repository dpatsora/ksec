# Ksec
Are you also tired of working with base64 encoded data in K8s secrets?
```bash
➜ ~ kubectl get secret db-user-pass  -oyaml
apiVersion: v1
data:
  password.txt: MWYyZDFlMmU2N2Rm
  testKey: dGVzS2V5MQ==
  username.txt: c3VwZXJhZG1pbg==
kind: Secret
metadata:
  creationTimestamp: "2022-10-15T20:13:58Z"
  managedFields:
  ...
➜ ~ echo -n "MWYyZDFlMmU2N2Rm" | base64 -d
1f2d1e2e67df%
```
Endless copy-pasting pisses you off?

### I feel you. Let me offer you a solution


Ksec is a CLI tool for k8s secrets manipulations.

It simplifies the operations, because you don't need to manually encode/decode base64 values.


## Installation

```bash
go install github.com/dpatsora/ksec@latest
```

## Features
### Retrieve secret data in human-readable format
```bash
➜  ksec read db-user-pass
password.txt: 1f2d1e2e67df
testKey: tesKey1
name.txt: default_name
```

### Write key/value pair in your secret
```bash
➜  ksec write db-user-pass username.txt admin
```

### Update key/value pair in your secret 

```bash
➜  ksec write db-user-pass username.txt superadmin
Current value: admin
New value: superadmin

Do you want to continue with this operation? [y|n]: y
```

### Get usage example for every command
```bash
➜  ksec write -h
Write key/value pair to secret data

To add "USER_PASSWORD: admin123" to "db-pass" secret data, located in "core" namespace, command will be:
ksec write db-pass USER_PASSWORD admin123 -n core

Usage:
  ksec write [flags]

Aliases:
  write, w

Flags:
  -h, --help   help for write

Global Flags:
      --kubeconfig string   path to k8s configuration file
  -n, --namespace string    resource k8s namespace (default "default")
```

## Configuration

The only required configuration is Kubeconfig.

By default, it would be taken from `KUBECONFIG` env var. But you can pass it with `--kubeconfig` flag (flag takes precedence over env var).
