package conn

func Conn() *nats.Conn{
	nc, err := nats.Connect("nats://localhost:4221")
	if err != nil {
		log.Fatal(err)
	}

	return nc
}
