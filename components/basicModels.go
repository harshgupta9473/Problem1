package components

type PostRequest struct {
	HospitalID    string `json:"hospitalID"`
	PatientName string `json:"name"`
	Age         int    `json:"age"`
	Problem     string `json:"problem"`
}

