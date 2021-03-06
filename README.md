# hummus
A concise way to declare structs and marshal them into JSON in Go.

Making nested structs/arrays in order to marshal JSON from objects can be such a pain sometimes, can't it? Especially when you are cherry-picking fields from one flat JSON message and trying to output them into another complex JSON message.

Well, along comes hummus which makes this super-duper easy. The biggest win, I feel, is that you can have all your nesting information in the same struct. Creating nested arrays and/or objects is a simple matter of using dots (`.`) or square brackets (`[]`) in your tags.

## How to install

```
go get github.com/aditya87/hummus
```

#### Fetch dependencies

```
go get github.com/Jeffail/gabs
```

##### For testing
```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
```

## How to use

#### Marshalling JSON

```
...
type Info struct {
  Company           string `hummus:"company"`
	Address           string `hummus:"address"`
	Brand0Name        string `hummus:"brands[0].name"`
	Brand0Flavor      string `hummus:"brands[0].flavor"`
	Brand0Store0Name  string `hummus:"brands[0].stores[0].name"`
	Brand0Store0Price int    `hummus:"brands[0].stores[0].price,omitempty"`
	Brand0Store1Name  string `hummus:"brands[0].stores[1].name"`
	Brand0Store1Price int    `hummus:"brands[0].stores[1].price,omitempty"`
}

func main() {
  info := Info{
		Company:           "hello foods",
		Address:           "338 New St",
		Brand0Name:        "sabra",
		Brand0Flavor:      "jalapeno",
		Brand0Store0Name:  "safeway",
		Brand0Store0Price: 5,
		Brand0Store1Name:  "wholefoods",
		Brand0Store1Price: 10,
	}

	jsonOutput, err := hummus.Marshal(info)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonOutput))
}
```

Gives us:
```
{
  "company": "hello foods",
  "address": "338 New St",
  "brands": [
    {
      "flavor": "jalapeno",
      "name": "sabra",
      "stores": [
        {
          "name": "safeway",
          "price": 5
        },
        {
          "name": "wholefoods",
          "price": 10
        }
      ]
    }
  ]
}
```

#### Translating one type of message into another

```
type Info struct {
	Name   string `json:"name" hummus:"name"`
	Flavor string `json:"flavor" hummus:"type"`

	MainSupplierName     string `json:"main_supplier_name" hummus:"suppliers[0].name"`
	MainSupplierLocation string `json:"main_supplier_location" hummus:"suppliers[0].location"`

	BackupSupplierName     string `json:"backup_supplier_name" hummus:"suppliers[1].name"`
	BackupSupplierLocation string `json:"backup_supplier_location" hummus:"suppliers[1].location"`
}

func main() {
	inputJSON := `{
		"name": "sabra",
		"flavor": "jalapeno",
		"main_supplier_name": "Hipster Foods",
		"main_supplier_location": "CO",
		"backup_supplier_name": "Good Foods",
		"backup_supplier_location": "CA"
	}`

	var input Info
	err := json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		panic(err)
	}

	jsonOutput, err := hummus.Marshal(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonOutput))
}
```

Gives us:
```
{
  "name": "sabra",
  "type": "jalapeno",
  "suppliers": [
    {
      "location": "CO",
      "name": "Hipster Foods"
    },
    {
      "location": "CA",
      "name": "Good Foods"
    }
  ]
}
```

#### Special cases

##### Escaping dots

Sometimes, you want to not interpret the dots in a JSON tag as hierarchical keys, instead you just want them to get interpreted as a flat key like so:
```
{
  "properties.name": "value"
}
```

Hummus offers you the ability to escape dots. Simply replace dots with hashtags in the struct tags:
```
...
type S struct {
	V string `hummus:"properties#name"`
}

func main() {
	s := S{
		V: "value",
	}

	jsonOutput, err := hummus.Marshal(s)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonOutput))
}
```

This gives us:
```
{
  "properties.name": "value"
}
```

Which is the desired result.

## Notes

1. Also provided an `omitempty` option to ignore empty fields, just like the [encoding/json](https://golang.org/pkg/encoding/json/) library. e.g.:
```
type foo struct {
  bar string `hummus:"bar,omitempty"`
}
```
2. Leverages [reflect](https://golang.org/pkg/reflect/) for dynamic struct interpretation and [gabs](https://github.com/Jeffail/gabs) for dynamic JSON generation.

## Contributing

PRs are welcome. Make sure unit tests are run. To do so, firstly install the ginkgo and gomega libraries as described in the "fetch dependencies" section above. Then, simply run:

```
ginkgo -r .
```
In the main directory.
