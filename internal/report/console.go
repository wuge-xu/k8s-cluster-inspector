package report

import (
	"fmt"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func PrintConsole(r model.Report) {
	fmt.Println("================================")
	fmt.Println(" Kubernetes Cluster Inspector")
	fmt.Println("================================")
	fmt.Println("Time:", r.Time)
	fmt.Println()

	fmt.Printf("Health Score: %d / 100\n\n", r.Score)

	fmt.Printf("Nodes: total=%d ready=%d notReady=%d\n",
		r.Nodes.Total, r.Nodes.Ready, r.Nodes.NotReady)

	fmt.Printf("Pods: total=%d running=%d pending=%d failed=%d restart>=3=%d\n",
		r.Pods.Total, r.Pods.Running, r.Pods.Pending, r.Pods.Failed, r.Pods.HighRestartPods)

	fmt.Printf("Deployments: total=%d available=%d unavailable=%d\n",
		r.Deployments.Total, r.Deployments.Available, r.Deployments.Unavailable)

	fmt.Printf("Services: total=%d clusterIP=%d nodePort=%d loadBalancer=%d\n",
		r.Services.Total, r.Services.ClusterIP, r.Services.NodePort, r.Services.LoadBalancer)

	fmt.Printf("PVCs: total=%d bound=%d pending=%d lost=%d\n",
		r.PVCs.Total, r.PVCs.Bound, r.PVCs.Pending, r.PVCs.Lost)

	fmt.Printf("Events: total=%d normal=%d warning=%d\n",
		r.Events.Total, r.Events.Normal, r.Events.Warning)

	fmt.Printf("Namespaces: total=%d\n\n", r.Namespaces.Total)

	fmt.Println("Namespace Summary")
	fmt.Println("-----------------")
	if len(r.Namespaces.Items) == 0 {
		fmt.Println("No namespaces found.")
	} else {
		for _, ns := range r.Namespaces.Items {
			fmt.Printf("%s pods=%d running=%d pending=%d failed=%d deployments=%d unavailableDeployments=%d pvcs=%d pendingPVCs=%d\n",
				ns.Name,
				ns.PodsTotal,
				ns.PodsRunning,
				ns.PodsPending,
				ns.PodsFailed,
				ns.DeploymentsTotal,
				ns.DeploymentsUnavailable,
				ns.PVCsTotal,
				ns.PVCsPending,
			)
		}
	}
	fmt.Println()

	fmt.Println("Warning Event Reasons")
	fmt.Println("---------------------")
	if len(r.EventReasons) == 0 {
		fmt.Println("No warning event reasons found.")
	} else {
		for _, er := range r.EventReasons {
			fmt.Printf("%s count=%d\n", er.Reason, er.Count)
		}
	}
	fmt.Println()

	fmt.Println("Abnormal Pods")
	fmt.Println("-------------")
	if len(r.AbnormalPods) == 0 {
		fmt.Println("No abnormal pods found.")
	} else {
		for _, p := range r.AbnormalPods {
			fmt.Printf("%s/%s reason=%s restart=%d\n", p.Namespace, p.Name, p.Reason, p.RestartCount)
			fmt.Printf("  suggestion: %s\n", p.Suggestion)
		}
	}
	fmt.Println()

	fmt.Println("Abnormal PVCs")
	fmt.Println("-------------")
	if len(r.AbnormalPVCs) == 0 {
		fmt.Println("No abnormal PVCs found.")
	} else {
		for _, pvc := range r.AbnormalPVCs {
			fmt.Printf("%s/%s phase=%s volume=%s\n", pvc.Namespace, pvc.Name, pvc.Phase, pvc.VolumeName)
			fmt.Printf("  suggestion: %s\n", pvc.Suggestion)
		}
	}
	fmt.Println()

	fmt.Println("Recent Warning Events")
	fmt.Println("---------------------")
	if len(r.AbnormalEvents) == 0 {
		fmt.Println("No warning events found.")
	} else {
		for _, e := range r.AbnormalEvents {
			fmt.Printf("%s/%s object=%s reason=%s count=%d\n",
				e.Namespace, e.Name, e.InvolvedObject, e.Reason, e.Count)
			fmt.Printf("  message: %s\n", e.Message)
			fmt.Printf("  suggestion: %s\n", e.Suggestion)
		}
	}
}
