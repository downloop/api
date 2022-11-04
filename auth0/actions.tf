terraform {
  required_providers {
    auth0 = {
      source  = "auth0/auth0"
      version = "~> 0.39.0"
    }
  }
}

provider "auth0" {
	domain = "downloop.us.auth0.com"
	client_id = local.provider_client_id
	client_secret = local.provider_client_secret

}

resource "auth0_action" "social_post_login" {
  name    = "Social User Post-Login"
  runtime = "node16"
  deploy  = true
  code    = file("./social-post-login.js")

  supported_triggers {
    id      = "post-login"
    version = "v3"
  }

  dependencies {
    name    = "axios"
    version = "1.1.3"
  }

  secrets {
    name  = "AUTH0_CLIENT_ID"
    value = local.auth0_client_id
  }

  secrets {
    name  = "AUTH0_CLIENT_SECRET"
    value = local.auth0_client_secret
  }
}

resource "auth0_action" "login_add_claims" {
  name    = "Add Token Claims"
  runtime = "node16"
  deploy  = true
  code    = file("./login-add-claims.js")

  supported_triggers {
    id      = "post-login"
    version = "v3"
  }
}

resource "auth0_trigger_binding" "login_flow" {
  trigger = "post-login"

  actions {
    id           = auth0_action.social_post_login.id
    display_name = auth0_action.social_post_login.name
  }

  actions {
    id           = auth0_action.login_add_claims.id
    display_name = auth0_action.login_add_claims.name
  }
}