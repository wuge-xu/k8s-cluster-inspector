package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/checker"
	"github.com/boserwuge/k8s-cluster-inspector/internal/client"
	"github.com/boserwuge/k8s-cluster-inspector/internal/metrics"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
	reporter "github.com/boserwuge/k8s-cluster-inspector/internal/report"
)

func main() {

	enableMetrics := flag.Bool("metrics", false, "Enable Prometheus metrics endpoint")
	metricsAddr := flag.String("metrics-addr", ":9090", "Metrics listen address")
	flag.Parse()

	clientset, err := client.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := context.Background()

	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	deployments, err := clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	services, err := clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	report := model.Report{
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}

	checker.CheckNodes(nodes.Items, &report)
	checker.CheckPods(pods.Items, &report)
	checker.CheckDeployments(deployments.Items, &report)
	checker.CheckServices(services.Items, &report)

	report.Score = checker.CalculateScore(report)

	reporter.PrintConsole(report)
	reporter.WriteJSON(report, "report.json")

	if *enableMetrics {
		metrics.Start(report, *metricsAddr)
	}
}
