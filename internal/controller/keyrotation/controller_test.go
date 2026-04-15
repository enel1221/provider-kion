/*
Copyright 2024 Upbound Inc.
*/

package keyrotation

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestHasRotateAnnotation(t *testing.T) {
	tests := map[string]struct {
		annotations map[string]string
		want        bool
	}{
		"NoAnnotations": {
			annotations: nil,
			want:        false,
		},
		"AnnotationMissing": {
			annotations: map[string]string{"other": "value"},
			want:        false,
		},
		"AnnotationFalse": {
			annotations: map[string]string{AnnotationRotate: "false"},
			want:        false,
		},
		"AnnotationTrue": {
			annotations: map[string]string{AnnotationRotate: "true"},
			want:        true,
		},
		"AnnotationTrueCaseInsensitive": {
			annotations: map[string]string{AnnotationRotate: "True"},
			want:        true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: tc.annotations,
				},
			}
			got := hasRotateAnnotation(secret)
			if got != tc.want {
				t.Errorf("hasRotateAnnotation() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := map[string]struct {
		input interface{}
		want  bool
	}{
		"StringTrue":      {input: "true", want: true},
		"StringTrueUpper": {input: "TRUE", want: true},
		"StringFalse":     {input: "false", want: false},
		"BoolTrue":        {input: true, want: true},
		"BoolFalse":       {input: false, want: false},
		"Nil":             {input: nil, want: false},
		"Number":          {input: float64(1), want: false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := parseBool(tc.input)
			if got != tc.want {
				t.Errorf("parseBool(%v) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestCallRotate(t *testing.T) {
	tests := map[string]struct {
		serverHandler http.HandlerFunc
		wantKey       string
		wantErr       bool
	}{
		"SuccessfulRotation": {
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				// Verify request.
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.Header.Get("Authorization") != "Bearer old-key" {
					t.Errorf("unexpected auth header: %s", r.Header.Get("Authorization"))
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("unexpected content-type: %s", r.Header.Get("Content-Type"))
				}

				var body map[string]string
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("decode body: %v", err)
				}
				if body["key"] != "old-key" {
					t.Errorf("expected key=old-key, got %s", body["key"])
				}

				w.WriteHeader(http.StatusCreated)
				resp := rotateResponse{}
				resp.Status = 201
				resp.Data.ID = 42
				resp.Data.Key = "new-key-abc"
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					t.Fatalf("encode response: %v", err)
				}
			},
			wantKey: "new-key-abc",
		},
		"UnauthorizedResponse": {
			serverHandler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"status":401,"error":"unauthorized"}`)) //nolint:errcheck
			},
			wantErr: true,
		},
		"EmptyKeyInResponse": {
			serverHandler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"status":201,"data":{"id":1,"key":""}}`)) //nolint:errcheck
			},
			wantErr: true,
		},
		"MalformedJSON": {
			serverHandler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`not json`)) //nolint:errcheck
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(tc.serverHandler)
			defer server.Close()

			r := &Reconciler{
				http: &http.Client{Timeout: 5 * time.Second},
			}

			key, err := r.callRotate(t.Context(), server.URL, "old-key", false)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.wantKey, key); diff != "" {
				t.Errorf("key mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// newFakeSecret creates a Secret with the given credentials and annotations.
func newFakeSecret(ns, name string, creds map[string]string, annotations map[string]string) *corev1.Secret {
	raw, err := json.Marshal(creds)
	if err != nil {
		panic(err)
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   ns,
			Annotations: annotations,
		},
		Data: map[string][]byte{
			"credentials": raw,
		},
	}
}

func TestReconcile(t *testing.T) {
	tests := map[string]struct {
		secret       *corev1.Secret
		serverStatus int
		serverKey    string
		wantRequeue  bool
		wantNewKey   string // if non-empty, verify the key was updated in the Secret
		wantStatus   string // expected rotation-status annotation value
	}{
		"SuccessfulRotation": {
			secret: newFakeSecret("default", "creds",
				map[string]string{"apikey": "old-key", "url": "PLACEHOLDER"},
				map[string]string{AnnotationRotate: "true"},
			),
			serverStatus: http.StatusCreated,
			serverKey:    "fresh-key-123",
			wantNewKey:   "fresh-key-123",
			wantStatus:   "success",
		},
		"RotationNotYetDue": {
			secret: newFakeSecret("default", "creds",
				map[string]string{"apikey": "old-key", "url": "PLACEHOLDER"},
				map[string]string{
					AnnotationRotate:       "true",
					AnnotationLastRotation: time.Now().UTC().Format(time.RFC3339),
				},
			),
			wantRequeue: true,
			wantNewKey:  "", // key should NOT change
		},
		"MissingAPIKey": {
			secret: newFakeSecret("default", "creds",
				map[string]string{"url": "PLACEHOLDER"},
				map[string]string{AnnotationRotate: "true"},
			),
			wantStatus: "error: credentials JSON must contain non-empty apikey and url",
		},
		"AnnotationDisabled": {
			secret: newFakeSecret("default", "creds",
				map[string]string{"apikey": "old-key", "url": "PLACEHOLDER"},
				map[string]string{AnnotationRotate: "false"},
			),
			// should be a no-op
		},
		"KionAPIError": {
			secret: newFakeSecret("default", "creds",
				map[string]string{"apikey": "old-key", "url": "PLACEHOLDER"},
				map[string]string{AnnotationRotate: "true"},
			),
			serverStatus: http.StatusUnauthorized,
			wantStatus:   "error:", // partial match is fine
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Spin up a fake Kion API server.
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if tc.serverStatus == 0 {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(tc.serverStatus)
				if tc.serverKey != "" {
					resp := rotateResponse{}
					resp.Status = tc.serverStatus
					resp.Data.Key = tc.serverKey
					resp.Data.ID = 1
					if err := json.NewEncoder(w).Encode(resp); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}
			}))
			defer server.Close()

			// Patch the Secret's url to point at our test server.
			if tc.secret.Data != nil {
				var creds map[string]string
				if err := json.Unmarshal(tc.secret.Data["credentials"], &creds); err == nil {
					if _, ok := creds["url"]; ok {
						creds["url"] = server.URL
						raw, _ := json.Marshal(creds)
						tc.secret.Data["credentials"] = raw
					}
				}
			}

			scheme := runtime.NewScheme()
			_ = corev1.AddToScheme(scheme)
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tc.secret).
				Build()

			rec := &Reconciler{
				client: fakeClient,
				log:    logging.NewNopLogger(),
				http:   &http.Client{Timeout: 5 * time.Second},
			}

			result, err := rec.Reconcile(t.Context(), reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: tc.secret.Namespace,
					Name:      tc.secret.Name,
				},
			})

			// For most cases we don't expect a hard error (errors become retry-with-status).
			if err != nil && tc.wantNewKey != "" {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantRequeue && result.RequeueAfter <= 0 {
				t.Error("expected RequeueAfter > 0 for not-yet-due rotation")
			}

			// Re-fetch the Secret and check the outcome.
			updated := &corev1.Secret{}
			if fetchErr := fakeClient.Get(t.Context(), types.NamespacedName{
				Namespace: tc.secret.Namespace,
				Name:      tc.secret.Name,
			}, updated); fetchErr != nil {
				t.Fatalf("failed to re-fetch Secret: %v", fetchErr)
			}

			if tc.wantNewKey != "" {
				var creds map[string]string
				if err := json.Unmarshal(updated.Data["credentials"], &creds); err != nil {
					t.Fatalf("unmarshal updated creds: %v", err)
				}
				if got := creds["apikey"]; got != tc.wantNewKey {
					t.Errorf("apikey = %q, want %q", got, tc.wantNewKey)
				}
				if updated.Annotations[AnnotationLastRotation] == "" {
					t.Error("expected last-rotation annotation to be set")
				}
			}

			if tc.wantStatus != "" {
				got := updated.Annotations[AnnotationRotationStatus]
				if tc.wantStatus == "error:" {
					// Partial match for error prefix.
					if len(got) < 6 || got[:6] != "error:" {
						t.Errorf("rotation-status = %q, want prefix 'error:'", got)
					}
				} else if got != tc.wantStatus {
					t.Errorf("rotation-status = %q, want %q", got, tc.wantStatus)
				}
			}
		})
	}
}
