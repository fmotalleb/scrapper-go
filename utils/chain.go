package utils

type ChainCallback = func(any) (any, error)

func ChainExec(item any, chain []ChainCallback) (any, error) {
	result := item
	for _, r := range chain {
		var err error
		if result, err = r(result); err != nil {
			return nil, err
		}
	}
	return result, nil
}
