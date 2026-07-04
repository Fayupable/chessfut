package chesscom

type profileResponse struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Avatar    string `json:"avatar"`
	Followers int    `json:"followers"`
	Country   string `json:"country"`
	Joined    int64  `json:"joined"`
}

type statsRecord struct {
	Last struct {
		Rating int `json:"rating"`
	} `json:"last"`
	Best struct {
		Rating int `json:"rating"`
	} `json:"best"`
	Record struct {
		Win  int `json:"win"`
		Loss int `json:"loss"`
		Draw int `json:"draw"`
	} `json:"record"`
}

type statsResponse struct {
	ChessDaily  statsRecord `json:"chess_daily"`
	ChessRapid  statsRecord `json:"chess_rapid"`
	ChessBullet statsRecord `json:"chess_bullet"`
	ChessBlitz  statsRecord `json:"chess_blitz"`
	Fide        int         `json:"fide"`
}

type archivesResponse struct {
	Archives []string `json:"archives"`
}

type gamePlayer struct {
	Rating   int    `json:"rating"`
	Result   string `json:"result"`
	Username string `json:"username"`
}

type gameResponse struct {
	PGN         string     `json:"pgn"`
	TimeControl string     `json:"time_control"`
	EndTime     int64      `json:"end_time"`
	Rated       bool       `json:"rated"`
	TimeClass   string     `json:"time_class"`
	White       gamePlayer `json:"white"`
	Black       gamePlayer `json:"black"`
	ECO         string     `json:"eco"`
}

type gamesResponse struct {
	Games []gameResponse `json:"games"`
}

type titledResponse struct {
	Players []string `json:"players"`
}
