package nmap

// CopyExpiringBytesMap returns a new copy of a giving string map.
func CopyExpiringBytesMap(src map[string]ExpiringValue) map[string]ExpiringValue {
	var dest = make(map[string]ExpiringValue, len(src))
	for key, value := range src {
		if value.Expired() {
			continue
		}
		dest[key] = value
	}
	return dest
}

// CopyStringBytesMap returns a new copy of a giving string map.
func CopyStringBytesMap(src map[string][]byte) map[string][]byte {
	var dest = make(map[string][]byte, len(src))
	for key, value := range src {
		dest[key] = copyBytes(value)
	}
	return dest
}

// CopyStringMap returns a new copy of a giving string map.
func CopyStringMap(src map[string]string) map[string]string {
	var dest = make(map[string]string, len(src))
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

// CopyStringKeyMap returns a new copy of a giving string keyed map.
func CopyStringKeyMap(src map[string]interface{}) map[string]interface{} {
	var dest = make(map[string]interface{}, len(src))
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

// CopyInterfaceKeyMap returns a new copy of a giving interface keyed map.
func CopyInterfaceKeyMap(src map[interface{}]interface{}) map[interface{}]interface{} {
	var dest = make(map[interface{}]interface{}, len(src))
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

// copyBytes returns a new copy of giving byte slice.
func copyBytes(bu []byte) []byte {
	var cu = make([]byte, len(bu))
	copy(cu, bu)
	return cu
}
