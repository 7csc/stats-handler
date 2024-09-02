# Go System Stats API  
  
This is a simple API Handler written in Go that provides the following system statistics:  
  
- Memory usage
- Number of goroutines
- File descriptor usage
- Environment variables
  
## Features  
  
- **Memory Usage**: Provides memory a allocation details and usage percentage.    
- **Goroutine Count**: Displays the current number of goroutines.  
- **File Descriptor Usage**: Shows the number of currently open file descriptors.  
- **Environment Variables**: Returns all environment variables in the process.  
- **Uptime**: Shows the server uptime since the process started.  
  
## Getting Started  
### Prerequisites  
  
- Go 1.18 or higher
- A Unix-Based System (Linux or MacOS) is recommendedfor accurayte file descriptor count.  
  
### Installation  
  
1. Clone the repository:  
  
```bash
git clone https://github.com/7csc/stats-handler.git
cd stats-handler
```
  
2. BUild the project:
  
```bash
go build -o stats-handler
```
  
```bash
./stats-handler
```  
  
## Usage  
This API Handler is designed to be integrated with a router like "chi-router". Here is an example of how to use it in your main application.  
  
```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "your_module_path/stats"
)

func main() {
    r := chi.NewRouter()

    // Your application routes
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    // Stats route
    r.Get("/stats", stats.Handler)

    http.ListenAndServe(":8080", r)
}

```
  
Once the server is running, you can access the API endpoint by sending a GET request to `http://(localhost):8080/stats`.  
  
Note: When integrating the stats handler into your main server, ensure that the main server's port configuration 

## API Response  
The API will return a JSON response with the following structure:  
  
```json
{
    "time": 1625503378245678900,
    "go_version": "go1.18",
    "go_os": "linux",
    "go_arch": "amd64",
    "cpus": 4,
    "goroutine_num": 12,
    "memory_alloc": 945256,
    "memory_total_alloc": 1024000,
    "memory_sys": 5079040,
    "memory_usage_percent": 18.61,
    "file_descriptor_num": 32,
    "env_vars": {
        "PATH": "/usr/local/bin:/usr/bin:/bin",
        "HOME": "/home/user"
    },
    "uptime": 3600
}  
```
 
## Contributing  
Contributions are welcome! Please submit a pull request or open an issue if you have suggestions or improvements.  
  
## License  
Tis project is licensed under the MIT License. See the LICENSE file for details.  