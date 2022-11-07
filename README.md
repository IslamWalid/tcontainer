# Tiny Container
Simple container implementation in go using linux namespaces and cgroups.

## Usage and Installation:
- **Install tcontainer:**
```sh
go install github.com/IslamWalid/tcontainer/cmd/tcontainer@latest
```
- **Run the container:**
```sh
sudo tcontainer run <cmd> <args>
```
**NOTE:** running the container for the first time may take a while for the initialization step.
