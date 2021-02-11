package PSUtils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// path and type(rx, tx)
	NetStatsPath = "/sys/class/net/%s/statistics/%s_bytes"

	RX = "rx"
	TX = "tx"
)

type NetStats struct {
	RxTotal int64
	TxTotal int64
	RxSpeed int64
	TxSpeed int64
}

func (ps *PSUtils) GetNetStats() (*NetStats, error) {
	fmt.Printf("rx: %d %d, tx: %d %d \n", ps.RxLastTmstamp, ps.RxLastTotal, ps.TxLastTmstamp, ps.TxLastTotal)
	st := NetStats{}
	err := ps.GetMainInterface()
	if err != nil {
		return &st, err
	}

	ps.getTotalAndSpeed(RX, &st)
	ps.getTotalAndSpeed(TX, &st)
	fmt.Printf("netstats: %d %d %d %d\n", st.RxTotal, st.TxTotal, st.RxSpeed, st.TxSpeed)
	return &st, nil
}

func (ps *PSUtils) getTotalAndSpeed(tp string, ns *NetStats) error {
	cmdStr := fmt.Sprintf("cat "+NetStatsPath, ps.NetworkInterface, tp)
	res, err := ps.Exec(cmdStr)
	if err != nil {
		return err
	}
	fmt.Println(res)
	total, err := strconv.ParseInt(res, 10, 64)
	if err != nil {
		return err
	}
	if tp == RX {
		ns.RxTotal = total
	} else {
		ns.TxTotal = total
	}

	cur := time.Now().Unix()
	// fmt.Printf("tm: %s %d %d %d\n", tp, cur, ps.RX_LAST_TMSTAMP, ps.TX_LAST_TMSTAMP)
	if tp == RX && ps.RxLastTmstamp != 0 && ps.RxLastTotal != 0 && cur-ps.RxLastTmstamp > 0 {
		ns.RxSpeed = (total - ps.RxLastTotal) / (cur - ps.RxLastTmstamp)
	} else if tp == TX && ps.TxLastTmstamp != 0 && ps.TxLastTotal != 0 && cur-ps.TxLastTmstamp > 0 {
		ns.TxSpeed = (total - ps.TxLastTotal) / (cur - ps.TxLastTmstamp)
	}
	if tp == RX {
		ps.RxLastTotal = total
		ps.RxLastTmstamp = cur
	} else {
		ps.TxLastTotal = total
		ps.TxLastTmstamp = cur
	}

	return nil
}

func (ps *PSUtils) GetMainInterface() error {
	if ps.NetworkInterface != "" {
		return nil
	}

	s, err := ps.Exec("cat /proc/net/route")
	if err != nil {
		return err
	}

	lines := SplitStringToLines(s)
	for _, line := range lines[1:] {
		sp := strings.Split(line, "\t")
		if len(sp) < 11 {
			continue
		}
		if sp[2] == "00000000" && sp[7] != "00000000" {
			ps.NetworkInterface = sp[0]
			return nil
		}
	}

	return errors.New("network interface not find")
}
