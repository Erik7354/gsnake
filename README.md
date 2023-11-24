# gsnake

[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)
[![Generic badge](https://img.shields.io/badge/Made_with-HTMX-3465a4.svg)](https://htmx.org)
[![Generic badge](https://img.shields.io/badge/Memes%3F-yes-<COLOR>.svg)]()


Snake using hypermedia as the engine of application state.

**Powered by Go+HTMX.**

## How To Start

```bash
make run
```

[Open localhost:8000](http://localhost:8000)

Game settings are done by query parameters:

- n: number of rows/columns; default/min: 5
- rr: refresh rate; default: 1000ms; [more formats](https://pkg.go.dev/time#ParseDuration)

## TODO

- [x] fix: apple respawns on start (Chrome)
- [x] highlight head of snake
- [ ] delete sessions
- [ ] improve ui 
- [ ] high scores