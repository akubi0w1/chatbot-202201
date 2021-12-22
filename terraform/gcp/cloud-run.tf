############################
# cloud run
resource "google_cloud_run_service" "chatbot" {
  name     = "${local.slack_workspace}-${local.service}"
  location = var.region

  template {
    spec {
      service_account_name = google_service_account.chatbot.email
      containers {
        image = "gcr.io/cloudrun/hello"
        env {
          name = "SLACK_BOT_TOKEN"
          value = var.slack_bot_token
        }
        env {
          name = "SLACK_SIGNING_SECRET"
          value = var.slack_signing_secret
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [
      template,
    ]
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.chatbot.location
  service = google_cloud_run_service.chatbot.name
  project     = google_cloud_run_service.chatbot.project
  policy_data = data.google_iam_policy.noauth.policy_data
}

############################
# service account
resource "google_service_account" "chatbot" {
    account_id = "akubi-post-chatbot-executer"
    display_name = "akubi-post chatbot executer"
    description = "slack bot executer"
}
