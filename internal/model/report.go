package model

type Report struct {
	Time                string               `json:"time"`
	Nodes               Nodes                `json:"nodes"`
	Pods                Pods                 `json:"pods"`
	Deployments         Deployments          `json:"deployments"`
	Services            Services             `json:"services"`
	PVCs                PVCs                 `json:"pvcs"`
	Events              Events               `json:"events"`
	AbnormalPods        []AbnormalPod        `json:"abnormal_pods"`
	AbnormalDeployments []AbnormalDeployment `json:"abnormal_deployments"`
	AbnormalPVCs        []AbnormalPVC        `json:"abnormal_pvcs"`
	AbnormalEvents      []AbnormalEvent      `json:"abnormal_events"`
	Score               int                  `json:"health_score"`
}

type Nodes struct {
	Total    int `json:"total"`
	Ready    int `json:"ready"`
	NotReady int `json:"not_ready"`
}

type Pods struct {
	Total            int `json:"total"`
	Running          int `json:"running"`
	Pending          int `json:"pending"`
	Succeeded        int `json:"succeeded"`
	Failed           int `json:"failed"`
	CrashLoopBackOff int `json:"crash_loop_backoff"`
	ImagePullBackOff int `json:"image_pull_backoff"`
	HighRestartPods  int `json:"high_restart_pods"`
}

type Deployments struct {
	Total       int `json:"total"`
	Available   int `json:"available"`
	Unavailable int `json:"unavailable"`
}

type Services struct {
	Total        int `json:"total"`
	ClusterIP    int `json:"cluster_ip"`
	NodePort     int `json:"node_port"`
	LoadBalancer int `json:"load_balancer"`
}

type PVCs struct {
	Total   int `json:"total"`
	Bound   int `json:"bound"`
	Pending int `json:"pending"`
	Lost    int `json:"lost"`
}

type Events struct {
	Total   int `json:"total"`
	Normal  int `json:"normal"`
	Warning int `json:"warning"`
}

type AbnormalPod struct {
	Namespace    string `json:"namespace"`
	Name         string `json:"name"`
	Phase        string `json:"phase"`
	Reason       string `json:"reason"`
	RestartCount int32  `json:"restart_count"`
	Suggestion   string `json:"suggestion"`
}

type AbnormalDeployment struct {
	Namespace         string `json:"namespace"`
	Name              string `json:"name"`
	DesiredReplicas   int32  `json:"desired_replicas"`
	AvailableReplicas int32  `json:"available_replicas"`
	Reason            string `json:"reason"`
	Suggestion        string `json:"suggestion"`
}

type AbnormalPVC struct {
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	Phase      string `json:"phase"`
	VolumeName string `json:"volume_name"`
	Suggestion string `json:"suggestion"`
}

type AbnormalEvent struct {
	Namespace      string `json:"namespace"`
	Name           string `json:"name"`
	InvolvedObject string `json:"involved_object"`
	Reason         string `json:"reason"`
	Message        string `json:"message"`
	Count          int32  `json:"count"`
	Suggestion     string `json:"suggestion"`
}
