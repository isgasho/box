# box

A command-line utility (_non-daemon_) for creating Linux containers sandboxing
a process written in Go.

It basically is a tiny version of docker, it uses neither
[containerd](https://containerd.io/) nor
[runc](https://github.com/opencontainers/runc).
Only a set of the Linux features.

> __NOTE__: This is a heavily modified fork (_of was_) of [vessel](//github.com/0xc0d/vessel.git)
            and a reimplementation of the `box` utility from [ulinux](https://github.com/prologuc/ulinux).

## Features

`box` supports:

* __Control Groups__ for resource restriction (CPU, Memory, Swap, PIDs)
* __Namespace__ for global system resources isolation (Mount, UTS, Network, IPS, PID)
* __Union File System__ for branches to be overlaid in a single coherent file system. (OverlayFS)

## Install

```#!console
go get -u github.com/prologic/box
```

## Usage

...

## Examples

Run `/bin/sh` in `alpine:latest`

```#!console
box run alpine /bin/sh
box run alpine # same as above due to alpine default command
```

## Notes

`box` is/does __NOT__:

- Designed to be used in critical production workloads.
- Known to have any orchestrator(s) for managing services.
- Useful for multi-host networking and has no support for it.
- Have any support for volumes besides bind-mount(s) from the host
- Have any otehr features you'd expect from Docker, Docker Swarm or Kubernetes.

## License

`box` is licensed under the MIT License.
