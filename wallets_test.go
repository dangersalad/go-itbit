package itbit

import (
	env "github.com/dangersalad/go-environment"
	"testing"
)

const (
	envKeyItbitClientKey    = "ITBIT_CLIENT_KEY"
	envKeyItbitClientSecret = "ITBIT_CLIENT_SECRET"
	envKeyItbitUserID       = "ITBIT_USER_ID"
)

type testLogger struct {
	t *testing.T
}

func (l *testLogger) Printf(f string, a ...interface{}) {
	l.t.Logf(f, a...)
}

func (l *testLogger) Debugf(f string, a ...interface{}) {
	l.Printf(f, a...)
}

func (l *testLogger) Debug(a ...interface{}) {
	l.t.Log(a...)
}

func TestGetAllWallets(t *testing.T) {
	params, err := env.ReadOptions(env.Options{
		envKeyItbitUserID:       "",
		envKeyItbitClientSecret: "",
		envKeyItbitClientKey:    "",
	})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	SetLogger(&testLogger{t})

	wallets, err := GetAllWallets(&Config{
		ClientSecret: params[envKeyItbitClientSecret],
		ClientKey:    params[envKeyItbitClientKey],
		UserID:       params[envKeyItbitUserID],
	})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, w := range wallets {
		t.Log("ID",w.ID)
		t.Log("Name", w.Name)
		for _, b := range w.Balances {
			t.Log(b)
		}
	}
}
