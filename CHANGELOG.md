## v0.0.1  (2023-02-25)

### Features

* update minimal Go version to 1.18 ([a1b93d7](https://github.com/rangzen/otel-status/commit/a1b93d727508d033c016d933d05db5c4a3f2ca53))
* from log to slog ([7c3c330](https://github.com/rangzen/otel-status/commit/7c3c3306293d95306f7369ced95e930dad977844))
* add metrics ([c6b6a19](https://github.com/rangzen/otel-status/commit/c6b6a1956c85ac4d282c7a918c2c272cc2d7e740))
* add plugin name in span attributes ([84f0876](https://github.com/rangzen/otel-status/commit/84f08767d7a6bfa2d1bce049ab5503eb9a9adda3))
* **cli:** add config file reading ([6e86208](https://github.com/rangzen/otel-status/commit/6e86208daf7144873e81c4ff3de9a8b27b8fa15e))
* **cli:** better message if configuration file is missing ([d7a3dd7](https://github.com/rangzen/otel-status/commit/d7a3dd71ac4d7d7d3f327c0d0dc2d15a1b8168b9))
* **http:** add tests ([ac70e4b](https://github.com/rangzen/otel-status/commit/ac70e4b4611dbfa9936d342295b7772fa8c0fc87))
* **http:** add Values for additional tags ([5b0d807](https://github.com/rangzen/otel-status/commit/5b0d807c3558e40a00a82f2986e350f52d5f0097))
* **http:** add error metric ([a7c0483](https://github.com/rangzen/otel-status/commit/a7c048365a88c1e4d8d3ea8014533eca8f43ae3e))
* **http:** add cyclic status service ([c9cdb89](https://github.com/rangzen/otel-status/commit/c9cdb89f512fb366610730ac75f516f046b0e243))
* **http:** add a nonexisting domain in config ([b1dbcc7](https://github.com/rangzen/otel-status/commit/b1dbcc7e46c7dfee12b20f705dbdfb574ff64c54))
* **http:** add http.url key to duration metric ([70655fa](https://github.com/rangzen/otel-status/commit/70655fad43dcbe1081f75abc44a8770c4ab707ee))
* **otel:** activate httpcheck ([3acca34](https://github.com/rangzen/otel-status/commit/3acca342fb565a7731ab23b298abe9c49daf1c9a))
* **jetbrains:** add run and watcher confs ([9e74b06](https://github.com/rangzen/otel-status/commit/9e74b06dc41d42d0d221db5afd2b0b687520a775))

### Bug Fixes

* **http:** fix creation of error metric ([4c0d4df](https://github.com/rangzen/otel-status/commit/4c0d4df43b4ae4ee022585ed3ec782fa301cb542))
* **http:** fix url naming collision ([b69eab3](https://github.com/rangzen/otel-status/commit/b69eab3fe73cee983a11367a2115cd666278294a))
