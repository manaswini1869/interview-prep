package models

type Worker struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ScriptContent string `json:"script_content"`
	CPULimit      int    `json:"cpu_limit"`
}
