import grpc
from fastapi import FastAPI, HTTPException
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
async def manual_trigger():

    print("Manual trigger received.")
    try:
        trigger_grpc_parsing()
    except grpc.RpcError as e:
        # Convert a gRPC error into a user-friendly HTTP 503 Service Unavailable error.
        raise HTTPException(
            status_code=500,
            detail=f"Failed to connect to parser service: {e}"
        )
    except Exception as e:
        # Catch any other unexpected errors and return a 500 Internal Server Error.
        print(f"An unexpected error occurred: {e}")
        raise HTTPException(status_code=500, detail="An internal server error occurred.")

# --- Scheduler Setup ---
@app.on_event("startup")
async def startup_event():

    print("Scheduler starting...")
    scheduler.add_job(
        run_scheduled_parsing,
        trigger=IntervalTrigger(minutes=20),
        id="parsing_job",
        name="Trigger gRPC parsing every 20 minutes",
        replace_existing=True,
    )
    scheduler.start()
    print("Scheduler started. Job 'parsing_job' is scheduled every 20 minutes.")

@app.on_event("shutdown")
async def shutdown_event():
    """This function runs when the FastAPI application is shutting down."""
    print("Scheduler shutting down...")
    scheduler.shutdown()
    print("Scheduler has been shut down.")