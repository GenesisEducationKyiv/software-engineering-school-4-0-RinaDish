package main

import (
	"fmt"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

func main() {
    rate := vegeta.Rate{Freq: 50, Per: time.Second} 
    duration := 10 * time.Second

    targeter := vegeta.NewStaticTargeter(vegeta.Target{
        Method: "POST",
        URL:    "http://localhost:8183/subscribe",
        Body:   []byte("email=test@test.com"),
        Header: map[string][]string{
            "Content-Type": {"application/x-www-form-urlencoded"},
        },
    })
    attacker := vegeta.NewAttacker()

    f, err := os.Create("results.bin")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    encoder := vegeta.NewEncoder(f)

    var metrics vegeta.Metrics
    for res := range attacker.Attack(targeter, rate, duration, "Load test") {
        metrics.Add(res)

        if err := encoder.Encode(res); err != nil {
            panic(err)
        }
    }
    metrics.Close()

    fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}
