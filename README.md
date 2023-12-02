# Evilbin

Evilbin is a POC method of binding static binaries to a Go application, then executing them in memory while intercepting stdout.

## WARNING
Do not run this outside of your own controlled environment. This is purely for research & learning purposes.

## Run

See example `build.sh`

# How does it work?

Evilbin works by generating the binder

```console
go build -o example.bin src/main.go
```

Then appending any static binary to the end of the file

```console
echo -n 'END_OF_PROGRAM' >> example.bin
cat /usr/bin/top >> example.bin
```

And finally choosing a way to modify stdout
```console
echo -n 'END_OF_PAYLOAD' >> example.bin
echo -n 'toor' >> example.bin
```

This simply replaces any instances of root with toor in stdout.
