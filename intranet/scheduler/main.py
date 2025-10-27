import grpc
from fastapi import FastAPI, HTTPException, BackgroundTasks
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.triggers.interval import IntervalTrigger

# Import the refactored modules.
from .grpc_client import trigger_grpc_parsing
from .config import settings

# --- FastAPI App Initialization ---
app = FastAPI(
    title="Scheduler Service",
    description="A service to trigger a gRPC parser manually or on a schedule.",
)
scheduler = AsyncIOScheduler()

def run_scheduled_parsing():
    print("Scheduled trigger fired.")
    try:
        trigger_grpc_parsing()
    except Exception as e:
        # For a background job, we log the error instead of raising it.
        print(f"Scheduled job failed with error: {e}")

@app.post("/trigger-parsing", summary="Manually trigger the parsing service")
async def manual_trigger(bg: BackgroundTasks):
    print("Manual trigger received.")
    # run in background to avoid blocking HTTP response
    bg.add_task(trigger_grpc_parsing)
    return {"status": "accepted"}

# --- Scheduler Setup ---
@app.on_event("startup")
async def startup_event():
    print("Scheduler starting...")
    scheduler.add_job(
        run_scheduled_parsing,
        trigger=IntervalTrigger(minutes=60),
        id="parsing_job",
        name="Trigger gRPC parsing every 20 minutes",
        replace_existing=True,
        coalesce=True,
        max_instances=1,
        misfire_grace_time=300,
    )
    scheduler.start()
    print("Scheduler started. Job 'parsing_job' is scheduled every 20 minutes.")

@app.on_event("shutdown")
async def shutdown_event():
    """This function runs when the FastAPI application is shutting down."""
    print("Scheduler shutting down...")
    scheduler.shutdown()
    print("Scheduler has been shut down.")
