# TODO list
`flag` package to get cmd line options, built-in `-h` for help, can customize by setting `flag.Usage` to a customized function (`flag.PrintDefaults()` for including the defaults as part of the cutomized help message.)

Separation of core logic ad cmd line interaction: 
- `package todo` is a lib that only includes the behaviour of the list
- `cmd/todo` has the logic for how this logic is used for user iteraction 

`fmt.Stringer` iterface: implement a `String()` function to inherit & controls how a type is `fmt` printed. 

Use environment variables instead of hardcoding to provide more flexibility

Take interface instead of concrete types as function args for flexibility: decouple implementatio from specific types