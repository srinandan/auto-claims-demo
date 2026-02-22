[🔙 Back to Main Project README](../README.md)

# Repair Shop Agent

ReAct agent with A2A protocol [experimental].
Agent generated with [`googleCloudPlatform/agent-starter-pack`](https://github.com/GoogleCloudPlatform/agent-starter-pack) version `0.36.0`.

## Service Description

The **Repair Shop Agent** is a specialized AI agent responsible for finding suitable repair shops and handling appointment bookings for the final steps of auto claim processing. It interacts with the main AI Service and receives input from other services to orchestrate repairs.

It implements the [A2A (Agent-to-Agent) Protocol](https://a2a-protocol.org/) for standardized communication.

## Quick Start

Install required packages and launch the local development environment:

```bash
make install && make playground
```

## Commands

| Command              | Description                                                                                 |
| -------------------- | ------------------------------------------------------------------------------------------- |
| `make install`       | Install dependencies using uv                                                               |
| `make playground`    | Launch local development environment                                                        |
| `make local-repair-shop-agent` | Launch local development server on port 8083                                   |

For full command options and usage, refer to the [Makefile](Makefile).
