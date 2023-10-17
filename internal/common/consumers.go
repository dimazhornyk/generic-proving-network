package common

import "slices"

type Consumer struct {
	Name  string
	Image string
}

var consumers = []Consumer{
	{
		Name:  "zksync",
		Image: "zksyncio/zksync-prover:latest",
	},
	{
		Name:  "scroll",
		Image: "scroll-network/scroll-prover:v1.0.0",
	},
}

func GetConsumers(list []string) []Consumer {
	res := make([]Consumer, 0)
	for _, c := range consumers {
		if slices.Contains(list, c.Name) {
			res = append(res, c)
		}
	}

	return res
}
