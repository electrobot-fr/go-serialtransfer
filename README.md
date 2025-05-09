# go-serialtransfer 

Golang library for https://github.com/PowerBroker2/SerialTransfer

## Usage

Encode

```go
type message struct {
    X int16 `struc:"int16,little"`
    Y int16 `struc:"int16,little"`
	
    CmdAimantInt     bool `struc:"bool"`
    AscBoites        int  `struc:"int16,little"`
    Compteur         int  `struc:"uint8,little"`
}

serialtransfer.Encode(message{
    X: 1,
    Y: 2,
})
```

Decode

```go
type message struct {
    X int16 `struc:"int16,little"`
    Y int16 `struc:"int16,little"`
	
    CmdAimantInt     bool `struc:"bool"`
    AscBoites        int  `struc:"int16,little"`
    Compteur         int  `struc:"uint8,little"`
}

var msg message
serialtransfer.Decode(bin, &msg)
```