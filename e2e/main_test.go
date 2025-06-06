package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	serverURL = "http://localhost:8080"
)

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error,omitempty"`
}

func waitForServer(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url + "/query")
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("server did not start within %v", timeout)
}

func setupTestServer(t *testing.T) func() {
	// Запускаем сервер в отдельном процессе
	cmd := exec.Command("go", "run", "../cmd/main.go", "-mode", "http", "-port", "8080")

	// Перенаправляем вывод сервера в тест
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)
	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)

	err = cmd.Start()
	require.NoError(t, err)

	// Читаем вывод сервера
	go func() {
		io.Copy(io.Discard, stdout)
	}()
	go func() {
		io.Copy(io.Discard, stderr)
	}()

	// Ждем пока сервер запустится
	err = waitForServer(serverURL, 10*time.Second)
	require.NoError(t, err, "Server failed to start")

	// Возвращаем функцию для очистки
	return func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}
}

func executeQuery(t *testing.T, query string) *QueryResponse {
	reqBody := QueryRequest{Query: query}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	resp, err := http.Post(serverURL+"/query", "application/json", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	defer resp.Body.Close()

	var response QueryResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	return &response
}

func TestCreateAndQueryTable(t *testing.T) {
	t.Run("Create Table", func(t *testing.T) {
		query := "CREATE TABLE test_table (id INT, name TEXT, is_admin BOOLEAN, registered_at TIMESTAMP);"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)
	})

	t.Run("Insert Data", func(t *testing.T) {
		queries := []string{
			"INSERT INTO test_table VALUES (1, 'Rick', true, '2024-03-20 15:30:45');",
			"INSERT INTO test_table VALUES (2, 'Morty', false, '2023-05-01 13:30:45');",
		}

		for _, query := range queries {
			response := executeQuery(t, query)
			assert.Empty(t, response.Error)
		}
	})

	t.Run("Select Data", func(t *testing.T) {
		query := "SELECT id, name, is_admin, registered_at FROM test_table;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)

		want := `{"name":"test_table","columns":[{"name":"id","type":1},{"name":"name","type":0},{"name":"is_admin","type":2},{"name":"registered_at","type":3}],"rows":[[1,"Rick",true,"2024-03-20 15:30:45"],[2,"Morty",false,"2023-05-01 13:30:45"]]}`

		assert.Equal(t, want, response.Result)
	})

	t.Run("Select Data with only some columns", func(t *testing.T) {
		query := "SELECT id, name FROM test_table;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)

		want := `{"name":"test_table","columns":[{"name":"id","type":1},{"name":"name","type":0}],"rows":[[1,"Rick"],[2,"Morty"]]}`

		assert.Equal(t, want, response.Result)
	})

	t.Run("Select Data with *", func(t *testing.T) {
		query := "SELECT * FROM test_table;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)

		want := `{"name":"test_table","columns":[{"name":"id","type":1},{"name":"name","type":0},{"name":"is_admin","type":2},{"name":"registered_at","type":3}],"rows":[[1,"Rick",true,"2024-03-20 15:30:45"],[2,"Morty",false,"2023-05-01 13:30:45"]]}`

		assert.Equal(t, want, response.Result)
	})

	t.Run("Drop Table", func(t *testing.T) {
		query := "DROP TABLE test_table;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)
	})
}

