# MAWSGO (mhafan AWS for Golang)

## Moduly

### Main

Predpoklada nastaveni sysenv AWS_DEFAULT_REGION. Z nejakeho duvodu to nebere z lokalne konfigurovaneho "aws configurate".

### S3

### DynamoDB

Modul je zalozen na jednoduchych zaznamech s jedinym klicovym (primarnim) atributem typu string.

#### Scan
```
var _out = []MawsgoTable{}

filt := expression.Name("Age").GreaterThan(expression.Value(10))

fmt.Println(tbl.Scan(filt, &_out))
fmt.Println(_out)
```

## Makefile

Vytvoreni zip-souboru s binarkou. Upload jako lambda.
```
zip:	main.go
	    GOOS=linux CGO_ENABLED=0 go build -o mainzip main.go
	    zip mainzip.zip mainzip
```