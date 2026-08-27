# manifest-separator

Separate YAML documents into individual files for easier comparison.

## Why

This project was born out of the necessity to split up an Istio manifest file into readable chunks. An IstioOperator file generated a manifest whose line count exceeded the 16 bit integer limit. It broke my editor's syntax highlighting.

## Building from source

Run `make build` to build from source.

The resulting binary will be in the `bin` directory; add it to your PATH or move it to a directory within your PATH.

## Modes

Manifest separator has three modes to handle parsing and separating manifests.

`-mode=dash` will separate YAML on triple dashes (`---`). This is the default behavior if `-mode` flag is omitted.

`-mode=list` will separate YAML that is in a `Kind: List`.

`-mode=appset` will separate generated ArgoCD ApplicationSet.

## Usage

The default mode is reading from Stdin.

Reading from a file or directory can be specified with `-f` flag
If reading from a directory, all files must be the same type (triple dashes, list, or appset)

## Example 1

`manifest-separator -f manifest.yaml`

This will attempt to separate `manifest.yaml` with `-mode=dash`

## Example 2

It is possible to combine manifest-separator with `kubectl get` commands like so:

`kubectl get namespace -A -oyaml | manifest-separator -mode=list`

## Example 3

Like the above example, it's possible to combine with ArgoCD CLI commands like so:

`argocd appset generate manifest.yaml | manifest-separator -mode=appset`
