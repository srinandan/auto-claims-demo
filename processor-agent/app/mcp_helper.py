import asyncio
from mcp import ClientSession
from mcp.client.sse import sse_client

async def lookup_address(address_query: str):
    """Call Google Maps Grounding Lite MCP search_places tool to resolve the address."""
    import os
    # Assuming the API key or ADC is configured.
    headers = {}
    api_key = os.environ.get("GOOGLE_MAPS_API_KEY")
    if api_key:
        headers["x-goog-api-key"] = api_key

    url = "https://mapstools.googleapis.com/mcp"
    try:
        # mcp client via SSE
        async with sse_client(url, headers=headers) as streams:
             async with ClientSession(streams[0], streams[1]) as session:
                 await session.initialize()
                 result = await session.call_tool("search_places", arguments={"text_query": address_query})
                 for content in result.content:
                     if content.type == "text":
                         return content.text
        return None
    except Exception as e:
        import logging
        logging.error(f"Error calling MCP server: {e}")
        return None
