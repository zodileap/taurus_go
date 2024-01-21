package field

// 表字段
type (
	// Scanner is an interface used by Scan.
	Scanner interface {
		// Scan assigns a value from a database driver.
		//
		// The src value will be of one of the following types:
		//
		//    int64
		//    float64
		//    bool
		//    []bytetaurus_go/entity/field/field.go
		//    string
		//    time.Time
		//    nil - for NULL values
		//
		// An error should be returned if the value cannot be stored
		// without loss of information.
		//
		// Reference types such as []byte are only valid until the next call to Scan
		// and should not be retained. Their underlying memory is owned by the driver.
		// If retention is necessary, copy their values before the next call to Scan.
		Scan(src any) error
	}
)
