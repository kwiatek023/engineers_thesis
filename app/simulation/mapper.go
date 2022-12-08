package simulation

func (m Manager) mapNameToProtocol(name string) Protocol {
	if name == "hll" {
		return HllProtocol{}
	} else if name == "minPropagation" {
		return MinPropagationProtocol{}
	}
	return nil
}
