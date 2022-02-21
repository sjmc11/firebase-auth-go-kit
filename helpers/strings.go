package helpers

import (
	"sort"
	"strconv"
	"strings"
)

func SplitStringIdsToIntSlice(stringids string) []int {
	var intids []int
	splitIds := strings.Split(stringids, ",")

	for _, v := range splitIds {
		intid, err := strconv.Atoi(strings.TrimSpace(v))
		if err == nil {
			intids = append(intids, intid)
		}
	}

	sort.Slice(intids, func(i, j int) bool {
		return intids[i] < intids[j]
	})

	return intids
}
