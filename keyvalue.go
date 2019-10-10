package vdf

import (
	"errors"
	"fmt"
	"reflect"
)

type KeyValue struct {
	raw map[string]interface{}
}

func NewKeyValue(raw map[string]interface{}) (kv *KeyValue, err error) {
	if len(raw) == 0 {
		err = errors.New("raw map is empty")
		return
	}

	kv = &KeyValue{raw: raw}
	return
}

func iterateToKey(submap map[string]interface{}, keys []string) (map[string]interface{}, error) {
	if len(keys) == 0 {
		return submap, nil
	}

	key := keys[0]

	var sm interface{}
	var ok bool
	if sm, ok = submap[key]; !ok {
		return nil, errors.New(fmt.Sprintf("key '%s' does not exist", key))
	}

	switch v := sm.(type) {
	case map[string]interface{}:
		return iterateToKey(v, keys[1:])
	default:
		err := fmt.Sprintf("final value needs to be map[string]interface, failed at key '%s'", key)
		return nil, errors.New(err)
	}
}

func (k *KeyValue) GetSubMap(keys ...string) (kv *KeyValue, err error) {
	resp, err := iterateToKey(k.raw, keys)
	if err != nil {
		return
	}

	kv = &KeyValue{raw: resp}
	return
}

func (k *KeyValue) GetObject(keys ...string) (map[string]string, error) {
	resp, err := iterateToKey(k.raw, keys)

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

func (k *KeyValue) GetValue(keys ...string) (string, error) {
	finalKey := keys[len(keys)-1]
	resp, err := iterateToKey(k.raw, keys[:len(keys)-1])
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

func (k *KeyValue) GetKeys() (keys []string) {
	for key := range k.raw {
		keys = append(keys, key)
	}
	return
}
