import grpc
from google.protobuf import empty_pb2

from api.protos.url_fetcher.gen import start_fetching_pb2_grpc
from parser import parse_dzen_for_links
from repo import LinkStorer

class ProfileServiceServicer(start_fetching_pb2_grpc.ProfileServiceServicer):
    """
    Implements the gRPC ProfileService.
    This class is decoupled from the data storage mechanism via an interface.
    """
    def __init__(self, link_storer: LinkStorer):
        """
        Initializes the servicer with a dependency that conforms to the
        LinkStorerInterface.

        Args:
            link_storer: An object that implements the LinkStorerInterface.
        """
        print("ProfileServiceServicer initialized.")
        self.link_storer = link_storer

    def StartFetching(self, request: empty_pb2.Empty, context) -> empty_pb2.Empty:
        print("gRPC call received: StartFetching.")

        try:
            links_to_store = parse_dzen_for_links()

            if links_to_store:
                num_added = self.link_storer.store_links(links_to_store)
                print(f"Link storer reported {num_added} new links.")
            else:
                print("No new links found to store.")

        except ConnectionError as e:
            print(f"Error during execution: {e}")
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details(f"A downstream service is unavailable: {e}")

        except Exception as e:
            print(f"An unexpected error occurred: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details("An internal server error occurred during fetching.")

        return empty_pb2.Empty()