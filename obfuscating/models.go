package obfuscating

type ConnectionInfo struct {
	User     string `binding:"required"`
	Password string `binding:"required"`
	Host     string `binding:"required"`
	Schema   string `binding:"required"`
}

type ObfuscateRequest struct {
	Model       map[string][]Column `binding:"required"`
	Origin      ConnectionInfo      `binding:"required"`
	Destination ConnectionInfo      `binding:"required"`
}

type Column struct {
	Name string `binding:"required"`
	Type string `binding:"required"`
	//can obfuscate in obfuscating context, need to obfuscate in request context
	NeedToObfuscate bool `binding:"required"`
	IsPrimaryKey    bool `binding:"required"`
}
