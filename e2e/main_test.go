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
	// cleanup := setupTestServer(t)
	// defer cleanup()

	// Тест 1: Создание таблицы
	t.Run("Create Table", func(t *testing.T) {
		query := "CREATE TABLE test_table (id INT, name TEXT, isAdmin BOOLEAN);"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)
	})

	// Тест 2: Вставка данных
	t.Run("Insert Data", func(t *testing.T) {
		queries := []string{
			"INSERT INTO test_table (id, name, isAdmin) VALUES (1, 'Rick', 'TRUE');",
			"INSERT INTO test_table (id, name, isAdmin) VALUES (2, 'Morty', 'FALSE');",
		}

		for _, query := range queries {
			response := executeQuery(t, query)
			assert.Empty(t, response.Error)
		}
	})

	// Тест 3: Выборка данных
	t.Run("Select Data", func(t *testing.T) {
		query := "SELECT id, name, isAdmin FROM test_table;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)

		want := `{"TableName":"test_table","Columns":[{"Name":"id","Type":"INT"},{"Name":"name","Type":"TEXT"},{"Name":"isAdmin","Type":"BOOLEAN"}],"Rows":[[1,"Rick","TRUE"],[2,"Morty","FALSE"]]}`

		assert.Equal(t, want, response.Result)
	})

	// Тест 4: Удаление таблицы
	t.Run("Drop Table", func(t *testing.T) {
		query := "DROP TABLE test_table;"
		response := executeQuery(t, query)
		assert.Empty(t, response.Error)
	})
}
