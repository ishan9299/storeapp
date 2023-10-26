package main

import (
    "fmt"
    "log"
    "strings"
    // "time"

    "net/http"
    "html/template"

    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type Manufacturer struct {
    Id int
    Name string
    Address string
}

type ManufacturerPageData struct {
    Manufacturers []Manufacturer
}

type Product struct {
    Id int
    Name string
    // ManufacId int
    Quantity int
    ManufacName string
}

type ProductPageData struct {
    Products []Product
    Manufacturers []Manufacturer
    PurchaseList []string
}

type Transaction struct {
    Name string
    Quantity int
}

type TransactionPageData struct {
    Transactions []Transaction
}

var tmpl *template.Template
var context ManufacturerPageData

var productTmpl *template.Template
var productContext ProductPageData

var transactionTmpl *template.Template
var transactionContext TransactionPageData

var manufacNameIdPair map[string]int

func main() {
    db, err = sql.Open("mysql", "ishan:password@tcp(127.0.0.1:3306)/storeapp")

    if err != nil {
        fmt.Println("error occured!")
        panic(err.Error())
    }

    defer db.Close()

    PORT := ":8080"

    http.HandleFunc("/manufacturer", ManufacturerHandler)
    http.HandleFunc("/manufacadd/", ManufacturerAddHandler)
    http.HandleFunc("/manufacedit/", ManufacturerEditHandler)
    http.HandleFunc("/manufacdelete/", ManufacturerDeletionHandler)
    http.HandleFunc("/product", ProductHandler)
    http.HandleFunc("/productdelete/", ProductDeleteHandler)
    http.HandleFunc("/productedit/", ProductEditHandler)
    http.HandleFunc("/productadd/", ProductAddHandler)
    http.HandleFunc("/transactions/", TransactionHandler)

    log.Fatal(http.ListenAndServe(PORT, nil))
}

func ManufacturerHandler(w http.ResponseWriter, r *http.Request) {
    var manufacturers []Manufacturer
    tmpl, _ = template.ParseFiles("./public/manufacturer.tmpl", "./public/bootstrap.tmpl")

    if r.Method == "GET" {

        sqlGetManufacturer(&manufacturers)

        context = ManufacturerPageData{Manufacturers: manufacturers}
        tmpl.Execute(w, context)

    }
}

func sqlGetManufacturer(manufacturers *[]Manufacturer) {
    var Id int
    var Name string
    var Address string

    manufacNameIdPair = make(map[string]int)

    getManufacturer := "select * from manufacturer;"
    rows, err := db.Query(getManufacturer)

    if err != nil {
        panic(err.Error())
    }

    defer rows.Close()

    for rows.Next() {
        rows.Scan(&Id, &Name, &Address)

        if err != nil {
            panic(err.Error())
        }

        manufacturer := Manufacturer {
            Id: Id,
            Name: Name,
            Address: Address,
        }

        manufacNameIdPair[Name] = Id

        *manufacturers = append(*manufacturers, manufacturer)
    }
}

func ManufacturerEditHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        var m Manufacturer
        r.ParseForm()
        id := r.FormValue("idvalue")

        getManufacturer := "select * from manufacturer where id=?"

        row := db.QueryRow(getManufacturer, id)

        err = row.Scan(&m.Id, &m.Name, &m.Address)
        if err != nil {
            log.Fatal(err)
        }

        manufacEditTmpl, err := template.ParseFiles("./public/manufacedit.tmpl", "./public/bootstrap.tmpl")
        if err != nil {
            log.Fatal(err)
        }

        manufacEditTmpl.Execute(w, m)
    } else if r.Method == "POST" {
        r.ParseForm()
        Id := r.FormValue("edit-id")
        Name := r.FormValue("edit-name")
        Address := r.FormValue("edit-address")

        stmt, err := db.Prepare(`update manufacturer set name=?, address=? where id=?;`)

        _, err = stmt.Exec(Name, Address, Id)

        if err != nil {
            panic(err.Error())
        }

        http.Redirect(w, r, "/manufacturer", http.StatusFound)
        tmpl.Execute(w, context)
    }
}

func ManufacturerDeletionHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    id := r.FormValue("idvalue")
    stmt1, err := db.Prepare(`delete from products where manufacid=?`)

    if err != nil {
        panic(err.Error())
    }

    _, err = stmt1.Exec(id)

	if err != nil {
		panic(err.Error())
	}

    stmt2, err := db.Prepare(`delete from manufacturer where id=?`)

    if err != nil {
        panic(err.Error())
    }

    _, err = stmt2.Exec(id)

	if err != nil {
		panic(err.Error())
	}

    http.Redirect(w, r, "/manufacturer", http.StatusFound)
}

func ManufacturerAddHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        r.ParseForm()
        name := r.FormValue("add-name")
        manufacturer := r.FormValue("add-address")

        if name != "" || manufacturer != "" {
            AddManufacturer(name, manufacturer)
        }

        http.Redirect(w, r, "/manufacturer", http.StatusFound)
    } else if r.Method == "GET" {
        manufacAddTmpl, err := template.ParseFiles("./public/manufacadd.tmpl", "./public/bootstrap.tmpl")
        if err != nil {
            log.Fatal(err)
        }

        manufacAddTmpl.Execute(w, nil)
    }
}

func AddManufacturer(name, address string) error {
    stmt, err := db.Prepare(`insert into manufacturer (name, address) values(?, ?);`)

    if err != nil {
        panic(err.Error())
    }

    _, err = stmt.Exec(name, address)

    // fmt.Println("Adding")
    // fmt.Println(err)
	if err != nil {
		panic(err.Error())
	}
    return err
}

