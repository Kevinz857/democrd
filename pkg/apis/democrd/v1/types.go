package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Mydemo 描述一个 Mydemo的资源字段
type Mydemo struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MydemoSpec `json:"spec"`
}

//MydemoSpec 为 Mydemo的资源的spec属性的字段
type MydemoSpec struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//复数形式
type MydemoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Mydemo `json:"items"`
}
