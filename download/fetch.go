package download

import "time"

type Fetch func() (recvd bool, err error)

func DoFetch(f Fetch, ms int) (recvd bool, err error) {
	deadline := time.NewTicker(time.Duration(ms) * time.Millisecond)
	defer deadline.Stop()
	for {
		select {
		case <-deadline.C:
			return recvd, nil

		default:
			recvd, err = f()
			if err != nil {
				return false, err
			}
			return recvd, nil
		}
	}
}
