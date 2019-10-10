package vdf

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// KeyValue wrapper around the raw parsed vdf file. When creating a new KeyValue use the provided function
type KeyValue struct {
	Raw             map[string]interface{}
	caseInsensitive bool
}

// NewKeyValue creates a new key value struct, case flag default should be false
func NewKeyValue(raw map[string]interface{}, caseInsensitive bool) (kv *KeyValue, err error) {
	if len(raw) == 0 {
		err = errors.New("raw map is empty")
		return
	}

	kv = &KeyValue{Raw: raw, caseInsensitive: caseInsensitive}
	return
}

func iterateToKey(submap map[string]interface{}, keys []string, ci bool) (map[string]interface{}, error) {
	if len(keys) == 0 {
		return submap, nil
	}

	key := keys[0]
	if ci {
		key = strings.ToLower(key)
	}

	var sm interface{}
	var ok bool
	if sm, ok = submap[key]; !ok {
		return nil, errors.New(fmt.Sprintf("key '%s' does not exist", key))
	}

	switch v := sm.(type) {
	case map[string]interface{}:
		return iterateToKey(v, keys[1:], ci)
	default:
		err := fmt.Sprintf("final value needs to be map[string]interface, failed at key '%s'", key)
		return nil, errors.New(err)
	}
}

// GetSubMap queries for a sub map based on the ordered keys provided
func (k *KeyValue) GetSubMap(keys ...string) (kv *KeyValue, err error) {
	resp, err := iterateToKey(k.Raw, keys, k.caseInsensitive)
	if err != nil {
		return
	}

	kv = &KeyValue{Raw: resp}
	return
}

// GetObject queries for an object (map[string]string) based on the ordered keys provided
func (k *KeyValue) GetObject(keys ...string) (map[string]string, error) {
	resp, err := iterateToKey(k.Raw, keys, k.caseInsensitive)

	if err != nil {
		return nil, err
	}

	obj := make(map[string]string)
	for key, val := range resp {
		switch v := val.(type) {
		case string:
			obj[key] = v
		default:
			err := fmt.Sprintf("key '%s' with value of type '%s' is invalid, needs to be string",
				key, reflect.TypeOf(v))
			return nil, errors.New(err)
		}
	}

	return obj, nil
}

// GetValue queries for the value based on the ordered keys provided
func (k *KeyValue) GetValue(keys ...string) (string, error) {
	finalKey := keys[len(keys)-1]
	if k.caseInsensitive {
		finalKey = strings.ToLower(finalKey)
	}

	resp, err := iterateToKey(k.Raw, keys[:len(keys)-1], k.caseInsensitive)
	if err != nil {
		return "", err
	}

	var val interface{}
	var ok bool
	if val, ok = resp[finalKey]; !ok {
		return "", errors.New(fmt.Sprintf("key '%s' does not exist", finalKey))
	}

	switch v := val.(type) {
	case string:
		return v, nil
	default:
		return "", errors.New(fmt.Sprintf("final value not string, got '%s'", reflect.TypeOf(v)))
	}
}

// GetKeys returns the keys at the root map
func (k *KeyValue) GetKeys() (keys []string) {
	for key := range k.Raw {
		keys = append(keys, key)
	}
	return
}
