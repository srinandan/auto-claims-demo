import asyncio
import httpx
from google.auth.transport.requests import Request
from google.auth import default

class GoogleAuth(httpx.Auth):
    def __init__(self):
        self.credentials, _ = default(scopes=['https://www.googleapis.com/auth/cloud-platform'])
        self.request = Request()
        print(f"init valid: {self.credentials.valid}, token: {self.credentials.token is not None}")

    def auth_flow(self, request):
        if self.credentials:
            if not self.credentials.valid:
                try:
                    self.credentials.refresh(self.request)
                    print(f"refreshed valid: {self.credentials.valid}")
                except Exception as e:
                    print(f"Error refreshing credentials: {e}")
            if self.credentials.token:
                request.headers["Authorization"] = f"Bearer {self.credentials.token}"
                print("added auth header")

        response = yield request
        
        if response.status_code == 401:
            print("Received 401 UNAUTHENTICATED. Attempting to refresh credentials and retry...")
            if self.credentials:
                try:
                    self.credentials.refresh(self.request)
                except Exception as e:
                    print(f"Error refreshing credentials on 401: {e}")
                if self.credentials.token:
                    request.headers["Authorization"] = f"Bearer {self.credentials.token}"
                    yield request

async def main():
    async with httpx.AsyncClient(auth=GoogleAuth()) as client:
        print("sending request")
        resp = await client.get("https://us-east4-aiplatform.googleapis.com/v1beta1/projects/srinandans-apphub/locations/us-east4/reasoningEngines")
        print("status:", resp.status_code)
        
asyncio.run(main())
