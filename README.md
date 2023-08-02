# CLI for HashiCorp Vault - Secrets Management

This project is a command line interface (CLI) tool developed for interacting with secrets stored in HashiCorp Vault. It simplifies and streamlines the process of managing secrets by providing a user-friendly interface to perform secret management operations.

## Installation

Before you get started, ensure that you have HashiCorp Vault installed. 

Then, follow these steps to install the CLI and configure the vault and keycloak servers:

1. Using Docker: 

```bash
./run.sh
```

2. Using a Kubernetes cluster (minikube needs to be installed):

```bash
./runKube.sh
```

## Usage

To use this CLI, you will need to have access to a running instance of HashiCorp Vault and a keycloak sever, started by the run commands provided above. This will configure these servers, so you can use them seemlessy in the CLI. Once that's ready, you can start using the CLI by executing the desired command:

- With Userpass Authentication method:

```bash
./cliapp write --mount=kv --path=secret --key=name --value=Company.Inc --user=user --pass=pass
// expected output 
Secret Written Successfully.
```

- With Keycloak Authentication method:

```bash
./cliapp write --mount=kv --path=secret --key=name --value=Company.Inc
// expected output 
Open the following URL in your browser to authenticate with Keycloak:

http://127.0.0.1:8080/auth/realms/my_realm/protocol/openid-connect/auth?access_type=offline&client_id=vault-client&redirect
uri=http%3A%2F%2F127.0.0.1%3A3000%2Fcallback&response_type=code&scope=openid+profiletemail&state=state

Secret Written Successfully.
```

## Contributing

If you'd like to contribute, please fork the repository and use a feature branch. Pull requests are warmly welcome.

## Contact

If you have any questions or issues, please feel free to open an issue on this project:

dawid.skraba@gmail.com

