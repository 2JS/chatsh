# chatsh

Zsh powered by ChatGPT. Language model observes the shell inputs and outputs, giving you useful informations, being your helpful assistant.

## Requirements

1. macOS.

   Linux and Windows systems are not tested for now, it may or may not work.

2. Go 1.19.6

   Lower versions are not tested, it may or may not work.

## Install

### Build from source

1. Clone this repo.

2. Run build script

   ```bash
   ./build.sh
   ```

3. Add `bin` directory to your `PATH` environment variable.

   ```bash
   export PATH="$PWD/bin:$PATH"
   ```

### Download binary

Todo.

## Todos

- [ ] History summary for unlimited token length
- [ ] Multiple chatsh instances
- [ ] Prebuilt binaries
- [ ] Linux & Windows support
- [ ] Supports models other than GPT-4