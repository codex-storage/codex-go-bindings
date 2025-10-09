# Codex Go Bindings

This repository provides Go bindings for the Codex library, enabling seamless integration with Go projects.

## Installation

To build the required dependencies for this module, the `make` command needs to be executed.
If you are integrating this module into another project via `go get`, ensure that you navigate
to the `codex-go-bindings` module directory and run the `make` commands.

### Steps to Install

Follow these steps to install and set up the module:

1. Make sure your system has the [prerequisites](https://github.com/codex-storage/nim-codex) to run a local Codex node.

2. Fetch the dependencies:
   ```
   make update
   ```

3. Build the library:
   ```
   make libcodex
   ```

You can pass flags to the Codex building step by using `CODEX_LIB_PARAMS`. For example,
if you want to enable debug API for peers, you can build the library using:

```
CODEX_LIB_PARAMS="-d:codex_enable_api_debug_peers=true" make libcodex
```

Now the module is ready for use in your project.

## Usage

### Init

First you need to create a Codex node:

```go
dataDir := "..."
node, err := CodexNew(CodexConfig{
   DataDir:        dataDir,
   BlockRetries:   10,
})
/// ....
err := node.Destroy()
```

The `CodexConfig` object provides several options to configure your node. You should at least
adjust the `DataDir` folder and the `BlockRetries` setting to avoid long retrieval times when
the data is unavailable.

When you are done with your node, you **have to** call `Destroy` method to free resources.

### Start / Stop

use `Start` method to start your node. You **have to** call `Stop` before `Destroy` when you are done
with your node.

```go
err := node.Start()
err := node.Stop()
err := node.Destroy()
```

### Info

You can get the version and revision without starting the node:

```go
version, err := node.Version()
revision, err := node.Revision()
```

Other information are available after the node is started:

```go
repo, err := node.Version()
spr, err := node.Spr()
peerId, err := node.PeerId()
```

### Upload

There are 3 strategies for uploading: `reader`, `file` or `chunks`. Each one requires its own upload session.

#### reader

The `reader` strategy is the easiest option when you already have a Go `Reader`.
It handles creating the upload session and cancels it if an error occurs.

The `filepath` should contain the data’s name with its extension, because Codex uses that to
infer the MIME type.

An `onProgress` callback is available to receive progress updates and notify the user.
The total size of the reader is determined via `stat` when it wraps a file, or from the buffer length otherwise.
From there, the callback can compute and report the percentage complete.

The `UploadReader` returns the cid of the content uploaded.

```go
buf := bytes.NewBuffer([]byte("Hello World!"))
onProgress := func(read, total int, percent float64, err error) {
   // Do something with the data
}
cid, err := codex.UploadReader(UploadOptions{filepath: "hello.txt", onProgress: onProgress}, buf)
```

Caveat: once started, the upload cannot be cancelled.

#### file

The `file` strategy allows you to upload a file on Codex using the path.
It handles creating the upload session and cancels it if an error occurs.

The `onProgress` callback is the same as for `reader` strategy.

The `UploadFile` returns the cid of the content uploaded.

```go
onProgress := func(read, total int, percent float64, err error) {
   // Do something with the data
}
cid, err := codex.UploadFile(UploadOptions{filepath: "./testdata/hello.txt", onProgress: onProgress})
```

Caveat: once started, the upload cannot be cancelled.

#### chunks

The `chunks` strategy allows to manage the upload by yourself. It requires more code
but provide more flexibility. You have to create the upload session, send the chunks
and then finalize to get the cid.

```go
sessionId, err := codex.UploadInit(&UploadOptions{filepath: "hello.txt"})

err = codex.UploadChunk(sessionId, []byte("Hello "))

err = codex.UploadChunk(sessionId, []byte("World!"))

cid, err := codex.UploadFinalize(sessionId)
```

Using this strategy, you can handle resumable uploads and cancel the upload
whenever you want !

### Download

When you receive a cid, you can download the `Manifest` to get information about the data:

```go
manifest, err := codex.DownloadManifest(cid)
```

It is not mandatory for downloading the data but it is really useful.

There are 2 strategies for downloading: `stream` and `chunks`.

#### stream

The `stream` strategy is the easiest to use.

It provides an `onProgress` callback to receive progress updates and notify the user.
The percentage is calculated from the `datasetSize` (taken from the manifest).
If you don’t provide it, you can enable `datasetSizeAuto` so `DownloadStream` fetches the
manifest first and uses its `datasetSize`.

You can pass a `write` callback and/or a `filepath` as destinations. They are not mutually exclusive,
letting you write the content to two places for the same download.

```go
opt := DownloadStreamOptions{
   writer:      f,
   datasetSize: len,
   filepath:    "testdata/hello.downloaded.writer.txt",
   onProgress: func(read, total int, percent float64, err error) {
      // Do something
   },
}
err := codex.DownloadStream(cid, opt)
```

Caveat: once started, the download cannot be cancelled.

#### chunks

The `chunks` strategy allows to manage the download by yourself. It requires more code
but provide more flexibility.

This strategy **assumes you already know the total size to download** (from the manifest).
After you believe all chunks have been retrieved, you **must** call `DownloadCancel`
to terminate the download session.

```go
cid := "..."
err := codex.DownloadInit(cid, DownloadInitOptions{})
chunk, err := codex.DownloadChunk(cid)
err := codex.DownloadCancel(cid)
```

Using this strategy, you can handle resumable downloads and cancel the download
whenever you want !

### Storage

Several methods are available to manage the data on your node:

```go
manifests, err := node.Manifests()
space, err := node.Space()

cid := "..."
err := node.Delete(cid)
err := node.Fetch(cid)
```

The `Fetch` method download remote data into your local node.

### P2P

You can connect to a node using the `peerId` or the `listenAddresses`:

```go
peerId := "..."
addrs := ["..."]
err := node.Connect(peerId, addrs)
```

### Debug

Several methods are available to debug your node:

```go
// Get node info
info, err := node.Debug()

// Update the chronicles level log on runtime
err := node.UpdateLogLevel("DEBUG")

peerId := "..."
record, err := node.CodexPeerDebug(peerId)
```

`CodexPeerDebug` is only available if you built with `-d:codex_enable_api_debug_peers=true` flag.

## Example

For an example on how to use this package, please take a look at our [example-go-bindings](https://github.com/codex-storage/example-codex-go-bindings) repo.