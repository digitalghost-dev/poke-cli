import requests
from dagster import DagsterRunStatus, RunStatusSensorContext, run_status_sensor


@run_status_sensor(run_status=DagsterRunStatus.SUCCESS, name="discord_success_sensor")
def discord_success_sensor(context: RunStatusSensorContext):
    context.log.info(f"Detected successful run: {context.dagster_run.run_id}")
    try:
        response = requests.post(
            "https://digitalghost-dev.app.n8n.cloud/webhook/3a58517d-c027-44fa-974c-aedc0035c4f7",
            json={
                "job_name": context.dagster_run.job_name,
                "status": "SUCCESS",
                "run_id": context.dagster_run.run_id,
            },
            timeout=10,
        )
        context.log.info(f"n8n response: {response.status_code}")
    except Exception as e:
        context.log.error(f"Failed to send notification: {e}")


@run_status_sensor(run_status=DagsterRunStatus.FAILURE, name="discord_failure_sensor")
def discord_failure_sensor(context: RunStatusSensorContext):
    context.log.info(f"Detected failed run: {context.dagster_run.run_id}")
    try:
        response = requests.post(
            "https://digitalghost-dev.app.n8n.cloud/webhook/3a58517d-c027-44fa-974c-aedc0035c4f7",
            json={
                "job_name": context.dagster_run.job_name,
                "status": "FAILURE",
                "run_id": context.dagster_run.run_id,
            },
            timeout=10,
            
        )
        context.log.info(f"n8n response: {response.status_code}")
    except Exception as e:
        context.log.error(f"Failed to send notification: {e}")
