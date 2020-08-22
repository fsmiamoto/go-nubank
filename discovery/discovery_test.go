package discovery

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func Test_create(t *testing.T) {
	type args struct {
		client HTTPClientGet
	}
	tests := []struct {
		name    string
		args    args
		want    *Discovery
		wantErr bool
	}{
		{
			name: "valid client",
			args: args{
				client: &MockClient{
					responses: map[string]*http.Response{
						DISCOVERY_URL: &http.Response{
							Status:     "200 OK",
							StatusCode: 200,
							Body: buildBodyFromString(`
							{
								"login":          "url_login",
								"reset_password": "url_reset_password",
								"email_verify":   "url_email_verify_token"
							}
						`),
						},
						DISCOVERY_APP_URL: &http.Response{
							Status:     "200 OK",
							StatusCode: 200,
							Body: buildBodyFromString(`
							{
								"scopes":       "url_scopes",
								"userinfo":     "url_userinfo",
								"revoke_token": "url_revoke_token",
								"faq": {
									"ios": "url_ios",
									"android": "url_android",
									"wp": "url_windows_phone"
								}
							}
						`),
						},
					},
				},
			},
			want: &Discovery{
				services: map[string]string{
					"login":          "url_login",
					"reset_password": "url_reset_password",
					"email_verify":   "url_email_verify_token",
					"scopes":         "url_scopes",
					"userinfo":       "url_userinfo",
					"revoke_token":   "url_revoke_token",
				},
			},
			wantErr: false,
		},
		{
			name: "error on api",
			args: args{
				client: &MockClient{
					responses: map[string]*http.Response{
						DISCOVERY_URL: &http.Response{
							StatusCode: 500,
							Body:       buildBodyFromString(""),
						},
						DISCOVERY_APP_URL: &http.Response{
							StatusCode: 500,
							Body:       buildBodyFromString(""),
						},
					},
				},
			},
			want:    &Discovery{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := create(tt.args.client)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if !reflect.DeepEqual(got.services, tt.want.services) {
				t.Errorf("New() = %v, want %v", got.services, tt.want.services)
			}
		})
	}
}

func TestServiceURL(t *testing.T) {
	type fields struct {
		services map[string]string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "existing service",
			fields: fields{
				services: map[string]string{
					"login":          "url_login",
					"reset_password": "url_reset_password",
					"email_verify":   "url_email_verify_token",
				},
			},
			args:    args{name: "login"},
			want:    "url_login",
			wantErr: false,
		},
		{
			name: "missing service",
			fields: fields{
				services: map[string]string{
					"login":          "url_login",
					"reset_password": "url_reset_password",
					"email_verify":   "url_email_verify",
				},
			},
			args:    args{name: "foo"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discovery{
				services: tt.fields.services,
			}
			got, err := d.ServiceURL(tt.args.name)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if got != tt.want {
				t.Errorf("Discovery.serviceURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

type MockClient struct {
	responses map[string]*http.Response
}

func (m *MockClient) Get(url string) (*http.Response, error) {
	return m.responses[url], nil
}

func buildBodyFromString(s string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(s)))
}
