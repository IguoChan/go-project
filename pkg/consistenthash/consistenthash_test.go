package consistenthash

import (
	"strconv"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// Adds 8, 18, 28
	hash.Add("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}

var (
	uuids = []string{
		"ae5429e0-fe9b-4d27-abac-3126d4436dac",
		"7ff1ea5d-0de7-4cfb-8f69-84e53e501f2e",
		"15534837-dffe-4756-a776-0750a5ddf11f",
		"baa84487-33fe-4ba1-a4f0-ec028ea3dd83",
		"26d898ef-daa6-4b6e-a0ee-02d31b42c61f",
		"0a93250c-f6aa-40d9-950e-045ca958fff4",
		"4642cc8c-1deb-4cdf-88df-dfcff599aef4",
		"d10ecfa5-ec64-4001-ba44-bc5dfed629f0",
		"2f1adac1-9bfb-4cf7-a078-7372a07619ea",
		"ee76bdd9-34f0-477d-84c7-c5868c5a2ef6",
	}
)

func TestUuids(t *testing.T) {
	var uuids []string
	for i := 0; i < 10; i++ {
		uuids = append(uuids, uuid.NewV4().String())
	}
	t.Log(uuids)
}

func TestHash2(t *testing.T) {
	hash := New(3, nil)
	hash.Add("10.171.2.23:2345", "10.171.6.65:2345", "10.171.2.7.68:2345")
	hash.Add("1.1.1.1:2345")
	var vs []string
	for i := range uuids {
		vs = append(vs, hash.Get(uuids[i]))
	}
	t.Log(vs)

	var vs1 []string
	for i := range uuids {
		vs1 = append(vs1, hash.Get(uuids[i]))
	}
	t.Log(vs1)

	for i := range vs {
		if vs[i] != vs1[i] {
			t.Fatal("not equal")
		}
	}

	hash.Add()

}
