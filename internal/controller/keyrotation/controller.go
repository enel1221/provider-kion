/*
Copyright 2024 Upbound Inc.
*/

// Package keyrotation implements a controller that watches Kubernetes Secrets
// annotated with kion.upbound.io/rotate=true and automatically rotates the
// Kion API key using the POST /api/v3/app-api-key/rotate endpoint.
//
// Usage: annotate the credentials Secret referenced by your ProviderConfig:
//
//	apiVersion: v1
//	kind: Secret
//	metadata:
//	  name: kion-creds
//	  namespace: crossplane-system
//	  annotations:
//	    kion.upbound.io/rotate: "true"
//	    kion.upbound.io/rotation-interval: "168h"   # optional, default 7 days
//	    kion.upbound.io/credentials-key: "credentials" # optional, default "credentials"
//	type: Opaque
//	stringData:
//	  credentials: |
//	    {
//	      "apikey": "your_api_key",
//	      "url": "https://kion.example.com",
//	      "skipsslvalidation": "false"
//	    }
package keyrotation

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
)

const (
	// User-configurable annotations.

	// AnnotationRotate enables rotation when set to "true".
	AnnotationRotate = "kion.upbound.io/rotate"
	// AnnotationRotationInterval sets how often the key is rotated (Go duration, e.g. "168h").
	AnnotationRotationInterval = "kion.upbound.io/rotation-interval"
	// AnnotationCredentialsKey overrides the Secret data key that holds the credentials JSON.
	AnnotationCredentialsKey = "kion.upbound.io/credentials-key"

	// Controller-managed annotations (written back to the Secret).

	// AnnotationLastRotation records the timestamp of the last successful rotation.
	AnnotationLastRotation = "kion.upbound.io/last-rotation"
	// AnnotationRotationStatus records the result of the last rotation attempt.
	AnnotationRotationStatus = "kion.upbound.io/rotation-status"

	// Defaults.

	// DefaultRotationInterval is used when no rotation-interval annotation is set.
	// 7 days keeps a comfortable margin within a 14-day expiry window.
	DefaultRotationInterval = 168 * time.Hour
	// DefaultCredentialsKey is the default Secret data key for the credentials JSON.
	DefaultCredentialsKey = "credentials"

	rotateAPIPath = "/api/v3/app-api-key/rotate"
	httpTimeout   = 30 * time.Second
	retryInterval = 5 * time.Minute
)

// Reconciler watches annotated Secrets and rotates the Kion API key.
type Reconciler struct {
	client client.Client
	log    logging.Logger
	http   *http.Client
}

