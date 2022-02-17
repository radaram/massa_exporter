package main


import (
    "flag"
    "log"
    "net/http"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "fmt"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)


type WalletInfo struct {
    finalBalance float32
    candidateBalance float32
    lockedBalance float32
    activeRolls int32
    finalRolls int32
    candidateRolls int32
}


func main() {
    var addr = flag.String(
        "port", ":8875", "The port to listen on for HTTP requests.",
    )
    var dir = flag.String(
        "dir", "", "The path to the directory with the command",
    )
    flag.Parse()

    if *dir == "" {
        log.Printf("dis is required param\n")
        os.Exit(1)
    }

    walletInfoRaw := runWalletInfo(*dir)
    walletInfo := extractWalletInfoFromRawData(walletInfoRaw)

    finalBalance := prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "final_balance",
        },
    )
    candidateBalance := prometheus.NewGauge(
       prometheus.GaugeOpts{
           Name: "candidate_balance",
       },
    )
    lockedBalance := prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "locked_balance",
        },
    )

    activeRolls := prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_rolls",
        },
    )
    finalRolls := prometheus.NewGauge(
       prometheus.GaugeOpts{
           Name: "final_rolls",
       },
    )
    candidateRolls := prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "candidate_rolls",
        },
    )

    prometheus.MustRegister(finalBalance)
    finalBalance.Set(float64(walletInfo.finalBalance))

    prometheus.MustRegister(candidateBalance)
    candidateBalance.Set(float64(walletInfo.candidateBalance))

    prometheus.MustRegister(lockedBalance)
    lockedBalance.Set(float64(walletInfo.lockedBalance))

    prometheus.MustRegister(activeRolls)
    activeRolls.Set(float64(walletInfo.activeRolls))

    prometheus.MustRegister(finalRolls)
    finalRolls.Set(float64(walletInfo.finalRolls))

    prometheus.MustRegister(candidateRolls)
    candidateRolls.Set(float64(walletInfo.candidateRolls))

    http.Handle("/metrics", promhttp.Handler())

    log.Printf("Starting web server at %s\n", *addr)
    err := http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Printf("http.ListenAndServer: %v\n", err)
    }
}


func runWalletInfo(dirPath string) string {
    log.Printf("Path: %s\n", dirPath)
    cmd := exec.Command("./massa-client", "wallet_info")
    cmd.Dir = dirPath
    out, err := cmd.Output()
    if err != nil {
        log.Fatal(err)
    } else {
        log.Println(string(out))
    }

    return string(out)
}


func extractWalletInfoFromRawData(rawData string) WalletInfo {
    var finalBalance float32
    var candidateBalance float32
    var lockedBalance float32
    var activeRolls int32
    var finalRolls int32
    var candidateRolls int32

    re := regexp.MustCompile("Final balance: \\d+\\.\\d+")
    match := re.FindAllStringSubmatch(rawData, 1)
    for i := range match {
        value := strings.Replace(match[i][0], "Final balance: ", "", 1)
        fmt.Sscan(value, &finalBalance)
    }

    re = regexp.MustCompile("Candidate balance: \\d+\\.\\d+")
    match = re.FindAllStringSubmatch(rawData, 1)
    for i := range match {
        value := strings.Replace(match[i][0], "Candidate balance: ", "", 1)
        fmt.Sscan(value, &candidateBalance)
    }

    re = regexp.MustCompile("Locked balance: \\d+\\.\\d+")
    match = re.FindAllStringSubmatch(rawData, 1)
    for i := range match {
        value := strings.Replace(match[i][0], "Locked balance: ", "", 1)
        fmt.Sscan(value, &lockedBalance)
    }

    re = regexp.MustCompile("Active rolls: \\d+")
    match = re.FindAllStringSubmatch(rawData, 1)
    for i := range match {
        value := strings.Replace(match[i][0], "Acive rolls: ", "", 1)
        fmt.Sscan(value, &activeRolls)
    }


    re = regexp.MustCompile("Final rolls: \\d+")
    match = re.FindAllStringSubmatch(rawData, 1)
    for i := range match {
        value := strings.Replace(match[i][0], "Final rolls: ", "", 1)
        fmt.Sscan(value, &finalRolls)
    }


    re = regexp.MustCompile("Candidate rolls: \\d+")
    match = re.FindAllStringSubmatch(rawData, 1)
    for i := range match {
        value := strings.Replace(match[i][0], "Candidate rolls: ", "", 1)
        fmt.Sscan(value, &candidateRolls)
    }

    return WalletInfo{
        finalBalance: finalBalance,
        candidateBalance: candidateBalance,
        lockedBalance: lockedBalance,
        activeRolls: activeRolls,
        finalRolls: finalRolls,
        candidateRolls: candidateRolls,
    }
}
