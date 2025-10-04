import grpc
from concurrent import futures

from api.protos.url_fetcher.gen import start_fetching_pb2_grpc
from .grpc_server import ProfileServiceServicer
from .config import settings
from .redis_adapter import RedisPublisherAdapter


def serve():
    """
    Starts the gRPC server. This is the Composition Root where dependencies
    are created and injected.
    """

    redis_storer = RedisPublisherAdapter()

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    servicer = ProfileServiceServicer(link_storer=redis_storer)

    start_fetching_pb2_grpc.add_ProfileServiceServicer_to_server(servicer, server)

    port = settings.DZEN_URL_PARSER_PORT
    server.add_insecure_port(f"[::]:{port}")

    print(f"dzen_url_parser gRPC server started on port {port}...")
    server.start()

    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        print("Server shutting down...")
        server.stop(0)

if __name__ == "__main__":
    serve()