# TopoRank

Lightweight Go library/demo scaffold for a Topological Ranking algorithm.

Project layout:

.
# TopoRank

TopoRank is a small Go library demonstrating a topological ranking
approach based on random-walk score propagation and distance-weighted
potentials. The codebase is intentionally compact and suitable for
experimentation, learning, and small-scale demos.

Repository layout (relevant files):

.
├─ examples/
│  └─ toporank-demo/
│     └─ main.go
├─ internal/
│  ├─ graph/                 # shortest-path helpers
│  ├─ potentials/            # topological potential computation
│  └─ walker/                # random-walk propagation + ranking helpers
├─ api/                      # public TopoRank entrypoint
├─ types/                    # small graph and config types
└─ go.mod

Getting started
---------------

Build and run the small example included in `examples/toporank-demo`:

```bash
go run ./examples/toporank-demo
```

Library overview
----------------

- `types` provides minimal `Graph` and `Node` types used by the example.
- `internal/potentials` computes a Gaussian-style distance-decayed
	potential for each node.
- `internal/walker` implements a simple random-walk propagation engine and
	a helper to rank nodes by score.
- `api.RunTopoRank` ties components together and exposes a single entry
	point for running the algorithm on a `types.Graph`.

Contributing
------------

This repository is a compact scaffold. Improvements could include:

- more efficient shortest-path methods (BFS is currently used with a
	high-cost sentinel when unreachable),
- unit tests and benchmarks, and
- additional APIs for weighted or undirected graphs.

License
-------

This example project is provided without an explicit license file. If you
intend to use it in production, please add an appropriate `LICENSE`.
