# gosong (Work In Progress)

A straightforward deployment just like deployer.org but made by me, LOL

## Features
- Multi host auto deployment (local and SSH)
- Task execution and management, you can add yours (currenly via yaml)
- Deploy versioning
- Rollback mechanism (not yet)
- Process management, like supervisor (not yet)
- Simple yet reverse proxy (not yet)
- Seemlesly integrating blue green deployment by switching port, that's why i need to write my own process management and reverse proxy (not yet)

## Usage
1. Configure your deployment, see `docs/example.yaml` for example.
2. Build the project:
   ```sh
   make install && make build
   ```
3. Run commands:
   ```sh
   ./build/current/gosong [command] [flags]
   ```

## Requirements
- Go 1.18 or newer
- Electricity

## License
See [LICENSE](LICENSE).
