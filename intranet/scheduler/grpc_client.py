import os
import sys
import time
import random
import grpc
from google.protobuf import empty_pb2

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '../../..')))
from api.protos.url_fetcher.gen import start_fetching_pb2_grpc
from .config import settings

def _channel_options():
    return [
        ("grpc.keepalive_time_ms", settings.GRPC_KEEPALIVE_TIME_MS),
        ("grpc.keepalive_timeout_ms", settings.GRPC_KEEPALIVE_TIMEOUT_MS),
        ("grpc.keepalive_permit_without_calls", settings.GRPC_KEEPALIVE_PERMIT_WITHOUT_CALLS),
        ("grpc.http2.max_pings_without_data", settings.GRPC_MAX_PINGS_WITHOUT_DATA),
    ]

def trigger_grpc_parsing():
    address = settings.construct_server_address
    print(f"Attempting to call Parser server at {address}...")

    request = empty_pb2.Empty()
    attempts = settings.GRPC_RETRY_MAX
    backoff_ms = settings.GRPC_RETRY_BACKOFF_MS

    with grpc.insecure_channel(address, options=_channel_options()) as channel:
        # Ждать готовности канала, чтобы избежать раннего UNAVAILABLE
        grpc.channel_ready_future(channel).result(timeout=5)

        stub = start_fetching_pb2_grpc.ProfileServiceStub(channel)

        for i in range(attempts):
            try:
                stub.StartFetching(request, timeout=settings.GRPC_TIMEOUT_SECONDS)
                print("Successfully called StartFetching on ProfileService.")
                return {"service": "profile_service", "status_message": "success"}
            except grpc.RpcError as e:
                code = e.code()
                if code in (grpc.StatusCode.UNAVAILABLE, grpc.StatusCode.DEADLINE_EXCEEDED) and i < attempts - 1:
                    sleep_ms = backoff_ms * (2 ** i)
                    jitter_ms = random.randint(0, 250)
                    time.sleep((sleep_ms + jitter_ms) / 1000.0)
                    continue
                raise e
        return None
