# Bober - Meaning Beaver in Polish
A silly YAML to Makefile CLI tool for my C++ needs.

Example project.yaml that is required for `bober` to work.
```yaml
# project.yaml
project:
  name: bobus
  version: 1.0.0

cpp:
  compiler: clang++
  standard: c++23
  flags: -Wall -Wextra -pedantic -O3

sources:
  - src

output:
  executable: bobus

libraries:
#  - name: raylib -- You can only have one, you either use pkg-config or name. name field will pass the library name to Makefile as -lraylib in this example.
  - config: pkg-config raylib
```
Example project structure:

```
project.yaml
src -> main.cpp
```
In above example, src can be named anything, just make sure to adjust this in `project.yaml` `sources:` field.

Bober has the following commands:

`bober` -> Generates the Makefile

`bober build` -> Builds the project

`bober clean` -> Cleans the build directory (deletes it)

`bober run` -> Runs the project

`bober run --html5` Starts a local HTTP server to serve files from the build directory. Useful when building web games.

I don't expect anyone to ever use this or have a sane reason to do so.

This is a tool I've made for my own purposes.

