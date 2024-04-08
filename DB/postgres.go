package DB

import (
	"database/sql"
	// "encoding/json"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
	// "reflect"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "bb"
)

var db *sql.DB
var err error

func ConnectPsql() error {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println("[DB] Error connecting to postgres server.")
		return err
	}

	err = db.Ping()

	return err

}

func CreateTableUserDetails() error {
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS user_details (
		name VARCHAR(100) NOT NULL,
		gender VARCHAR(10) NOT NULL,
		age INTEGER NOT NULL,
		activity_level VARCHAR(20) NOT NULL,
		goals TEXT NOT NULL,
		height FLOAT NOT NULL,
		weight FLOAT NOT NULL,
		target_weight FLOAT NOT NULL,
		diseases TEXT NOT NULL,
		email VARCHAR(100) PRIMARY KEY,
		diet_plan TEXT,
		healthscore INTEGER NOT NULL,
		track JSONB,
		dm TEXT
	);	
	
	`)
	if err != nil {
		fmt.Println("[CREATE UD] Cant create.", err.Error())
		return err
	}
	fmt.Println("[UD] Created Table UserDetails")
	return nil
}

func InsertUserData(values map[string]interface{}) error {

	name, _ := values["name"].(string)
	gender, _ := values["gender"].(string)
	ageStr, _ := values["age"].(string)
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		fmt.Println("Conversion error:", err)
	}

	activityLevel, _ := values["activityLevel"].(string)
	goals, _ := values["goals"].([]string)

	heightStr, _ := values["height"].(string)
	height, err := strconv.ParseFloat(heightStr, 64)
	if err != nil {
		fmt.Println("Conversion error:", err)
	}
	weightStr, _ := values["weight"].(string)
	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		fmt.Println("Conversion error:", err)
	}
	targetWeightStr, _ := values["tweight"].(string)
	targetWeight, err := strconv.ParseFloat(targetWeightStr, 64)
	if err != nil {
		fmt.Println("Conversion error:", err)
	}
	// "tweight" key in the map
	diseases, _ := values["disease"].([]string) // "disease" key in the map
	email, _ := values["email"].(string)

	hsStr, _ := values["healthscore"].(string)
	hs, err := strconv.Atoi(hsStr)
	if err != nil {
		fmt.Println("Conversion error:", err)
	}

	// Convert arrays to comma-separated strings
	goalsStr := strings.Join(goals, ",")
	diseasesStr := strings.Join(diseases, ",")

	fmt.Println(name, gender, age, activityLevel, goalsStr, height, weight, targetWeight, diseasesStr, email)
	// Prepare the SQL query
	query := `
        INSERT INTO user_details (name, gender, age, activity_level, goals, height, weight, target_weight, diseases, email, healthscore)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `

	// Execute the SQL query
	_, err = db.Exec(query, name, gender, age, activityLevel, goalsStr, height, weight, targetWeight, diseasesStr, email, hs)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}
	fmt.Println("[UD] inserted data!")
	return nil
}

func ReadTrack(email string) ([]byte, error) {

	query := "SELECT track FROM user_details WHERE email = $1"
	var jsonData []byte
	err := db.QueryRow(query, email).Scan(&jsonData)
	if err != nil {
		fmt.Println("[Reading Track]", err)
		return nil, err
	}
	return jsonData, nil
}

func CheckEmailExists(email string) (bool, error) {

	query := "SELECT COUNT(*) FROM user_details WHERE email = $1"
	var count int
	err := db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func FetchExistingJSONData(email string) ([]map[string]interface{}, error) {

	var track map[string]interface{}
	var trackData []map[string]interface{}
	var trackDataString sql.NullString

	query := "SELECT track FROM user_details WHERE email = $1"
	err := db.QueryRow(query, email).Scan(&trackDataString)
	if err == fmt.Errorf("sql: no rows in result set") {
		return []map[string]interface{}{}, nil
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if trackDataString.Valid {
		err = json.Unmarshal([]byte(trackDataString.String), &trackData)
		if err != nil {
			err2 := json.Unmarshal([]byte(trackDataString.String), &track)
			if err2 != nil {
				fmt.Println(err2)
				return nil, fmt.Errorf("error parsing JSONB data: %v", err2)
			}
			trackData = []map[string]interface{}{track}
			return trackData, nil
		}
	} else {
		trackData = []map[string]interface{}{}
	}
	return trackData, nil
}

func AppendMeals(values map[string]interface{}) error {

	email, _ := values["email"].(string)
	delete(values, "email")

	existingMap, err := FetchExistingJSONData(email)
	if err != nil {
		return fmt.Errorf("failed to fetch existing JSON data: %w", err)
	}

	existingMap = append(existingMap, values)

	newData, err := json.Marshal(existingMap)
	if err != nil {
		return fmt.Errorf("failed to marshal merged data to JSON: %w", err)
	}

	query := `
        UPDATE user_details
        SET track = $1 
		WHERE email = $2
    `

	_, err = db.Exec(query, newData, email)
	if err != nil {
		return fmt.Errorf("failed to update JSON data: %w", err)
	}

	return nil
}

func ReadAllUsers(query string) ([]map[string]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to hold the results
	results := []map[string]interface{}{}

	// Iterate over the rows
	for rows.Next() {
		// Create a map to hold the values of the current row
		rowData := make(map[string]interface{})

		// Get column names
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		// Create a slice to store values for Scan
		values := make([]interface{}, len(columns))
		for i := range columns {
			values[i] = new(interface{})
		}

		// Scan the values of the current row into the map
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		// Convert interface{} values to appropriate types and store them in the map
		for i, col := range columns {
			rowData[col] = *(values[i].(*interface{}))
		}

		// Append the map to the results slice
		results = append(results, rowData)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func ReadRowData(primaryKeyValue string) (map[string]interface{}, error) {
	rowData := make(map[string]interface{})

	// Build SQL query
	query := fmt.Sprintf("SELECT * FROM user_details WHERE email = $1")

	// Execute the query and get the rows object
	rows, err := db.Query(query, primaryKeyValue)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting column names: %v", err)
	}

	// Create a slice of interface{} to hold column values
	values := make([]interface{}, len(columns))

	// Create pointers to each element in the values slice
	for i := range values {
		values[i] = new(interface{})
	}

	// Iterate over rows
	for rows.Next() {
		// Scan row data into the slice
		err = rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("error scanning row data: %v", err)
		}

		// Add data to the map
		for i, colName := range columns {
			rowData[colName] = *(values[i].(*interface{}))
		}
	}

	return rowData, nil
}

func UpdateDM(email string, message string) error {

	stmt, err := db.Prepare("UPDATE user_details SET dm = $1 WHERE email = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(message, email)
	if err != nil {
		return err
	}

	fmt.Println("Text field updated successfully!")
	return nil
}

func UpdateDiet(email string, diet string, healthscore int) error {
	var query string
	if healthscore == 0 {
		query = "UPDATE user_details SET diet_plan = $1 WHERE email = $2"
	} else {
		query = "UPDATE user_details SET diet_plan = $1, healthscore = $2 WHERE email = $3"
	}
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var err2 error
	if healthscore == 0 {
		_, err2 = stmt.Exec(diet, email)
	}else {
		_, err2 = stmt.Exec(diet, healthscore, email)
	}
	if err2 != nil {
		return err2
	}

	fmt.Println("Diet updated successfully!")
	return nil
}
