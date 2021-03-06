# loraserver-provisioner
A simple device provisioner. It expects a csv file `path` with the devices, and `hostname`, `username` and `password` as flags.

This is an example of the expected csv format:

| dev\_eui         | name    | application\_id | description            | device\_profile\_id                      | skip\_f\_cnt\_check | reference\_altitude | dev\_addr | nwk\_key                         | app\_key                         | gen\_app\_key | app\_s\_key                      | f\_nwk\_s\_int\_key              | s\_nwk\_s\_int\_key              | nwk\_s\_enc\_key                 | activation |
|------------------|---------|-----------------|------------------------|------------------------------------------|---------------------|---------------------|-----------|----------------------------------|----------------------------------|---------------|----------------------------------|----------------------------------|----------------------------------|----------------------------------|------------|
| 0000000000000001 | device1 | 1               | device\-description\-1 | 994b28a5\-cc81\-4a40\-8f23\-71030db4b38e | true                | 600\.0              |           | 00000000000000010000000000000001 | b06a309cb576cc82a607f6339609f25f |               |                                  |                                  |                                  |                                  | OTAA       |
| 0000000000000002 | device2 | 1               | device\-description\-2 | 2b4fb8e1\-3fa4\-497a\-b3bb\-fd1eca4727c6 | true                | 600\.0              | 00000001  |                                  |                                  |               | fe37b2fb6aa30900c04937944297817b | f99477469164b614343cf3581db64baa | f99477469164b614343cf3581db64baa | f99477469164b614343cf3581db64baa | ABP        |

## Building

Make sure you have [Go](https://golang.org/) installed and then just clone the repo and build from the `cli` directory:

```
git clone https://github.com/iegomez/lsp.git
cd lsp/cli
go build
```

Now you can run the program like this:

```
./cli --path /path/to/file.csv --hostname https://example.com --username your-user --password your-password
```


