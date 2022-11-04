package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	downloopv1 "github.com/downloop/api/pkg/api/v1"
	"github.com/golang-jwt/jwt"
	"github.com/int128/oauth2cli"
	"github.com/int128/oauth2cli/oauth2params"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

const authEndpoint = "https://downloop.us.auth0.com/"
const clientID = "UoYY0BzQWm0850zD3ASjTw9IjkXueWSK"

type oidcConfig struct {
	config       oauth2.Config
	idToken      string
	refreshToken string
}

func extractExpiry(idToken string, deltaSeconds int) (*time.Time, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not extract token claims")
	}

	var tm time.Time
	switch exp := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(exp), 0)
	case json.Number:
		v, _ := exp.Int64()
		tm = time.Unix(v, 0)
	}

	tm.Add(time.Duration((-1 * deltaSeconds)) * time.Second)
	return &tm, nil
}

func (c oidcConfig) newClient() (*http.Client, error) {

	exp, err := extractExpiry(c.idToken, 10)
	if err != nil {
		return nil, err
	}

	token := oauth2.Token{
		AccessToken:  c.idToken,
		RefreshToken: c.refreshToken,
		Expiry:       *exp,
	}

	ts := c.config.TokenSource(context.Background(), &token)
	rts := oauth2.ReuseTokenSource(&token, ts)

	baseTransport := oidcHTTPTransport{
		T:   http.DefaultTransport,
		rts: rts,
	}

	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: rts,
			Base:   &baseTransport,
		},
		//Timeout: time.Duration(HTTPRequestTimeout) * time.Second,
	}
	return client, nil
}

type oidcHTTPTransport struct {
	T   http.RoundTripper
	rts oauth2.TokenSource
}

func (t *oidcHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	tok, _ := t.rts.Token()

	idToken := tok.Extra("id_token")
	if idToken != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", idToken))
	}
	
	return t.T.RoundTrip(req)
}

func newOauthConfig(endpoint string) (*oauth2.Config, error) {
	scopes := []string{"profile", "email", "openid", "offline_access"}
	provider, err := oidc.NewProvider(context.Background(), endpoint)
	if err != nil {
		return nil, err
	}
	return &oauth2.Config{
		ClientID: clientID,
		Endpoint: provider.Endpoint(),
		Scopes:   scopes,
	}, nil
}

func newClient() (*downloopv1.ClientWithResponses, error) {
	conf, err := ReadConfig()
	if err != nil {
		return nil, err
	}

	oauth2Config, err := newOauthConfig(authEndpoint)
	if err != nil {
		return nil, err
	}
	oidcConfig := oidcConfig{
		config:       *oauth2Config,
		idToken:      conf.Token,
		refreshToken: conf.RefreshToken,
	}

	oidcClient, err := oidcConfig.newClient()
	if err != nil {
		return nil, err
	}

	return downloopv1.NewClientWithResponses(conf.Endpoint, downloopv1.WithHTTPClient(oidcClient))
}

func authenticateUser(c *cli.Context) error {
	pkce, err := oauth2params.NewPKCE()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	ready := make(chan string, 1)
	defer close(ready)

	oauth2Config, err := newOauthConfig(authEndpoint)
	if err != nil {
		return err
	}

	cfg := oauth2cli.Config{
		OAuth2Config:           *oauth2Config,
		AuthCodeOptions:        pkce.AuthCodeOptions(),
		TokenRequestOptions:    pkce.TokenRequestOptions(),
		LocalServerReadyChan:   ready,
		LocalServerBindAddress: []string{"localhost:30000"},
		//LocalServerCertFile:  o.localServerCert,
		//LocalServerKeyFile:   o.localServerKey,
		Logf: log.Printf,
	}

	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case url := <-ready:
			log.Printf("Open %s", url)
			if err := browser.OpenURL(url); err != nil {
				log.Printf("could not open the browser: %s", err)
			}
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
		}
	})
	eg.Go(func() error {
		token, err := oauth2cli.GetToken(ctx, cfg)
		if err != nil {
			return fmt.Errorf("could not get a token: %w", err)
		}
		log.Printf("You got a valid token until %s", token.Expiry)
		idToken := token.Extra("id_token").(string)
		conf := Config{
			Token:        idToken,
			RefreshToken: token.RefreshToken,
			Endpoint:     "http://localhost:8080/",
		}

		err = conf.Save()
		if err != nil {
			return err
		}

		return nil
	})
	if err := eg.Wait(); err != nil {
		log.Fatalf("authorization error: %s", err)
	}
	return nil
}
