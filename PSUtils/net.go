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
	NET_STATS_PATH = "/sys/class/net/%s/statistics/%s_bytes"

	RX = "rx"
	TX = "tx"
)

type NetStats struct {
	RX_TOTAL int64
	TX_TOTAL int64
	RX_SPEED int64
	TX_SPEED int64
}

func (ps *PSUtils) GetNetStats() (*NetStats, error) {
	fmt.Printf("rx: %d %d, tx: %d %d \n", ps.RX_LAST_TMSTAMP, ps.RX_LAST_TOTAL, ps.TX_LAST_TMSTAMP, ps.TX_LAST_TOTAL)
	st := NetStats{}
	err := ps.GetMainInterface()
	if err != nil {
		return &st, err
	}

	ps.getTotalAndSpeed(RX, &st)
	ps.getTotalAndSpeed(TX, &st)
	fmt.Printf("netstats: %d %d %d %d\n", st.RX_TOTAL, st.TX_TOTAL, st.RX_SPEED, st.TX_SPEED)
	return &st, nil
}

func (ps *PSUtils) getTotalAndSpeed(tp string, ns *NetStats) error {
	cmdStr := fmt.Sprintf("cat "+NET_STATS_PATH, ps.NetworkInterface, tp)
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
		ns.RX_TOTAL = total
	} else {
		ns.TX_TOTAL = total
	}

	cur := time.Now().Unix()
	// fmt.Printf("tm: %s %d %d %d\n", tp, cur, ps.RX_LAST_TMSTAMP, ps.TX_LAST_TMSTAMP)
	if tp == RX && ps.RX_LAST_TMSTAMP != 0 && ps.RX_LAST_TOTAL != 0 && cur-ps.RX_LAST_TMSTAMP > 0 {
		ns.RX_SPEED = (total - ps.RX_LAST_TOTAL) / (cur - ps.RX_LAST_TMSTAMP)
	} else if tp == TX && ps.TX_LAST_TMSTAMP != 0 && ps.TX_LAST_TOTAL != 0 && cur-ps.TX_LAST_TMSTAMP > 0 {
		ns.TX_SPEED = (total - ps.TX_LAST_TOTAL) / (cur - ps.TX_LAST_TMSTAMP)
	}
	if tp == RX {
		ps.RX_LAST_TOTAL = total
		ps.RX_LAST_TMSTAMP = cur
	} else {
		ps.TX_LAST_TOTAL = total
		ps.TX_LAST_TMSTAMP = cur
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

	lines := SplitString(s)
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
