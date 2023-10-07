# StoreApp Documentation

StoreApp is a web application developed using GoLang that allows users to manage manufacturers, products, and transactions. It integrates with a MySQL database for data storage.

## Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
    - [Manufacturer Handling](#manufacturer-handling)
    - [Product Handling](#product-handling)
    - [Transaction Handling](#transaction-handling)

---

## 1. Installation <a name="installation"></a>

1. Set up the MySQL database by executing the provided SQL scripts (`queries.sql`).

2. Configure the MySQL connection details in the `main.go` file.

3. Install the required GoLang packages:

   ```bash
   go get github.com/go-sql-driver/mysql
   ```

4. Build and run the application:

   ```bash
   go run main.go
   ```

---

## 2. Usage <a name="usage"></a>

### 2.1 Manufacturer Handling <a name="manufacturer-handling"></a>

- **View Manufacturers:**
  - Access the manufacturer page at `http://localhost:8080/manufacturer`.

- **Add Manufacturer:**
  - Fill in the "Name" and "Address" fields and click "Add".

- **Edit Manufacturer:**
  - Click the "Edit" button next to a manufacturer, make changes, and click "Save Changes".

- **Delete Manufacturer:**
  - Click the "Delete" button next to a manufacturer to remove it.

### 2.2 Product Handling <a name="product-handling"></a>

- **View Products:**
  - Access the product page at `http://localhost:8080/product`.

- **Add Product:**
  - Fill in the "Name", "Quantity", and select a manufacturer. Click "Add".

- **Edit Product:**
  - Click the "Edit" button next to a product, make changes, and click "Save Changes".

- **Delete Product:**
  - Click the "Delete" button next to a product to remove it.

### 2.3 Transaction Handling <a name="transaction-handling"></a>

- **Create Transactions:**
  - Access the transaction page by adding products to cart and checking out.
