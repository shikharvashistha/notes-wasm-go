// main.go
// Automatically handled with Promise rejects when returning an error!
func divide(x int, y int) (int, error) {
    if y == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return x / y, nil
}

func main() {
    wasm.Expose("divide", divide)
    wasm.Ready()
}
