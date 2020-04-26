package common

func NonRepeatElements(elements []string) []string {
	var es []string
	var elementExist = make(map[string]bool, len(elements))
	for _, e := range elements {
		_, ok := elementExist[e]
		if ok {
			continue
		}

		es = append(es, e)
		elementExist[e] = true
	}

	return es
}
