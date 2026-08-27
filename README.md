# manifest-separator

Separate YAML documents into individual files for easier comparison.

## Why

This project was born out of the necessity to split up an Istio manifest file into readable chunks. An IstioOperator file generated a manifest whose line count exceed the 16 bit integer limit. It broke my editor's syntax highlighting.

## Building from source

Run `make build` to build from source.

The resulting binary will be in the `bin` directory; add it to your PATH or move it to a directory within your PATH.

## Usage

The default mode is reading from Stdin.

Reading from a file or directory can be specified with `-f` flag
If reading from a directory, all files must be the same type (triple dashes, list, or appset)

## Triple Dash vs List

Default behavior is to separate YAML on the triple dashes (`---`)

To separate YAML that is in a `Kind: List`, use `-list` or `-l` flag

## ArgoCD AppSet

Manifest Separator can separated generated ArgoCD AppSet with the `-appset` flag.
