# Adaptive Hive

A dynamic adaptive distributed computing system where applications are orchestrated 
across multiple devices based on available resources, application priorities, network 
topology, and declarative state. This is a learning project exploring distributed systems 
concepts through hands-on implementation. This is just an idea that I wanted to feel 
out and see where it takes me.

## Overview

Imagine teams operating in disconnected environments - each team has several devices 
that form an ad-hoc cluster. When teams are isolated, only critical applications run 
due to limited resources. When teams meet, their clusters merge, more resources become 
available, and additional applications can spin up. When teams return to headquarters, 
full connectivity enables all applications plus HQ-only services.

## Status

**Work in Progress** - Currently building the gossip layer (SWIM protocol implementation).

### Completed
- [ ] Transport abstraction with test harness
- [ ] SWIM protocol message types
- [ ] Membership list with state merge logic
- [ ] Message serialization (envelope-based gob)
- [ ] Gossip dissemination queue
- [ ] Failure detector (direct ping, indirect ping, suspicion)
- [ ] Gossiper coordinator (in progress)
- [ ] UDP transport
- [ ] CRDT state store
- [ ] Scheduler
- [ ] Agent

## Building
```bash
make build
```

## Testing
```bash
make test
```

With race detection:
```bash
go test -race ./...
```

## Project Structure
```
├── cmd/                    # Entry points
│   ├── agent/              # Device agent
│   └── hivectl/            # CLI tool
├── internal/
│   ├── config/             # Configuration
│   └── gossip/             # SWIM protocol implementation
├── pkg/
│   └── types/              # Shared domain types
└── test/
    └── simulation/         # Integration tests
```

## Concept
The system enables:

- **Hybrid architecture**: Centralized when connected to HQ, peer-to-peer when isolated

- **Resource-aware scheduling**: Apps run where resources are available

- **Graceful degradation**: Fewer devices = only critical apps run

- **State synchronization**: Devices sync when they encounter each other

## Learning Goals
- Distributed consensus and coordination

- Gossip protocols and failure detection

- Conflict-free replicated data types (CRDTs)

- Container orchestration patterns

- Network programming and fault tolerance

- Testing distributed systems

## References
- [SWIM Protocol Paper](https://www.cs.cornell.edu/projects/Quicksilver/public_pdfs/SWIM.pdf)
- [HashiCorp Memberlist](https://github.com/hashicorp/memberlist)
- [Serf](https://www.serf.io/)

## License

MIT
