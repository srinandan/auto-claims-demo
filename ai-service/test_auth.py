import asyncio
import httpx
from google.auth.transport.requests import Request
from google.auth import default

class GoogleAuth(httpx.Auth):
    def __init__(self):
        self.credentials, _ = default(scopes=['https://www.googleapis.com/auth/cloud-platform'])
        self.request = Request()

    def auth_flow(self, request):
        if self.credentials:
            if not self.credentials.valid:
                self.credentials.refresh(self.request)
            if self.credentials.token:
                request.headers["Authorization"] = f"Bearer {self.credentials.token}"
        
        response = yield request
        
        if response.status_code == 401:
            print("Received 401, refreshing and retrying...")
            self.credentials.refresh(self.request)
            request.headers["Authorization"] = f"Bearer {self.credentials.token}"
            yield request

async def main():
    async with httpx.AsyncClient(auth=GoogleAuth()) as client:
        # Mock vertex ai endpoint or just a real one
        resp = await client.get("https://aiplatform.googleapis.com/v1/projects")
        print(resp.status_code, resp.text)

asyncio.run(main())