func TestWhereClause(t *testing.T) {
	// Подготовка данных
	t.Run("Setup Test Data", func(t *testing.T) {
		queries := []string{
			"CREATE TABLE test_table_2 (id INT, name TEXT, age INT, is_admin BOOLEAN);",
			"INSERT INTO test_table_2 VALUES (1, 'Alice', 25, true);",
			"INSERT INTO test_table_2 VALUES (2, 'Bob', 30, false);",
			"INSERT INTO test_table_2 VALUES (3, 'Charlie', 35, true);",
			"INSERT INTO test_table_2 VALUES (4, 'David', 40, false);",
		}

		for _, query := range queries {
			response := executeQuery(t, query)
			assert.Empty(t, response.Error)
		}
	})

	// Тест простого условия равенства
	t.Run("Simple Equality", func(t *testing.T) {
		query := "SELECT id, name, age, is_admin FROM test_table_2 WHERE id = 2;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)
		want := `{"name":"test_table_2","columns":[{"name":"id","type":1},{"name":"name","type":0},{"name":"age","type":1},{"name":"is_admin","type":2}],"rows":[[2,"Bob",30,false]]}`
		assert.Equal(t, want, response.Result)
	})

	// Тест условия неравенства
	t.Run("Inequality", func(t *testing.T) {
		query := "SELECT id, name, age, is_admin FROM test_table_2 WHERE age > 30;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)
		want := `{"name":"test_table_2","columns":[{"name":"id","type":1},{"name":"name","type":0},{"name":"age","type":1},{"name":"is_admin","type":2}],"rows":[[3,"Charlie",35,true],[4,"David",40,false]]}`
		assert.Equal(t, want, response.Result)
	})

	// Тест условия с AND
	t.Run("AND Condition", func(t *testing.T) {
		query := "SELECT id, name, age, is_admin FROM test_table_2 WHERE age > 25 AND age < 40;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)
		want := `{"name":"test_table_2","columns":[{"name":"id","type":1},{"name":"name","type":0},{"name":"age","type":1},{"name":"is_admin","type":2}],"rows":[[2,"Bob",30,false],[3,"Charlie",35,true]]}`
		assert.Equal(t, want, response.Result)
	})

	// Тест условия с LIMIT
	t.Run("LIMIT Condition", func(t *testing.T) {
		query := "SELECT id, name, age, is_admin FROM test_table_2 LIMIT 1;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)
		want := `{"name":"test_table_2","columns":[{"name":"id","type":1},{"name":"name","type":0},{"name":"age","type":1},{"name":"is_admin","type":2}],"rows":[[1,"Alice",25,true]]}`
		assert.Equal(t, want, response.Result)
	})

	// Тест условия с OFFSET
	t.Run("OFFSET Condition", func(t *testing.T) {
		query := "SELECT id, name, age, is_admin FROM test_table_2 OFFSET 1;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)
		want := `{"name":"test_table_2","columns":[{"name":"id","type":1},{"name":"name","type":0},{"name":"age","type":1},{"name":"is_admin","type":2}],"rows":[[2,"Bob",30,false],[3,"Charlie",35,true],[4,"David",40,false]]}`
		assert.Equal(t, want, response.Result)
	})

	// Тест условия с OR
	t.Run("OR Condition", func(t *testing.T) {
		query := "SELECT id, name, age, is_admin FROM test_table_2 WHERE age = 25 OR age = 40;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)
		want := `{"name":"test_table_2","columns":[{"name":"id","type":1},{"name":"name","type":0},{"name":"age","type":1},{"name":"is_admin","type":2}],"rows":[[1,"Alice",25,true],[4,"David",40,false]]}`
		assert.Equal(t, want, response.Result)
	})

	// Тест сложного условия с AND и OR
	t.Run("Complex AND OR Condition", func(t *testing.T) {
		query := "SELECT id, name, age, is_admin FROM test_table_2 WHERE (age > 30 AND name = 'Charlie') OR (age < 30 AND name = 'Alice');"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)
		want := `{"name":"test_table_2","columns":[{"name":"id","type":1},{"name":"name","type":0},{"name":"age","type":1},{"name":"is_admin","type":2}],"rows":[[1,"Alice",25,true],[3,"Charlie",35,true]]}`
		assert.Equal(t, want, response.Result)
	})

	// Очистка
	t.Run("Cleanup", func(t *testing.T) {
		query := "DROP TABLE test_table_2;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)
	})
}
