package membership

func senderService() error {
	memsToPing := MembershipList.members.getPingTargets(3)
	for {
		for i := 0; i < 3; i++ { // Send 3 times.
			for _, addr := range memsToPing {
				Xmtr.WriteTo([]byte{}, addr) // TODO: send marshall ping message
			}
		}
		//time.Sleep()
	}
	return nil
}

func (s *MembershipListType) getPingTargets(num int) []string {
	return nil
}
