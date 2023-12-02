go build -o example.bin src/main.go
echo -n 'END_OF_PROGRAM' >> example.bin
cat /usr/bin/top >> example.bin
echo -n 'END_OF_PAYLOAD' >> example.bin
echo -n 'toor' >> example.bin