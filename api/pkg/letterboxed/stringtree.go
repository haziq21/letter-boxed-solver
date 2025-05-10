package letterboxed

type StringTree map[string]*StringTree

func (wt *StringTree) PushSequence(seq []string) {
	if len(seq) == 0 {
		return
	}

	if _, ok := (*wt)[seq[0]]; !ok {
		(*wt)[seq[0]] = &StringTree{}
	}

	(*wt)[seq[0]].PushSequence(seq[1:])
}

func (wt *StringTree) PopLeaf() []string {
	for child, grandChildren := range *wt {
		subSequence := grandChildren.PopLeaf()
		if len(*grandChildren) == 0 {
			delete(*wt, child)
		}

		return append([]string{child}, subSequence...)
	}

	return nil
}
