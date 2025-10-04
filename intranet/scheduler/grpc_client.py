import os
import grpc
import sys

from google.protobuf import empty_pb2

# Add the parent 'generated' directory to the Python path.
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '../../..')))

from api.protos.url_fetcher.gen import start_fetching_pb2_grpc
# Import the centralized settings.
from config import settings


def trigger_grpc_parsing():
    address = settings.PARSER_SERVER_ADDRESS
    print(f"Attempting to call Parser server at {address}...")

    try:
        with grpc.insecure_channel(address) as channel:
            stub = start_fetching_pb2_grpc.ProfileServiceStub(channel)

            # The StartFetching RPC call expects an Empty message as a request
            request = empty_pb2.Empty()

            # Call the RPC method on the stub
            stub.StartFetching(request, timeout=10)

            print("Successfully called StartFetching on ProfileService.")
            # The response is also an Empty message, so there isn't much to print from it.
            return {"service": "profile_service", "status_message": "success"}

    except grpc.RpcError as e:
        print(f"An RPC error occurred with the ProfileService: {e}")
        raise e