// Setup creates the key-rotation controller and registers it with the manager.
//
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;update;patch
//
//nolint:lll // kubebuilder markers are long
func Setup(mgr ctrl.Manager, log logging.Logger) error {
	r := &Reconciler{
		client: mgr.GetClient(),
		log:    log.WithValues("controller", "keyrotation"),
		http:   &http.Client{Timeout: httpTimeout},
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named("keyrotation").
		For(&corev1.Secret{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc:  func(e event.CreateEvent) bool { return hasRotateAnnotation(e.Object) },
			UpdateFunc:  func(e event.UpdateEvent) bool { return hasRotateAnnotation(e.ObjectNew) },
			DeleteFunc:  func(_ event.DeleteEvent) bool { return false },
			GenericFunc: func(_ event.GenericEvent) bool { return false },
		}).
		Complete(r)
}

func hasRotateAnnotation(obj client.Object) bool {
	a := obj.GetAnnotations()
	return a != nil && strings.EqualFold(a[AnnotationRotate], "true")
}

// Reconcile checks whether the Kion API key needs rotation and, if so, rotates
// it via the Kion API and writes the new key back into the Secret.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("secret", req.NamespacedName)

	secret := &corev1.Secret{}
	if err := r.client.Get(ctx, req.NamespacedName, secret); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if !hasRotateAnnotation(secret) {
		return reconcile.Result{}, nil
	}

	annotations := secret.GetAnnotations()

	// ── Determine rotation interval ──────────────────────────────────
	interval := DefaultRotationInterval
	if v, ok := annotations[AnnotationRotationInterval]; ok {
		if parsed, err := time.ParseDuration(v); err == nil {
			interval = parsed
		} else {
			log.Info("Invalid rotation-interval annotation, using default", "value", v)
		}
	}

	// ── Check if rotation is due ─────────────────────────────────────
	if v, ok := annotations[AnnotationLastRotation]; ok {
		if lastRotation, err := time.Parse(time.RFC3339, v); err == nil {
			remaining := time.Until(lastRotation.Add(interval))
			if remaining > 0 {
				log.Debug("Rotation not yet due", "remaining", remaining.String())
				return reconcile.Result{RequeueAfter: remaining}, nil
			}
		}
	}

	// ── Extract credentials from Secret ──────────────────────────────
	credKey := DefaultCredentialsKey
	if v, ok := annotations[AnnotationCredentialsKey]; ok && v != "" {
		credKey = v
	}

	raw, ok := secret.Data[credKey]
	if !ok {
		log.Info("Credentials key not found in secret data", "key", credKey)
		r.setStatus(ctx, secret, fmt.Sprintf("error: key %q not found in secret data", credKey))
		return reconcile.Result{RequeueAfter: retryInterval}, nil
	}

	var creds map[string]interface{}
	if err := json.Unmarshal(raw, &creds); err != nil {
		log.Info("Failed to parse credentials JSON", "error", err)
		r.setStatus(ctx, secret, fmt.Sprintf("error: invalid credentials JSON: %v", err))
		return reconcile.Result{RequeueAfter: retryInterval}, nil
	}

	apiKey, _ := creds["apikey"].(string)
	kionURL, _ := creds["url"].(string)
	skipSSL := parseBool(creds["skipsslvalidation"])

	if apiKey == "" || kionURL == "" {
		log.Info("Missing apikey or url in credentials")
		r.setStatus(ctx, secret, "error: credentials JSON must contain non-empty apikey and url")
		return reconcile.Result{RequeueAfter: retryInterval}, nil
	}

	// ── Call Kion rotate API ─────────────────────────────────────────
	log.Info("Rotating Kion API key")
	newKey, err := r.callRotate(ctx, kionURL, apiKey, skipSSL)
	if err != nil {
		log.Info("API key rotation failed", "error", err)
		r.setStatus(ctx, secret, fmt.Sprintf("error: %v", err))
		return reconcile.Result{RequeueAfter: retryInterval}, nil
	}

	// ── Update Secret with the new key ───────────────────────────────
	creds["apikey"] = newKey
	updated, err := json.Marshal(creds)
	if err != nil {
		log.Info("Failed to marshal updated credentials", "error", err)
		return reconcile.Result{RequeueAfter: retryInterval}, nil
	}

	secret.Data[credKey] = updated
	if secret.Annotations == nil {
		secret.Annotations = make(map[string]string)
	}
	secret.Annotations[AnnotationLastRotation] = time.Now().UTC().Format(time.RFC3339)
	secret.Annotations[AnnotationRotationStatus] = "success"

	if err := r.client.Update(ctx, secret); err != nil {
		// The rotation succeeded at Kion but the Secret update failed.
		// The controller will retry on the next reconcile.
		log.Info("CRITICAL: key rotated at Kion but Secret update failed — will retry", "error", err)
		return reconcile.Result{Requeue: true}, err
	}

	log.Info("Successfully rotated Kion API key")
	return reconcile.Result{RequeueAfter: interval}, nil
}

// ── Kion API client ──────────────────────────────────────────────────────────

// rotateResponse mirrors the Kion API response envelope for key rotation.
type rotateResponse struct {
	Status int `json:"status"`
	Data   struct {
		ID  int    `json:"id"`
		Key string `json:"key"`
	} `json:"data"`
}

// callRotate issues POST /api/v3/app-api-key/rotate to the Kion instance and
// returns the freshly-issued API key.
func (r *Reconciler) callRotate(ctx context.Context, kionURL, currentKey string, skipSSL bool) (string, error) {
	body, err := json.Marshal(map[string]string{"key": currentKey})
	if err != nil {
		return "", fmt.Errorf("marshal rotate request: %w", err)
	}

	endpoint := strings.TrimRight(kionURL, "/") + rotateAPIPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+currentKey)

	httpClient := r.http
	if skipSSL {
		httpClient = &http.Client{
			Timeout: httpTimeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // user-configured via skipsslvalidation
			},
		}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call rotate API: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1 MB safety limit
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("rotate API returned HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var rr rotateResponse
	if err := json.Unmarshal(respBody, &rr); err != nil {
		return "", fmt.Errorf("unmarshal rotate response: %w", err)
	}

	if rr.Data.Key == "" {
		return "", fmt.Errorf("rotate API returned an empty key")
	}

	return rr.Data.Key, nil
}

// ── Helpers ──────────────────────────────────────────────────────────────────

// setStatus writes the rotation-status annotation without failing the reconcile.
func (r *Reconciler) setStatus(ctx context.Context, secret *corev1.Secret, status string) {
	if secret.Annotations == nil {
		secret.Annotations = make(map[string]string)
	}
	secret.Annotations[AnnotationRotationStatus] = status
	if err := r.client.Update(ctx, secret); err != nil {
		r.log.Info("Failed to update rotation-status annotation", "error", err)
	}
}

// parseBool handles both JSON string "true" and JSON bool true.
func parseBool(v interface{}) bool {
	switch val := v.(type) {
	case string:
		return strings.EqualFold(val, "true")
	case bool:
		return val
	default:
		return false
	}
}
