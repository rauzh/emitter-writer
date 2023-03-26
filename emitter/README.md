{
    "id": 123,
    "ASIN": "0827229534",
    "discontinued": false,
    "title": "BlahBlah",
    "group": "Book",
    "salesrank": 123456,
    "similarnum": 5,
    "similars": [
        "0804675715", "0845615715", "080427685", "0804211715", "4675215715"
    ]
}

// discontinued == true -> пустое тело, есть только id и ASIN (пример - первая запись в базе)

```
{
    "id": 123,
    "ASIN": "0827229534",
    "discontinued": false,
    "title": "BlahBlah",
    "group": "Book",
    "salesrank": 123456,
    "similarnum": 5,
    "similars": [
        "0804675715", "0845615715", "080427685", "0804211715", "4675215715"
    ]
}
```


```
type Item struct {
	IsDiscontinued bool
	ID             uint
	ASIN           string
	Title          string
	Group          string
	Salesrank      uint64
	Similarsnum    uint
	Similars       []string
}
```
