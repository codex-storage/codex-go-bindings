## v0.0.24 (2025-10-30)
### Notes

- Return the standard context.Canceled when a context is cancelled

## v0.0.23 (2025-10-30)
### Notes

- Add context cancellation support
- Prevent segmentation fault when trying to stop and node not started

## v0.0.22 (2025-10-20)
### Notes

- Downgrade Go version requirement to 1.24.0

## v0.0.21 (2025-10-15)
### Notes

- Remove libs/ from @rpath

## v0.0.20 (2025-10-15)
### Notes

- Set default install_name for mac

## v0.0.19 (2025-10-15)
### Notes

- Bump nim-codex

## v0.0.18 (2025-10-15)
### Notes

- Bump nim-codex to specific `install_name` for macOS

## v0.0.17 (2025-10-15)
### Notes

- Bump nim-codex to produce dylib for macos

## v0.0.16 (2025-10-15)
### Notes

- Remove the CGO LDFLAGS flags in the source code to control them with env variables

## v0.0.15 (2025-10-14)
### Notes

- Export fields in upload and download struct
- Fix typo

## v0.0.13 (2025-10-14)
### Notes

- Fix Rust version during build

## v0.0.11 (2025-10-13)
### Notes

- Fix libcodex.h path
- Rename CodexNew to New
- Rename CodexConfig to Config

## v0.0.10 (2025-10-13)
### Notes

- First release

### Features

- Codex data info
- Upload using `reader`, `file` and `chunks`
- Download using `stream` and `chunks`
- P2P connect
- Peer info and debug info