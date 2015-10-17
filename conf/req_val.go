package conf

// This struct we need to required values (to be sure, that user provided value)
type ReqVal struct {
	Name    string // Name
	Default string // Default value
	val     string // Defined value
	defined bool   // Flag, that value was defined
}

// Implementing stringer (flag.Value interface)
func (v *ReqVal) String() string {
	return v.Default
}

// Set value (implementing flag.Value interface)
func (v *ReqVal) Set(s string) error {
	v.val = s
	v.defined = true
	return nil
}

// To get value
func (v *ReqVal) Get() *string {
	return &v.val
}

// To make sure that value is defined
func (v *ReqVal) IsDefined() bool {
	return v.defined
}
