import httpx

class GoogleAuth(httpx.Auth):
    requires_request_body = True

    def __init__(self):
        self._credentials = None
        self._request = None

    def _get_token(self, force_refresh=False):
        import google.auth
        from google.auth.transport.requests import Request
        if not self._credentials or force_refresh:
            self._credentials, _ = google.auth.default(scopes=['https://www.googleapis.com/auth/cloud-platform'])
            self._request = Request()
        
        if not self._credentials.valid or force_refresh:
            try:
                self._credentials.refresh(self._request)
            except Exception as e:
                print(f"Error fetching/refreshing credentials: {e}")
        return self._credentials.token

    def auth_flow(self, request):
        try:
            token = self._get_token()
            if token:
                request.headers["Authorization"] = f"Bearer {token}"
        except Exception as e:
            print(f"GoogleAuth error: {e}")

        from opentelemetry.propagate import inject
        otel_headers = {}
        inject(otel_headers)
        for k, v in otel_headers.items():
            request.headers[k] = v

        response = yield request

        if response.status_code == 401:
            print("Received 401 UNAUTHENTICATED. Attempting to force refresh credentials and retry...")
            try:
                token = self._get_token(force_refresh=True)
                if token:
                    request.headers["Authorization"] = f"Bearer {token}"
                    yield request
            except Exception as e:
                print(f"Error on retry token fetch: {e}")

