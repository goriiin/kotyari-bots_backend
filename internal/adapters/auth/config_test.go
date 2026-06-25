package auth

import "testing"

func TestConfigDialAddr(t *testing.T) {
	cases := []struct {
		name string
		cfg  Config
		want string
	}{
		{
			name: "host and port",
			cfg:  Config{Host: "auth_rs", Port: 3001},
			want: "auth_rs:3001",
		},
		{
			name: "explicit addr wins over host/port",
			cfg:  Config{Host: "auth_rs", Port: 3001, Addr: "override:9999"},
			want: "override:9999",
		},
		{
			name: "addr only",
			cfg:  Config{Addr: "localhost:50051"},
			want: "localhost:50051",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.cfg.dialAddr(); got != tc.want {
				t.Fatalf("dialAddr() = %q, want %q", got, tc.want)
			}
		})
	}
}
