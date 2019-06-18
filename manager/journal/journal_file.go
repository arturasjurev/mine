package journal

// this is journal to file implementation

type JournalFileService struct {
	File         string
	DumpOnChange bool
}

func (j *JournalFileService) Init() error {
	return nil
}

func (j *JournalFileService) ListClients() ([]Client, error) {
	return []Client{}, nil
}

func (j *JournalFileService) Client(id string) (Client, error) {
	return Client{}, nil
}

func (j *JournalFileService) UpsertClient(c Client) (Client, error) {
	return Client{}, nil
}

func (j *JournalFileService) ListOrders() ([]Order, error) {
	return []Order{}, nil
}

func (j *JournalFileService) Order(id string) (Order, error) {
	return Order{}, nil
}

func (j *JournalFileService) UpsertOrder(o Order) (Order, error) {
	return Order{}, nil
}
