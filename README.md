# Cablebox Control

A Go-based control system for managing [FieldStation42](https://github.com/shane-mason/FieldStation42), a broadcast TV simulator.

## Overview

This project adds a simple remote control web interface for FieldStation42, allowing you to control the TV simulator from your browser. The project uses Nix for dependency management and reproducible builds, although it is not required. You can also install it using standard Go tools by running `go install github.com/scottjab/cablebox-control@latest`.

## Requirements

- Go 1.21 or later
- Nix package manager

## Building

The project uses Nix for dependency management. To build the project:

```bash
nix build
```

## Development

To set up the development environment:

```bash
nix develop
```

## Project Structure

- `cmd/cablebox-control/` - Main application code
- `flake.nix` - Nix flake configuration
- `go.mod` - Go module definition

## Integration with FieldStation42

This control system integrates with FieldStation42 by:
- Managing channel changes through FieldStation42's `runtime/channel.socket`
- Monitoring playback status via `runtime/play_status.socket`
- Providing additional control capabilities beyond the basic FieldStation42 interface

