############################
# cloud scheduler
resource "google_cloud_scheduler_job" "timer" {
  name             = "${local.slack_workspace}-${local.service}-schedule"
  description      = "interval"
  schedule         = "0 */1 * * *"
  time_zone        = "Asia/Tokyo"
  attempt_deadline = "320s"

  retry_config {
    retry_count = 0
    min_backoff_duration = "5s"
    max_retry_duration = "0s"
    max_doublings = 5
  }

  http_target {
    http_method = "POST"
    uri         = "${google_cloud_run_service.chatbot.status[0].url}/watch"
    headers = {
        "Content-Type" = "application/json"
    }

    oidc_token {
      service_account_email = google_service_account.chatbot.email
    }
  }
}
