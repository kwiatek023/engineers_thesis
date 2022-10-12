package simulation

import (
	"encoding/binary"
	"hash/fnv"
	"math"
	"math/bits"
	"math/rand"
)

type HllProtocol struct{}

type HyperLogLog struct {
	registers []float64
	m         uint // number of registers
	b         uint
}

func NewHyperLogLog(m uint) HyperLogLog {
	return HyperLogLog{
		registers: make([]float64, m),
		m:         m,
		b:         uint(math.Ceil(math.Log2(float64(m)))),
	}
}

func (HllProtocol) GetInitialData(station IStation) {
	min := 1
	max := 1000000
	observedValuesAsBytes := make([][]byte, 0)
	for i := 0; i < 20; i++ {
		randomValue := min + rand.Intn(max+1)
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(randomValue))
		station.ObserveValue([]float64{float64(randomValue)})
		observedValuesAsBytes = append(observedValuesAsBytes, b)
	}

	h := NewHyperLogLog(32)
	for _, b := range observedValuesAsBytes {
		h.Add(b)
	}

	station.SetCurrentData(h.registers)
}

func (HllProtocol) OnInitialize(station IStation) {
	station.Broadcast()
}

func (HllProtocol) OnDataReceive(station IStation) {
	station.SetUserDefinedVariable("vectorChanged", false)
	mq := station.GetMsgQueue()
	currentVector := station.GetCurrentData()
	for mq.Len() > 0 {
		msg := mq.Dequeue()
		vector := msg.Data
		for i, element := range vector {
			if currentVector[i] < element {
				station.SetUserDefinedVariable("vectorChanged", true)
				currentVector[i] = element
			}
		}
	}
}

func (HllProtocol) OnDataPropagate(station IStation) {
	vectorChanged := station.GetUserDefinedVariable("vectorChanged").(bool)
	if vectorChanged {
		// state changed, inform neighbours
		station.SynchronizedBroadcast()
	}

}

func (HllProtocol) StopCondition(station IStation) bool {
	return station.GetRoundCounter() < station.GetGraph().GetDiameter()
}

func (HllProtocol) OnFinalize(station IStation) {
	currentVector := station.GetCurrentData()
	sum := 0.
	m := 32.0
	for _, v := range currentVector {
		if v != 0 {
			sum += math.Pow(math.Pow(2, v), -1)
		}
	}
	estimate := int(.79402 * m * m / sum)

	station.SetResult(float64(estimate))
	//fmt.Println("result of", station.GetId(),
	//	station.GetCurrentData(), "msg: ", station.GetSentMsgCounter(), "estimate:", estimate)
}

func (HllProtocol) CalculateStationExactResult(station IStation) float64 {
	return -1
}

func (HllProtocol) CalculateGlobalExactResult(stations *[]IStation) float64 {
	m := map[float64]struct{}{}

	for _, station := range *stations {
		for _, i := range station.GetObservedValues() {
			if _, ok := m[i[0]]; !ok {
				m[i[0]] = struct{}{}
			}
		}
	}

	return float64(len(m))
}

func leftmostSignificantBitPosition(x uint32) int {
	return 1 + bits.LeadingZeros32(x)
}

func create32BitHash(stream []byte) uint32 {
	h := fnv.New32()
	h.Write(stream)
	sum := h.Sum32()
	h.Reset()
	return sum
}

func (h HyperLogLog) Add(data []byte) HyperLogLog {
	x := create32BitHash(data)
	k := 32 - h.b // first b bits
	r := float64(leftmostSignificantBitPosition(x << h.b))
	j := x >> uint(k)

	if r > h.registers[j] {
		h.registers[j] = r
	}

	return h
}

func (h HyperLogLog) Count() uint64 {
	sum := 0.
	m := float64(h.m)
	numOfRegistersEqualToZero := 0
	for _, v := range h.registers {
		if v != 0 {
			sum += math.Pow(math.Pow(2, float64(v)), -1)
		} else {
			numOfRegistersEqualToZero++
		}
	}
	estimate := .697 * m * m / sum
	return h.rangeCorrection(estimate, numOfRegistersEqualToZero)
}

func (h HyperLogLog) rangeCorrection(estimate float64, numOfRegistersEqualToZero int) uint64 {
	var result uint64
	m := 32
	if estimate <= (5/2)*float64(m) {
		if numOfRegistersEqualToZero != 0 {
			result = uint64(float64(m) * math.Log2(float64(m/numOfRegistersEqualToZero)))
		} else {
			result = uint64(estimate)
		}
	} else if estimate <= (1/30)*math.Pow(2, 32) {
		result = uint64(estimate)
	} else {
		x := math.Pow(2, 32)
		result = uint64(-1 * x * math.Log2(1-estimate/x))
	}
	return result
}
