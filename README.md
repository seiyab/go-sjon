# go-sjon
go-sjon provides JSON utilities. As of now, it only provides customizable JSON serializer.

## Serializer
Serializer is a customizable json serializer.

**Example: print struct key in lower camel case**
```go
type User struct {
    FirstName string
    LastName string
}

s := sjon.NewSerializer().
    With(sjon.StructKeyNamer(strcase.ToLowerCamel))
    
b, _ := s.Marshal(User{
    FirstName: "Betty",
    LastName: "Miller",
})
fmt.Println(string(b))
// Output: {"firstName":"Betty","lastName":"Miller"}
```

**Example: print only date for time.Time**
```go
s := sjon.NewSerializer().
    With(sjon.Replacer(func(d time.Time) string {
        return d.Format("2006-01-02")
    }))
    
b, _ := s.Marshal([]time.Time{
        time.Date(2024, 10, 20, 1, 2, 3, 4, time.UTC),
        time.Date(2024, 10, 21, 1, 2, 3, 4, time.UTC),
        time.Date(2024, 10, 22, 1, 2, 3, 4, time.UTC),
    })
fmt.Println(string(b))
// Output: ["2024-10-20","2024-10-21","2024-10-22"]
```
