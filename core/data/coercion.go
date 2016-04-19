package data

import (
	"fmt"
	"strconv"
	"encoding/json"
)

func CoerceToTypedValue(value interface{}, dataType Type) (*TypedValue, error) {

	coerced, err := CoerceToValue(value, dataType)

	if err != nil {
		return nil, err
	}

	return &TypedValue{Type:dataType, Value:coerced}, nil
}

func CoerceToValue(value interface{}, dataType Type) (interface{}, error) {

	var coerced interface{}
	var err error

	switch dataType {
	case STRING:
		coerced, err = CoerceToString(value)
	case INTEGER:
		coerced, err = CoerceToInteger(value)
	case NUMBER:
		coerced, err = CoerceToNumber(value)
	case BOOLEAN:
		coerced, err = CoerceToBoolean(value)
	case OBJECT:
		coerced, err = CoerceToObject(value)
	case ARRAY:
		coerced, err = CoerceToArray(value)
	case MAP:
		coerced, err = CoerceToMap(value)
	}

	if err != nil {
		return nil, err
	}

	return coerced, nil
}

//todo check int64,float64 on raspberry pi

func CoerceToString(val interface{}) (string, error) {
	switch t := val.(type) {
	case string:
		return t, nil
	case int:
		return strconv.Itoa(t), nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case json.Number:
		return t.String(),nil
	case bool:
		return strconv.FormatBool(t), nil
	case nil:
		return "", nil
	case map[string]interface{}:
		b,err := json.Marshal(t)
		if err != nil {
			return "", err
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("Unable to Coerce %#v to string", t)
	}
}

func CoerceToInteger(val interface{}) (int, error) {
	switch t := val.(type) {
	case int:
		return t, nil
	case int64:
		return int(t), nil
	case float64:
		return int(t), nil
	case json.Number:
		 i, err := t.Int64()
		 return int(i),err
	case string:
		return strconv.Atoi(t)
	case bool:
		if t {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("Unable to coerce %#v to integer", val)
	}
}

func CoerceToNumber(val interface{}) (float64, error) {
	switch t := val.(type) {
	case int:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case float64:
		return t, nil
	case json.Number:
		return t.Float64()
	case string:
		return  strconv.ParseFloat(t, 64)
	case bool:
		if t {
			return 1.0, nil
		}
		return 0.0, nil
	case nil:
		return 0.0, nil
	default:
		return 0.0, fmt.Errorf("Unable to coerce %#v to float", val)
	}
}

func CoerceToBoolean(val interface{}) (bool, error) {
	switch t := val.(type) {
	case bool:
		return t, nil
	case int:
		return t != 0, nil
	case int64:
		return t != 0, nil
	case float64:
		return t != 0.0, nil
	case json.Number:
		i, err := t.Int64()
		return i != 0, err
	case string:
		return strconv.ParseBool(t)
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("Unable to coerce %#v to bool", val)
	}
}

func CoerceToObject(val interface{}) (map[string]interface{}, error) {

	switch t := val.(type) {
	case map[string]interface{}:
		return t, nil
	default:
		return nil, fmt.Errorf("Unable to coerce %#v to map[string]interface{}", val)
	}
}

func CoerceToArray(val interface{}) ([]interface{}, error) {

	switch t := val.(type) {
	case []interface{}:
		return t, nil
	case []map[string]interface{}:
		var a []interface{}
		for _, v := range t {
			a = append(a, v)
		}
		return a, nil
	default:
		return nil, fmt.Errorf("Unable to coerce %#v to []interface{}", val)
	}
}

func CoerceToMap(val interface{}) (map[string]string, error) {

	switch t := val.(type) {
	case map[string]string:
		return t, nil
	case map[string]interface{}:

		var m = make(map[string]string, len(t))
		for k, v := range t {

			mVal, err := CoerceToString(v)
			if err != nil {
				return nil, err
			}
			m[k] = mVal
		}
		return m, nil
	case map[interface{}]string:

		var m = make(map[string]string, len(t))
		for k, v := range t {

			mKey , err := CoerceToString(k)
			if err != nil {
				return nil, err
			}
			m[mKey] = v
		}
		return m, nil
	case map[interface{}]interface{}:

		var m = make(map[string]string, len(t))
		for k, v := range t {

			mKey , err := CoerceToString(k)
			if err != nil {
				return nil, err
			}

			mVal, err := CoerceToString(v)
			if err != nil {
				return nil, err
			}
			m[mKey] = mVal
		}
		return m, nil
	default:
		return nil, fmt.Errorf("Unable to coerce %#v to map[string]string", val)
	}
}
