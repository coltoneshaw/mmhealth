package processpacket

type CheckType string

const (
	Proactive CheckType = "proactive"
	Health    CheckType = "health"
	Adoption  CheckType = "adoption"
)

type CheckResult struct {
	Name        string
	Result      string
	Type        CheckType
	Description string
	Status      string
}

func ProcessPacket(packet PacketData) error {

	configChecks(packet.Config)
	return nil
}
