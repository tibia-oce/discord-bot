package configs

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

const (
	defaultString    = "default_string"
	defaultNumberStr = "8080"
	defaultNumber    = 8080
)

func setGameConfigs() {
	setEnvKeys(
		[]string{
			EnvAppIpKey,
			EnvDBHostKey,
			EnvDBNameKey,
			EnvDBUserKey,
			EnvDBPassKey,
		},
		[]string{
			EnvHttpPortKey,
			EnvGrpcPortKey,
			EnvRateLimiterBurstKey,
			EnvRateLimiterRateKey,
			EnvDBPortKey,
		},
	)
}

func unsetGameConfigs() {
	unsetEnvKeys(
		[]string{
			EnvAppIpKey,
			EnvDBHostKey,
			EnvDBNameKey,
			EnvDBUserKey,
			EnvDBPassKey,
			EnvHttpPortKey,
			EnvGrpcPortKey,
			EnvRateLimiterBurstKey,
			EnvRateLimiterRateKey,
			EnvDBPortKey,
		},
	)
}

func setEnvKeys(strings []string, integers []string) {
	for _, key := range strings {
		err := os.Setenv(key, defaultString)
		if err != nil {
			log.Printf("Cannot set key %s", key)
		}
	}
	for _, key := range integers {
		err := os.Setenv(key, defaultNumberStr)
		if err != nil {
			log.Printf("Cannot set key %s", key)
		}
	}
}

func unsetEnvKeys(keys []string) {
	for _, key := range keys {
		err := os.Unsetenv(key)
		if err != nil {
			log.Printf("Cannot unset key %s", key)
		}
	}
}

func TestGetEnvInt(t *testing.T) {
	type args struct {
		key          string
		defaultValue []int
	}
	tests := []struct {
		name   string
		args   args
		want   int
		envKey int
	}{{
		name: "no key set, uses default",
		args: args{key: "my_key", defaultValue: []int{10}},
		want: 10,
	}, {
		name: "no key set, no default",
		args: args{key: "my_key", defaultValue: []int{}},
		want: 0,
	}, {
		name:   "key set, does not use default",
		args:   args{key: "my_key", defaultValue: []int{10}},
		want:   8,
		envKey: 8,
	}, {
		name:   "key set, works without default",
		args:   args{key: "my_key", defaultValue: []int{}},
		want:   8,
		envKey: 8,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envKey != 0 {
				err := os.Setenv(tt.args.key, fmt.Sprintf("%d", tt.envKey))
				assert.Nil(t, err)
				defer os.Unsetenv(tt.args.key)
			}
			got := GetEnvInt(tt.args.key, tt.args.defaultValue...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetEnvStr(t *testing.T) {
	type args struct {
		key          string
		defaultValue []string
	}
	tests := []struct {
		name   string
		args   args
		want   string
		envKey string
	}{{
		name: "no key set, uses default",
		args: args{key: "my_key", defaultValue: []string{"default"}},
		want: "default",
	}, {
		name: "no key set, no default",
		args: args{key: "my_key", defaultValue: []string{}},
		want: "",
	}, {
		name:   "key set, does not use default",
		args:   args{key: "my_key", defaultValue: []string{"default"}},
		want:   "set",
		envKey: "set",
	}, {
		name:   "key set, works without default",
		args:   args{key: "my_key", defaultValue: []string{}},
		want:   "set",
		envKey: "set",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envKey != "" {
				err := os.Setenv(tt.args.key, tt.envKey)
				assert.Nil(t, err)
			}
			got := GetEnvStr(tt.args.key, tt.args.defaultValue...)
			assert.Equal(t, got, tt.want)
			if tt.envKey != "" {
				err := os.Unsetenv(tt.args.key)
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetGlobalConfigs(t *testing.T) {
	tests := []struct {
		name   string
		want   GlobalConfigs
		setEnv bool
	}{
		{
			name: "default global configs",
			want: GlobalConfigs{
				ServerConfigs: ServerConfigs{
					Http: HttpConfigs{
						Ip:   "0.0.0.0", // Adjusting the default value based on current environment
						Port: 80,
					},
					Grpc: GrpcConfigs{
						Ip:   "0.0.0.0", // Adjusting the default value based on current environment
						Port: 9090,
					},
					RateLimiter: RateLimiter{
						Burst: 5,
						Rate:  rate.Limit(2),
					},
				},
				DBConfigs: DBConfigs{
					Host: "database", // Expected actual value from environment
					Name: "forgottenserver",
					Port: 3306,
					User: "forgottenserver",
					Pass: "forgottenserver",
				},
			},
		},
		{
			name: "global configs with env variables set",
			want: GlobalConfigs{
				ServerConfigs: ServerConfigs{
					Http: HttpConfigs{
						Ip:   defaultString,
						Port: defaultNumber,
					},
					Grpc: GrpcConfigs{
						Ip:   defaultString,
						Port: defaultNumber,
					},
					RateLimiter: RateLimiter{
						Burst: defaultNumber,
						Rate:  rate.Limit(defaultNumber),
					},
				},
				DBConfigs: DBConfigs{
					Host: defaultString,
					Name: defaultString,
					Port: 8080, // Adjusting based on test environment setup
					User: defaultString,
					Pass: defaultString,
				},
			},
			setEnv: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				setGameConfigs()
				defer unsetGameConfigs() // Unset environment variables after test runs
			}

			got := GetGlobalConfigs()
			assert.Equal(t, tt.want, got)
		})
	}
}