func ProductHandler(w http.ResponseWriter, r *http.Request) {
    productTmpl, err = template.ParseFiles("./public/products.tmpl", "./public/bootstrap.tmpl")
    if err != nil {
        log.Fatal(err)
    }
    var products []Product
    var manufacturers []Manufacturer
    if r.Method == "GET" {

        sqlGetManufacturer(&manufacturers)
        sqlGetProducts(&products)

        productContext = ProductPageData{Products: products, Manufacturers: manufacturers}

        productTmpl.Execute(w, productContext)
    }
}

func ProductDeleteHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    id := r.FormValue("idvalue")
    stmt, err := db.Prepare(`delete from products where id=?`)

    if err != nil {
        panic(err.Error())
    }

    _, err = stmt.Exec(id)

	if err != nil {
		panic(err.Error())
	}

    http.Redirect(w, r, "/product", http.StatusFound)
}

func ProductEditHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    Id := r.FormValue("idvalue")
    var m Product
    var products []Product
    var manufacturers []Manufacturer

    if r.Method == "GET" {

        productEditTmpl, err := template.ParseFiles("./public/productedit.tmpl", "./public/bootstrap.tmpl")
        if err != nil {
            log.Fatal(err)
        }

        getProduct := "select * from products where id=?"

        row := db.QueryRow(getProduct, Id)

        err = row.Scan(&m.Id, &m.Name, &m.Quantity, &m.ManufacName)
        if err != nil {
            log.Fatal(err)
        }

        products = append(products, m)
        sqlGetManufacturer(&manufacturers)

        productContext = ProductPageData{Products: products, Manufacturers: manufacturers}

        productEditTmpl.Execute(w, productContext)

    } else if r.Method == "POST" {
        r.ParseForm()

        Name := r.FormValue("edit-name")
        Quantity := r.FormValue("edit-quantity")
        Manufacturer := r.FormValue("edit-manufacturer")

        ManufacId := manufacNameIdPair[Manufacturer]
        fmt.Println("Manufacturer: ", Manufacturer, "manufacNameIdPair: ", manufacNameIdPair)

        stmt, err := db.Prepare(`update products set name=?, quantity=?, manufacid=? where id=?;`)

        fmt.Println("Name: ", Name, "Quantity: ", Quantity, "ManufacId: ", ManufacId, "Id: ", Id)

        _, err = stmt.Exec(Name, Quantity, ManufacId, Id)


        if err != nil {
            panic(err.Error())
        }

        http.Redirect(w, r, "/product", http.StatusFound)
        productTmpl.Execute(w, productContext)
    }
}

func ProductAddHandler(w http.ResponseWriter, r *http.Request) {
    var manufacturers []Manufacturer
    if r.Method == "GET" {
        productAddTmpl, err := template.ParseFiles("./public/productadd.tmpl", "./public/bootstrap.tmpl")
        if err != nil {
            log.Fatal(err)
        }

        sqlGetManufacturer(&manufacturers)
        context = ManufacturerPageData{Manufacturers: manufacturers}

        productAddTmpl.Execute(w, context)
    } else if r.Method == "POST" {
        r.ParseForm()

        Name := r.FormValue("add-name")
        Quantity := r.FormValue("add-quantity")
        Manufacturer := r.FormValue("add-manufacturer")

        ManufacId := manufacNameIdPair[Manufacturer]

        stmt, err := db.Prepare(`insert into products (name, quantity, manufacid) values(?, ?, ?);`)

        if err != nil {
            panic(err.Error())
        }

        if Name != "" || Quantity != "" || Manufacturer != "" {
            _, err = stmt.Exec(Name, Quantity, ManufacId)
            if err != nil {
                panic(err.Error())
            }
        }


        http.Redirect(w, r, "/product", http.StatusFound)
        productTmpl.Execute(w, productContext)
    }
}

func sqlGetProducts(products *[]Product) {
    var Id int
    var Name string
    var Quantity int
    var ManufacName string

    getProduct := `
    select
    p.id, p.name as product_name,
    p.quantity, m.name
    from products p
    inner join manufacturer m
    on p.manufacid = m.id;`

    rows, err := db.Query(getProduct)

    if err != nil {
        panic(err.Error())
    }

    defer rows.Close()

    for rows.Next() {
        rows.Scan(&Id, &Name, &Quantity, &ManufacName)

        if err != nil {
            panic(err.Error())
        }

        product := Product {
            Id: Id,
            Name: Name,
            Quantity: Quantity,
            ManufacName: ManufacName,
        }

        *products = append(*products, product)
    }
}

func sqlGetSelectedProducts(products *[]Product, names []string) {
    var Id int
    var Name string
    var Quantity int
    var ManufacName string

    getProduct := `
    select
    p.id, p.name as product_name,
    p.quantity, m.name
    from products p
    inner join manufacturer m
    on p.manufacid = m.id
    where m.name in ();`

    rows, err := db.Query(getProduct)

    if err != nil {
        panic(err.Error())
    }

    defer rows.Close()

    for rows.Next() {
        rows.Scan(&Id, &Name, &Quantity, &ManufacName)

        if err != nil {
            panic(err.Error())
        }

        product := Product {
            Id: Id,
            Name: Name,
            Quantity: Quantity,
            ManufacName: ManufacName,
        }

        *products = append(*products, product)
    }
}

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
    var transactions []Transaction
    transactionTmpl, _ = template.ParseFiles("./public/transactions.tmpl", "./public/bootstrap.tmpl")

    r.ParseForm()
    names := r.FormValue("names")
    nameArray := strings.Split(names, ",")

    for _, name := range nameArray {
        transaction := Transaction {
            Name: name,
            Quantity: 1,
        }

        transactions = append(transactions, transaction)
    }

    transactionContext = TransactionPageData{Transactions: transactions}

    if r.Method == "GET" {
        transactionTmpl.Execute(w, transactionContext)
    }
}
