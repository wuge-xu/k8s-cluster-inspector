package report

import (
	"fmt"
	"html/template"
	"os"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Kubernetes Cluster Inspector Report</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			margin: 40px;
			background: #f6f8fa;
			color: #24292f;
		}
		.container {
			max-width: 1100px;
			margin: auto;
			background: white;
			padding: 30px;
			border-radius: 10px;
			box-shadow: 0 2px 10px rgba(0,0,0,0.08);
		}
		h1 {
			margin-bottom: 5px;
		}
		.score {
			font-size: 36px;
			font-weight: bold;
			margin: 20px 0;
		}
		.grid {
			display: grid;
			grid-template-columns: repeat(2, 1fr);
			gap: 16px;
			margin-bottom: 30px;
		}
		.card {
			border: 1px solid #d0d7de;
			border-radius: 8px;
			padding: 16px;
			background: #ffffff;
		}
		.card h2 {
			margin-top: 0;
			font-size: 20px;
		}
		table {
			width: 100%;
			border-collapse: collapse;
			margin-top: 12px;
		}
		th, td {
			border: 1px solid #d0d7de;
			padding: 8px;
			text-align: left;
			vertical-align: top;
		}
		th {
			background: #f6f8fa;
		}
		.warning {
			color: #bc4c00;
			font-weight: bold;
		}
		.ok {
			color: #1a7f37;
			font-weight: bold;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>Kubernetes Cluster Inspector Report</h1>
		<p>Generated at: {{.Time}}</p>

		<div class="score">Health Score: {{.Score}} / 100</div>

		<div class="grid">
			<div class="card">
				<h2>Nodes</h2>
				<p>Total: {{.Nodes.Total}}</p>
				<p>Ready: {{.Nodes.Ready}}</p>
				<p>NotReady: {{.Nodes.NotReady}}</p>
			</div>

			<div class="card">
				<h2>Pods</h2>
				<p>Total: {{.Pods.Total}}</p>
				<p>Running: {{.Pods.Running}}</p>
				<p>Pending: {{.Pods.Pending}}</p>
				<p>Failed: {{.Pods.Failed}}</p>
				<p>Restart >= 3: {{.Pods.HighRestartPods}}</p>
			</div>

			<div class="card">
				<h2>Deployments</h2>
				<p>Total: {{.Deployments.Total}}</p>
				<p>Available: {{.Deployments.Available}}</p>
				<p>Unavailable: {{.Deployments.Unavailable}}</p>
			</div>

			<div class="card">
				<h2>PVCs & Events</h2>
				<p>PVC Total: {{.PVCs.Total}}</p>
				<p>PVC Pending: {{.PVCs.Pending}}</p>
				<p>PVC Lost: {{.PVCs.Lost}}</p>
				<p>Warning Events: {{.Events.Warning}}</p>
			</div>
		</div>

		<h2>Abnormal Pods</h2>
		{{if .AbnormalPods}}
		<table>
			<tr>
				<th>Namespace</th>
				<th>Name</th>
				<th>Reason</th>
				<th>Restart</th>
				<th>Suggestion</th>
			</tr>
			{{range .AbnormalPods}}
			<tr>
				<td>{{.Namespace}}</td>
				<td>{{.Name}}</td>
				<td class="warning">{{.Reason}}</td>
				<td>{{.RestartCount}}</td>
				<td>{{.Suggestion}}</td>
			</tr>
			{{end}}
		</table>
		{{else}}
		<p class="ok">No abnormal pods found.</p>
		{{end}}

		<h2>Abnormal PVCs</h2>
		{{if .AbnormalPVCs}}
		<table>
			<tr>
				<th>Namespace</th>
				<th>Name</th>
				<th>Phase</th>
				<th>Volume</th>
				<th>Suggestion</th>
			</tr>
			{{range .AbnormalPVCs}}
			<tr>
				<td>{{.Namespace}}</td>
				<td>{{.Name}}</td>
				<td class="warning">{{.Phase}}</td>
				<td>{{.VolumeName}}</td>
				<td>{{.Suggestion}}</td>
			</tr>
			{{end}}
		</table>
		{{else}}
		<p class="ok">No abnormal PVCs found.</p>
		{{end}}

		<h2>Warning Events</h2>
		{{if .AbnormalEvents}}
		<table>
			<tr>
				<th>Namespace</th>
				<th>Object</th>
				<th>Reason</th>
				<th>Count</th>
				<th>Message</th>
				<th>Suggestion</th>
			</tr>
			{{range .AbnormalEvents}}
			<tr>
				<td>{{.Namespace}}</td>
				<td>{{.InvolvedObject}}</td>
				<td class="warning">{{.Reason}}</td>
				<td>{{.Count}}</td>
				<td>{{.Message}}</td>
				<td>{{.Suggestion}}</td>
			</tr>
			{{end}}
		</table>
		{{else}}
		<p class="ok">No warning events found.</p>
		{{end}}
	</div>
</body>
</html>
`

func WriteHTML(r model.Report, filename string) {
	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		fmt.Println("html template parse failed:", err)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("create html report failed:", err)
		return
	}
	defer file.Close()

	if err := tmpl.Execute(file, r); err != nil {
		fmt.Println("write html report failed:", err)
		return
	}

	fmt.Println("Generated", filename)
}
