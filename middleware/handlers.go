package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"go_learn/models" // models package where Product schema is defined
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"  // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// CreateProduct create a product in the postgres db
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty product of type models.Product
	var product models.Product

	// decode the json request to product
	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert product function and pass the product
	insertID := insertProduct(product)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "Product created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// GetProduct will return a single product by its id
func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get the id from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getProduct function with product id to retrieve a single product
	product, err := getProduct(int64(id))

	if err != nil {
		log.Fatalf("Unable to get product. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(product)
}

// GetAllProduct will return all the products
func GetAllProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the products in the db
	products, err := getAllProducts()

	if err != nil {
		log.Fatalf("Unable to get all product. %v", err)
	}

	// send all the products as response
	json.NewEncoder(w).Encode(products)
}

// UpdateProduct update product's detail in the postgres db
func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the id from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// create an empty product of type models.Product
	var product models.Product

	// decode the json request to product
	err = json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update product to update the product
	updatedRows := updateProduct(int64(id), product)

	// format the message string
	msg := fmt.Sprintf("Product updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// DeleteProduct delete product's detail in the postgres db
func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the id from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the deleteProduct, convert the int to int64
	deletedRows := deleteProduct(int64(id))

	// format the message string
	msg := fmt.Sprintf("Product updated successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------
// insert one product in the DB
func insertProduct(product models.Product) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning id will return the id of the inserted product
	sqlStatement := `INSERT INTO products (pname, pdesc, mrp, stBidPrice) VALUES ($1, $2, $3, $4) RETURNING id`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, product.Pname, product.Pdesc, product.Mrp, product.StBidPrice).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// get one product from the DB by its id
func getProduct(id int64) (models.Product, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a product of models.Product type
	var product models.Product

	// create the select sql query
	sqlStatement := `SELECT * FROM products WHERE id=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to product
	err := row.Scan(&product.ID, &product.Pname, &product.Pdesc, &product.Mrp, &product.StBidPrice)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return product, nil
	case nil:
		return product, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty product on error
	return product, err
}

// get one product from the DB by its id
func getAllProducts() ([]models.Product, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var products []models.Product

	// create the select sql query
	sqlStatement := `SELECT * FROM products`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var product models.Product

		// unmarshal the row object to product
		err = rows.Scan(&product.ID, &product.Pname, &product.Pdesc, &product.Mrp, &product.StBidPrice)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the product in the products slice
		products = append(products, product)

	}

	// return empty product on error
	return products, err
}

// update product in the DB
func updateProduct(id int64, product models.Product) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE products SET pname=$2, pdesc=$3, mrp=$4, stBidPrice=$5 WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, product.Pname, product.Pdesc, product.Mrp, product.StBidPrice)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete product in the DB
func deleteProduct(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM products WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
